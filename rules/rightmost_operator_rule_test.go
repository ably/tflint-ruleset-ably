package rules

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_RightmostOperatorRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "valid version constraint in required_providers",
			Content: `
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.0"
    }
  }
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "invalid version constraint with patch version",
			Content: `
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.1.2"
    }
  }
}`,
			Expected: helper.Issues{
				{
					Rule:     NewRightmostOperatorRule(),
					Message:  "Provider aws version constraint should use '~> x.y' format where x is the major version and y is the minor version (no patch version), got: ~> 4.1.2",
					Range:    hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 4, Column: 11},
						End:      hcl.Pos{Line: 7, Column: 6},
					},
				},
			},
		},
		{
			Name: "valid version constraint with minor version",
			Content: `
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.1"
    }
  }
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "invalid version constraint with >= operator",
			Content: `
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 4.0"
    }
  }
}`,
			Expected: helper.Issues{
				{
					Rule:     NewRightmostOperatorRule(),
					Message:  "Provider aws version constraint should use '~> x.y' format where x is the major version and y is the minor version (no patch version), got: >= 4.0",
					Range:    hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 4, Column: 11},
						End:      hcl.Pos{Line: 7, Column: 6},
					},
				},
			},
		},
		{
			Name: "invalid version constraint with exact version",
			Content: `
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "4.0.0"
    }
  }
}`,
			Expected: helper.Issues{
				{
					Rule:     NewRightmostOperatorRule(),
					Message:  "Provider aws version constraint should use '~> x.y' format where x is the major version and y is the minor version (no patch version), got: 4.0.0",
					Range:    hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 4, Column: 11},
						End:      hcl.Pos{Line: 7, Column: 6},
					},
				},
			},
		},
		{
			Name: "valid version constraint in provider block",
			Content: `
provider "aws" {
  version = "~> 4.0"
  region  = "us-east-1"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "invalid version constraint in provider block",
			Content: `
provider "aws" {
  version = ">= 4.0"
  region  = "us-east-1"
}`,
			Expected: helper.Issues{
				{
					Rule:     NewRightmostOperatorRule(),
					Message:  "Provider aws version constraint should use '~> x.y' format where x is the major version and y is the minor version (no patch version), got: >= 4.0",
					Range:    hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 13},
						End:      hcl.Pos{Line: 3, Column: 21},
					},
				},
			},
		},
		{
			Name: "multiple providers with mixed constraints",
			Content: `
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
    google = {
      source  = "hashicorp/google"
      version = ">= 3.0"
    }
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "3.10.0"
    }
  }
}`,
			Expected: helper.Issues{
				{
					Rule:     NewRightmostOperatorRule(),
					Message:  "Provider google version constraint should use '~> x.y' format where x is the major version and y is the minor version (no patch version), got: >= 3.0",
					Range:    hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 8, Column: 14},
						End:      hcl.Pos{Line: 11, Column: 6},
					},
				},
				{
					Rule:     NewRightmostOperatorRule(),
					Message:  "Provider azurerm version constraint should use '~> x.y' format where x is the major version and y is the minor version (no patch version), got: 3.10.0",
					Range:    hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 12, Column: 15},
						End:      hcl.Pos{Line: 15, Column: 6},
					},
				},
			},
		},
		{
			Name: "provider without version constraint",
			Content: `
terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
  }
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "valid constraint with spaces",
			Content: `
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~>  5.0"
    }
  }
}`,
			Expected: helper.Issues{},
		},
	}

	rule := NewRightmostOperatorRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}

func Test_RightmostOperatorRule_WithModules(t *testing.T) {
	rule := NewRightmostOperatorRule()

	content := map[string]string{
		"main.tf": `
module "example" {
  source = "./modules/example"
}

terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
  }
}`,
		"modules/example/main.tf": `
terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = ">= 3.0"
    }
  }
}`,
	}

	runner := helper.TestRunner(t, content)

	if err := rule.Check(runner); err != nil {
		t.Fatalf("Unexpected error occurred: %s", err)
	}

	expected := helper.Issues{
		{
			Rule:     rule,
			Message:  "Provider google version constraint should use '~> x.y' format where x is the major version and y is the minor version (no patch version), got: >= 3.0",
			Range:    hcl.Range{
				Filename: "modules/example/main.tf",
				Start:    hcl.Pos{Line: 4, Column: 14},
				End:      hcl.Pos{Line: 7, Column: 6},
			},
		},
	}

	helper.AssertIssues(t, expected, runner.Issues)
}