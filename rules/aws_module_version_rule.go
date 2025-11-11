package rules

import (
	"fmt"
	"regexp"

	"github.com/Masterminds/semver"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
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

// Check checks whether provider version constraints use ~> x.0 format
func (r *AwsModuleVersionRule) Check(runner tflint.Runner) error {
	// Walk through all files to find provider configurations
	files, err := runner.GetFiles()
	if err != nil {
		return err
	}

	for _, file := range files {
		if err := r.checkFile(runner, file); err != nil {
			return err
		}
	}

	return nil
}

func (r *AwsModuleVersionRule) checkFile(runner tflint.Runner, file *hcl.File) error {
	// Parse the file body
	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		return nil
	}

	// Check terraform blocks
	for _, block := range body.Blocks {
		if block.Type != "module" {
			continue
		}

		if err := r.checkModuleBlock(runner, block.Body.Attributes); err != nil {
			return err
		}
	}

	return nil
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
}

func (r *AwsModuleVersionRule) checkModuleBlock(runner tflint.Runner, attributes hclsyntax.Attributes) error {
	source, exists := attributes["source"]
	if !exists {
		return nil
	}

	switch source.Name {
	case "terraform-aws-modules/s3-bucket/aws":
		version, exists := attributes["version"]
		if !exists {
			return nil
		}

		v, err := semver.NewVersion(version.Name)

	case "terraform-aws-modules/lambda/aws":
	case "terraform-aws-modules/vpc/aws":
	default:
		return nil
	}

	// Look for required_providers blocks
	// for _, innerBlock := range block.Body.Blocks {
	// 	if innerBlock.Type == "required_providers" {
	// 		// Process each provider in the required_providers block
	// 		for name, attr := range innerBlock.Body.Attributes {
	// 			// Evaluate the attribute
	// 			val, diags := attr.Expr.Value(nil)
	// 			if diags.HasErrors() {
	// 				continue
	// 			}

	// 			// Check if it's an object with a version field
	// 			if val.Type().IsObjectType() && val.Type().HasAttribute("version") {
	// 				versionVal := val.GetAttr("version")
	// 				if !versionVal.IsNull() && versionVal.Type() == cty.String {
	// 					version := versionVal.AsString()
	// 					if err := r.checkVersionConstraint(runner, version, attr.Expr.Range(), fmt.Sprintf("Provider %s", name)); err != nil {
	// 						return err
	// 					}
	// 				}
	// 			}
	// 		}
	// 	}
	// }

	return nil
}

func (r *AwsModuleVersionRule) checkVersionConstraint(runner tflint.Runner, version string, rng hcl.Range, providerName string) error {
	// Regular expression to match ~> x.y pattern (major.minor only, no patch)
	validPattern := regexp.MustCompile(`^\s*~>\s*(\d+)\.(\d+)\s*$`)

	if !validPattern.MatchString(version) {
		return runner.EmitIssue(
			r,
			fmt.Sprintf("%s version constraint should use '~> x.y' format where x is the major version and y is the minor version (no patch version), got: %s", providerName, version),
			rng,
		)
	}

	return nil
}
