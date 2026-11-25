package account_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/yellowdogai/terraform-provider-yellowdog/internal/testutil"
)

func TestApplicationResource(t *testing.T) {
	const resourceName = "yellowdog_application.test"
	const hcl = `
variable "name"        { default = null }
variable "description" { default = null }

resource "yellowdog_application" "test" {
  name        = var.name
  description = var.description
}
`
	configVariables := config.Variables{
		"name":        config.StringVariable("test-application"),
		"description": config.StringVariable("test description"),
	}

	testutil.RunAllStepsPlanOnly(t, hcl, map[string]resource.TestStep{
		"plan create": {
			ConfigVariables:    configVariables,
			ExpectNonEmptyPlan: true,
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PostApplyPreRefresh: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
					plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("api_key")),
					plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("api_secret")),
					plancheck.ExpectSensitiveValue(resourceName, tfjsonpath.New("api_key")),
					plancheck.ExpectSensitiveValue(resourceName, tfjsonpath.New("api_secret")),
				},
			},
		},
		"missing name": {
			ConfigVariables: testutil.CloneWithout(configVariables, "name"),
			ExpectError:     regexp.MustCompile(`value for the name\s+attribute`),
		},
	})
}

// Read-only attributes are rejected even when set to null, so each needs its own config.
func TestApplicationResourceReadOnlyAttributes(t *testing.T) {
	for _, attribute := range []string{"api_key", "api_secret"} {
		t.Run(attribute, func(t *testing.T) {
			hcl := fmt.Sprintf(`
resource "yellowdog_application" "test" {
  name = "test-application"
  %s = "some value"
}
`, attribute)
			testutil.RunAllStepsPlanOnly(t, hcl, map[string]resource.TestStep{
				"not configurable": {
					ExpectError: regexp.MustCompile(`Read-Only Attribute[\s\S]*` + attribute + ` = "some value"`),
				},
			})
		})
	}
}
