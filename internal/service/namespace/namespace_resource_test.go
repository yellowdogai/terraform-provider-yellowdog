package namespace_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/yellowdogai/terraform-provider-yellowdog/internal/testutil"
)

func TestNamespaceResource(t *testing.T) {
	const resourceName = "yellowdog_namespace.test"
	const hcl = `
variable "namespace" { default = null }

resource "yellowdog_namespace" "test" {
  namespace = var.namespace
}
`
	configVariables := config.Variables{
		"namespace": config.StringVariable("test-namespace"),
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
		"missing namespace": {
			ConfigVariables: testutil.CloneWithout(configVariables, "namespace"),
			ExpectError:     regexp.MustCompile(`value for the namespace\s+attribute`),
		},
	})
}
