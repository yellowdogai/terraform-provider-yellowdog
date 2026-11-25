package keyring_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/yellowdogai/terraform-provider-yellowdog/internal/testutil"
)

func TestKeyringResource(t *testing.T) {
	const resourceName = "yellowdog_keyring.test"
	const hcl = `
variable "name"        { default = null }
variable "description" { default = null }

resource "yellowdog_keyring" "test" {
  name        = var.name
  description = var.description
}
`
	configVariables := config.Variables{
		"name":        config.StringVariable("test-keyring"),
		"description": config.StringVariable("test description"),
	}

	testutil.RunAllStepsPlanOnly(t, hcl, map[string]resource.TestStep{
		"plan create": {
			ConfigVariables:    configVariables,
			ExpectNonEmptyPlan: true,
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PostApplyPreRefresh: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
					plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("keyring_password")),
					plancheck.ExpectSensitiveValue(resourceName, tfjsonpath.New("keyring_password")),
				},
			},
		},
		"missing name": {
			ConfigVariables: testutil.CloneWithout(configVariables, "name"),
			ExpectError:     regexp.MustCompile(`value for the name\s+attribute`),
		},
	})
}

// Read-only attributes are rejected even when set to null, so this needs its own config.
func TestKeyringResourceReadOnlyAttributes(t *testing.T) {
	const hcl = `
resource "yellowdog_keyring" "test" {
  name             = "test-keyring"
  keyring_password = "some value"
}
`
	testutil.RunAllStepsPlanOnly(t, hcl, map[string]resource.TestStep{
		"keyring password not configurable": {
			ExpectError: regexp.MustCompile(`Read-Only Attribute[\s\S]*keyring_password = "some value"`),
		},
	})
}
