package rules

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

// RightmostOperatorRule checks that provider version constraints use ~> x.0 format
type RightmostOperatorRule struct {
	tflint.DefaultRule
}

func NewRightmostOperatorRule() *RightmostOperatorRule {
	return &RightmostOperatorRule{}
}

func (r *RightmostOperatorRule) Name() string {
	return "rightmost_operator_rule"
}

func (r *RightmostOperatorRule) Enabled() bool {
	return true
}

func (r *RightmostOperatorRule) Severity() tflint.Severity {
	return tflint.WARNING
}

func (r *RightmostOperatorRule) Link() string {
	return ""
}

// Check checks whether provider version constraints use ~> x.0 format
func (r *RightmostOperatorRule) Check(runner tflint.Runner) error {
	// Walk through all files to find provider configurations
	files, err := runner.GetFiles()
	if err != nil {
		return err
	}

	for filename, file := range files {
		if err := r.checkFile(runner, file, filename); err != nil {
			return err
		}
	}

	return nil
}

func (r *RightmostOperatorRule) checkFile(runner tflint.Runner, file *hcl.File, filename string) error {
	// Parse the file body
	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		return nil
	}

	// Check terraform blocks
	for _, block := range body.Blocks {
		if block.Type == "terraform" {
			if err := r.checkTerraformBlock(runner, block); err != nil {
				return err
			}
		} else if block.Type == "provider" && len(block.Labels) > 0 {
			if err := r.checkProviderBlock(runner, block); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *RightmostOperatorRule) checkTerraformBlock(runner tflint.Runner, block *hclsyntax.Block) error {
	// Look for required_providers blocks
	for _, innerBlock := range block.Body.Blocks {
		if innerBlock.Type == "required_providers" {
			// Process each provider in the required_providers block
			for name, attr := range innerBlock.Body.Attributes {
				// Evaluate the attribute
				val, diags := attr.Expr.Value(nil)
				if diags.HasErrors() {
					continue
				}

				// Check if it's an object with a version field
				if val.Type().IsObjectType() && val.Type().HasAttribute("version") {
					versionVal := val.GetAttr("version")
					if !versionVal.IsNull() && versionVal.Type() == cty.String {
						version := versionVal.AsString()
						if err := r.checkVersionConstraint(runner, version, attr.Expr.Range(), fmt.Sprintf("Provider %s", name)); err != nil {
							return err
						}
					}
				}
			}
		}
	}

	return nil
}

func (r *RightmostOperatorRule) checkProviderBlock(runner tflint.Runner, block *hclsyntax.Block) error {
	// Check for version attribute
	if attr, exists := block.Body.Attributes["version"]; exists {
		val, diags := attr.Expr.Value(nil)
		if !diags.HasErrors() && val.Type() == cty.String {
			version := val.AsString()
			if err := r.checkVersionConstraint(runner, version, attr.Expr.Range(), fmt.Sprintf("Provider %s", block.Labels[0])); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *RightmostOperatorRule) checkVersionConstraint(runner tflint.Runner, version string, rng hcl.Range, providerName string) error {
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
