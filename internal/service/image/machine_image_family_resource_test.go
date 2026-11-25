package image_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/yellowdogai/terraform-provider-yellowdog/internal/testutil"
)

func TestMachineImageFamilyResource(t *testing.T) {
	const resourceName = "yellowdog_machine_image_family.test"
	const hcl = `
variable "namespace" { default = null }
variable "name"      { default = null }
variable "os_type"   { default = null }
variable "access"    { default = null }

resource "yellowdog_machine_image_family" "test" {
  namespace = var.namespace
  name      = var.name
  os_type   = var.os_type
  access    = var.access
}
`
	configVariables := config.Variables{
		"namespace": config.StringVariable("test-namespace"),
		"name":      config.StringVariable("test-family"),
		"os_type":   config.StringVariable("LINUX"),
		"access":    config.StringVariable("PRIVATE"),
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
		"missing name": {
			ConfigVariables: testutil.CloneWithout(configVariables, "name"),
			ExpectError:     regexp.MustCompile(`value for the name\s+attribute`),
		},
		"missing os type": {
			ConfigVariables: testutil.CloneWithout(configVariables, "os_type"),
			ExpectError:     regexp.MustCompile(`value for the os_type\s+attribute`),
		},
		"default access": {
			ConfigVariables:    testutil.CloneWithout(configVariables, "access"),
			ExpectNonEmptyPlan: true,
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PostApplyPreRefresh: []plancheck.PlanCheck{
					plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("access"), knownvalue.StringExact("PRIVATE")),
				},
			},
		},
		// Smoke tests that the generated enum validators are wired up, not exhaustive coverage of their values.
		"invalid os type": {
			ConfigVariables: testutil.CloneWith(configVariables, "os_type", config.StringVariable("INVALID_ENUM_VALUE")),
			ExpectError:     regexp.MustCompile(`value must\s+be one of`),
		},
		"invalid access": {
			ConfigVariables: testutil.CloneWith(configVariables, "access", config.StringVariable("INVALID_ENUM_VALUE")),
			ExpectError:     regexp.MustCompile(`value must\s+be one of`),
		},
	})
}
