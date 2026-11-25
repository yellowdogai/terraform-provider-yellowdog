package compute_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/yellowdogai/terraform-provider-yellowdog/internal/testutil"
)

func TestComputeRequirementStaticTemplateResource(t *testing.T) {
	const resourceName = "yellowdog_compute_requirement_static_template.test"
	const hcl = `
variable "name"          { default = null }
variable "namespace"     { default = null }
variable "strategy_type" { default = null }
variable "sources"       { default = null }

resource "yellowdog_compute_requirement_static_template" "test" {
  name          = var.name
  namespace     = var.namespace
  strategy_type = var.strategy_type
  sources       = var.sources
}
`
	sources := func(sourceTemplateId string) config.Variable {
		return config.ListVariable(config.ObjectVariable(map[string]config.Variable{
			"source_template_id": config.StringVariable(sourceTemplateId),
		}))
	}
	configVariables := config.Variables{
		"name":          config.StringVariable("test-template"),
		"namespace":     config.StringVariable("test-namespace"),
		"strategy_type": config.StringVariable("co.yellowdog.platform.model.SplitProvisionStrategy"),
		"sources":       sources("ydid:cst:00000000-0000-0000-0000-000000000000"),
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
		"missing sources": {
			ConfigVariables: testutil.CloneWithout(configVariables, "sources"),
			ExpectError:     regexp.MustCompile(`value for the sources\s+attribute`),
		},
		// Smoke test that the generated ydid validator is wired up on a nested attribute.
		"invalid source template id": {
			ConfigVariables: testutil.CloneWith(configVariables, "sources", sources("not a ydid")),
			ExpectError:     regexp.MustCompile(`must be a YellowDog ID`),
		},
	})
}
