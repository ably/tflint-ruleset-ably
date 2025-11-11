package rules

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

// AwsModuleVersionRule checks that when using community terraform-aws-modules,
// the version matches that of the AWS provider, as well as ensuring the
// rightmost operator
type AwsModuleVersionRule struct {
	tflint.DefaultRule
}

func NewAwsModuleVersionRule() *AwsModuleVersionRule {
	return &AwsModuleVersionRule{}
}

func (r *AwsModuleVersionRule) Name() string {
	return "aws_module_version_rule"
}

func (r *AwsModuleVersionRule) Enabled() bool {
	return true
}

func (r *AwsModuleVersionRule) Severity() tflint.Severity {
	return tflint.WARNING
}

func (r *AwsModuleVersionRule) Link() string {
	return ""
}

// Check checks whether module version constraints match AWS provider version
func (r *AwsModuleVersionRule) Check(runner tflint.Runner) error {
	// First, get the AWS provider version from terraform blocks
	awsProviderVersion, err := r.getAwsProviderVersion(runner)
	if err != nil {
		return err
	}

	// If no AWS provider version found, we don't need to check modules
	if awsProviderVersion == "" {
		return nil
	}

	// Walk through all files to find module configurations
	files, err := runner.GetFiles()
	if err != nil {
		return err
	}

	for _, file := range files {
		if err := r.checkFile(runner, file, awsProviderVersion); err != nil {
			return err
		}
	}

	return nil
}

func (r *AwsModuleVersionRule) checkFile(runner tflint.Runner, file *hcl.File, awsProviderVersion string) error {
	// Parse the file body
	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		return nil
	}

	// Check module blocks
	for _, block := range body.Blocks {
		if block.Type != "module" {
			continue
		}

		if err := r.checkModuleBlock(runner, block, awsProviderVersion); err != nil {
			return err
		}
	}

	return nil
}

// getAwsProviderVersion finds the AWS provider version from terraform required_providers blocks
func (r *AwsModuleVersionRule) getAwsProviderVersion(runner tflint.Runner) (string, error) {
	files, err := runner.GetFiles()
	if err != nil {
		return "", err
	}

	for _, file := range files {
		body, ok := file.Body.(*hclsyntax.Body)
		if !ok {
			continue
		}

		for _, block := range body.Blocks {
			if block.Type == "terraform" {
				for _, innerBlock := range block.Body.Blocks {
					if innerBlock.Type == "required_providers" {
						for name, attr := range innerBlock.Body.Attributes {
							if name != "aws" {
								continue
							}

							val, diags := attr.Expr.Value(nil)
							if diags.HasErrors() {
								continue
							}

							if val.Type().IsObjectType() && val.Type().HasAttribute("version") {
								versionVal := val.GetAttr("version")
								if !versionVal.IsNull() && versionVal.Type() == cty.String {
									return versionVal.AsString(), nil
								}
							}
						}
					}
				}
			}
		}
	}

	return "", nil
}

// Pair up the matching major versions between provider and module
type sourceVersion struct {
	Provider int
	Module   int
}

// A definitive map of AWS module major versions to AWS provider major versions. They only started linking
// breaking changes to breaking provider versions in the last few versions.
var sourceVersionMap = map[string][]sourceVersion{
	"terraform-aws-modules/s3-bucket/aws": {
		{
			Provider: 5,
			Module:   4,
		},
		{
			Provider: 6,
			Module:   5,
		},
	},
	"terraform-aws-modules/vpc/aws": {
		{
			Provider: 5,
			Module:   5,
		},
		{
			Provider: 6,
			Module:   6,
		},
	},
	"terraform-aws-modules/lambda/aws": {
		{
			Provider: 5,
			Module:   7,
		},
		{
			Provider: 6,
			Module:   8,
		},
	},
}

func (r *AwsModuleVersionRule) checkModuleBlock(runner tflint.Runner, block *hclsyntax.Block, awsProviderVersion string) error {
	// Get the source attribute
	sourceAttr, exists := block.Body.Attributes["source"]
	if !exists {
		return nil
	}

	// Evaluate the source expression
	sourceVal, diags := sourceAttr.Expr.Value(nil)
	if diags.HasErrors() || sourceVal.Type() != cty.String {
		return nil
	}
	source := sourceVal.AsString()

	// Check if this is a terraform-aws-modules module we care about
	versionMap, exists := sourceVersionMap[source]
	if !exists {
		return nil
	}

	// Get the version attribute
	versionAttr, exists := block.Body.Attributes["version"]
	if !exists {
		// No version specified for terraform-aws-modules, emit a warning
		return runner.EmitIssue(
			r,
			fmt.Sprintf("Module %s should specify a version constraint", source),
			sourceAttr.Expr.Range(),
		)
	}

	// Evaluate the version expression
	versionVal, diags := versionAttr.Expr.Value(nil)
	if diags.HasErrors() || versionVal.Type() != cty.String {
		return nil
	}
	moduleVersion := versionVal.AsString()

	// First check that the module version uses the rightmost operator
	if err := r.checkRightmostOperator(runner, moduleVersion, versionAttr.Expr.Range(),
		fmt.Sprintf("Module %s", source)); err != nil {
		return err
	}

	// Parse the AWS provider version to get the major version
	awsProviderMajor := r.extractMajorVersion(awsProviderVersion)
	if awsProviderMajor == -1 {
		// Can't determine AWS provider major version, skip the compatibility check
		return nil
	}

	// Parse the module version to get the major version
	moduleMajor := r.extractMajorVersion(moduleVersion)
	if moduleMajor == -1 {
		// Can't determine module major version, skip the compatibility check
		return nil
	}

	// Check if this module version is compatible with the AWS provider version
	compatible := false
	for _, mapping := range versionMap {
		if mapping.Provider == awsProviderMajor && mapping.Module == moduleMajor {
			compatible = true
			break
		}
	}

	if !compatible {
		// Find the correct module version for this AWS provider version
		recommendedModule := -1
		for _, mapping := range versionMap {
			if mapping.Provider == awsProviderMajor {
				recommendedModule = mapping.Module
				break
			}
		}

		if recommendedModule != -1 {
			return runner.EmitIssue(
				r,
				fmt.Sprintf("Module %s version ~> %d.0 is not compatible with AWS provider version %s. Use module version ~> %d.0 for AWS provider ~> %d.0",
					source, moduleMajor, awsProviderVersion, recommendedModule, awsProviderMajor),
				versionAttr.Expr.Range(),
			)
		} else {
			return runner.EmitIssue(
				r,
				fmt.Sprintf("Module %s version ~> %d.0 does not have a known compatibility mapping for AWS provider version %s",
					source, moduleMajor, awsProviderVersion),
				versionAttr.Expr.Range(),
			)
		}
	}

	return nil
}

// extractMajorVersion extracts the major version from a version constraint
func (r *AwsModuleVersionRule) extractMajorVersion(versionConstraint string) int {
	// Pattern to match ~> X.Y or ~> X.Y.Z
	pattern := regexp.MustCompile(`^\s*~>\s*(\d+)\.`)
	matches := pattern.FindStringSubmatch(versionConstraint)
	if len(matches) > 1 {
		var major int
		if _, err := fmt.Sscanf(matches[1], "%d", &major); err == nil {
			return major
		}
	}
	return -1
}

// checkRightmostOperator checks that a version constraint uses the rightmost operator (~>)
func (r *AwsModuleVersionRule) checkRightmostOperator(runner tflint.Runner, version string, rng hcl.Range, name string) error {
	// Regular expression to match ~> x.y pattern (major.minor only, no patch)
	validPattern := regexp.MustCompile(`^\s*~>\s*(\d+)\.(\d+)\s*$`)

	if !validPattern.MatchString(version) {
		return runner.EmitIssue(
			r,
			fmt.Sprintf("%s version constraint should use '~> x.y' format where x is the major version and y is the minor version (no patch version), got: %s", name, version),
			rng,
		)
	}

	return nil
}
