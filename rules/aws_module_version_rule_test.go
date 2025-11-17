package rules

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsModuleVersionRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "valid s3-bucket module with AWS provider v5",
			Content: `
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

module "s3_bucket" {
  source  = "terraform-aws-modules/s3-bucket/aws"
  version = "~> 4.0"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "valid s3-bucket module with AWS provider v6",
			Content: `
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.0"
    }
  }
}

module "s3_bucket" {
  source  = "terraform-aws-modules/s3-bucket/aws"
  version = "~> 5.0"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "invalid s3-bucket module version for AWS provider v5",
			Content: `
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

module "s3_bucket" {
  source  = "terraform-aws-modules/s3-bucket/aws"
  version = "~> 5.0"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsModuleVersionRule(),
					Message: "Module terraform-aws-modules/s3-bucket/aws version ~> 5.0 is not compatible with AWS provider version ~> 5.0. Use module version ~> 4.0 for AWS provider ~> 5.0",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 13, Column: 13},
						End:      hcl.Pos{Line: 13, Column: 21},
					},
				},
			},
		},
		{
			Name: "valid vpc module with AWS provider v5",
			Content: `
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "~> 5.0"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "valid vpc module with AWS provider v6",
			Content: `
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.0"
    }
  }
}

module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "~> 6.0"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "invalid vpc module version for AWS provider v5",
			Content: `
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "~> 6.0"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsModuleVersionRule(),
					Message: "Module terraform-aws-modules/vpc/aws version ~> 6.0 is not compatible with AWS provider version ~> 5.0. Use module version ~> 5.0 for AWS provider ~> 5.0",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 13, Column: 13},
						End:      hcl.Pos{Line: 13, Column: 21},
					},
				},
			},
		},
		{
			Name: "valid lambda module with AWS provider v5",
			Content: `
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

module "lambda" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "~> 7.0"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "valid lambda module with AWS provider v6",
			Content: `
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.0"
    }
  }
}

module "lambda" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "~> 8.0"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "invalid lambda module version for AWS provider v5",
			Content: `
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

module "lambda" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "~> 8.0"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsModuleVersionRule(),
					Message: "Module terraform-aws-modules/lambda/aws version ~> 8.0 is not compatible with AWS provider version ~> 5.0. Use module version ~> 7.0 for AWS provider ~> 5.0",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 13, Column: 13},
						End:      hcl.Pos{Line: 13, Column: 21},
					},
				},
			},
		},
		{
			Name: "module without version constraint",
			Content: `
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

module "s3_bucket" {
  source = "terraform-aws-modules/s3-bucket/aws"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsModuleVersionRule(),
					Message: "Module terraform-aws-modules/s3-bucket/aws should specify a version constraint",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 12, Column: 12},
						End:      hcl.Pos{Line: 12, Column: 49},
					},
				},
			},
		},
		{
			Name: "module with invalid version operator",
			Content: `
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

module "s3_bucket" {
  source  = "terraform-aws-modules/s3-bucket/aws"
  version = ">= 4.0"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsModuleVersionRule(),
					Message: "Module terraform-aws-modules/s3-bucket/aws version constraint should use '~> x.y' format where x is the major version and y is the minor version (no patch version), got: >= 4.0",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 13, Column: 13},
						End:      hcl.Pos{Line: 13, Column: 21},
					},
				},
			},
		},
		{
			Name: "module with patch version in constraint",
			Content: `
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

module "s3_bucket" {
  source  = "terraform-aws-modules/s3-bucket/aws"
  version = "~> 4.1.2"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsModuleVersionRule(),
					Message: "Module terraform-aws-modules/s3-bucket/aws version constraint should use '~> x.y' format where x is the major version and y is the minor version (no patch version), got: ~> 4.1.2",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 13, Column: 13},
						End:      hcl.Pos{Line: 13, Column: 23},
					},
				},
			},
		},
		{
			Name: "non-aws module should not be checked",
			Content: `
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

module "custom" {
  source  = "./modules/custom"
  version = "1.0.0"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "no AWS provider defined",
			Content: `
module "s3_bucket" {
  source  = "terraform-aws-modules/s3-bucket/aws"
  version = "~> 4.0"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "multiple modules with mixed validity",
			Content: `
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

module "s3_bucket" {
  source  = "terraform-aws-modules/s3-bucket/aws"
  version = "~> 4.0"
}

module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "~> 6.0"
}

module "lambda" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "~> 7.0"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsModuleVersionRule(),
					Message: "Module terraform-aws-modules/vpc/aws version ~> 6.0 is not compatible with AWS provider version ~> 5.0. Use module version ~> 5.0 for AWS provider ~> 5.0",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 18, Column: 13},
						End:      hcl.Pos{Line: 18, Column: 21},
					},
				},
			},
		},
	}

	rule := NewAwsModuleVersionRule()

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssues(t, tc.Expected, runner.Issues)
		})
	}
}

func Test_AwsModuleVersionRule_WithMultipleFiles(t *testing.T) {
	rule := NewAwsModuleVersionRule()

	content := map[string]string{
		"versions.tf": `
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}`,
		"main.tf": `
module "s3_bucket" {
  source  = "terraform-aws-modules/s3-bucket/aws"
  version = "~> 5.0"
}

module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "~> 5.0"
}`,
	}

	runner := helper.TestRunner(t, content)

	if err := rule.Check(runner); err != nil {
		t.Fatalf("Unexpected error occurred: %s", err)
	}

	expected := helper.Issues{
		{
			Rule:    rule,
			Message: "Module terraform-aws-modules/s3-bucket/aws version ~> 5.0 is not compatible with AWS provider version ~> 5.0. Use module version ~> 4.0 for AWS provider ~> 5.0",
			Range: hcl.Range{
				Filename: "main.tf",
				Start:    hcl.Pos{Line: 4, Column: 13},
				End:      hcl.Pos{Line: 4, Column: 21},
			},
		},
	}

	helper.AssertIssues(t, expected, runner.Issues)
}