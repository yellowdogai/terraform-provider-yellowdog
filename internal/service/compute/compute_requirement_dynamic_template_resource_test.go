package compute_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/yellowdogai/terraform-provider-yellowdog/internal/testutil"
)

func TestComputeRequirementDynamicTemplateResource(t *testing.T) {
	const resourceName = "yellowdog_compute_requirement_dynamic_template.test"
	const hcl = `
variable "name"          { default = null }
variable "namespace"     { default = null }
variable "strategy_type" { default = null }
variable "preferences"   { default = null }

resource "yellowdog_compute_requirement_dynamic_template" "test" {
  name          = var.name
  namespace     = var.namespace
  strategy_type = var.strategy_type
  preferences   = var.preferences
}
`
	preferences := func(rankOrder string) config.Variable {
		return config.ListVariable(config.ObjectVariable(map[string]config.Variable{
			"numeric_attribute_preference": config.ObjectVariable(map[string]config.Variable{
				"attribute":  config.StringVariable("yd.cost"),
				"weight":     config.FloatVariable(1.0),
				"rank_order": config.StringVariable(rankOrder),
			}),
		}))
	}
	configVariables := config.Variables{
		"name":          config.StringVariable("test-template"),
		"namespace":     config.StringVariable("test-namespace"),
		"strategy_type": config.StringVariable("co.yellowdog.platform.model.SplitProvisionStrategy"),
		"preferences":   preferences("PREFER_LOWER"),
	}

	testutil.RunAllStepsPlanOnly(t, hcl, map[string]resource.TestStep{
		"plan create": {
			ConfigVariables:    configVariables,
			ExpectNonEmptyPlan: true,
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PostApplyPreRefresh: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
				},
			},
		},
		"missing name": {
			ConfigVariables: testutil.CloneWithout(configVariables, "name"),
			ExpectError:     regexp.MustCompile(`value for the name\s+attribute`),
		},
		"missing strategy type": {
			ConfigVariables: testutil.CloneWithout(configVariables, "strategy_type"),
			ExpectError:     regexp.MustCompile(`value for the strategy_type\s+attribute`),
		},
		// Smoke test that the generated enum validators are wired up on nested attributes, not exhaustive coverage.
		"invalid rank order": {
			ConfigVariables: testutil.CloneWith(configVariables, "preferences", preferences("INVALID_ENUM_VALUE")),
			ExpectError:     regexp.MustCompile(`value must\s+be one of`),
		},
	})
}
