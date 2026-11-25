package image_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/yellowdogai/terraform-provider-yellowdog/internal/testutil"
)

func TestMachineImageGroupResource(t *testing.T) {
	const resourceName = "yellowdog_machine_image_group.test"
	const hcl = `
variable "family_id" { default = null }
variable "name"      { default = null }
variable "os_type"   { default = null }

resource "yellowdog_machine_image_group" "test" {
  family_id = var.family_id
  name      = var.name
  os_type   = var.os_type
}
`
	configVariables := config.Variables{
		"family_id": config.StringVariable("ydid:imgfam:00000000-0000-0000-0000-000000000000"),
		"name":      config.StringVariable("test-group"),
		"os_type":   config.StringVariable("LINUX"),
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
		"missing family id": {
			ConfigVariables: testutil.CloneWithout(configVariables, "family_id"),
			ExpectError:     regexp.MustCompile(`value for the family_id\s+attribute`),
		},
		"invalid family id": {
			ConfigVariables: testutil.CloneWith(configVariables, "family_id", config.StringVariable("not a ydid")),
			ExpectError:     regexp.MustCompile("must be a imgfam ydid"),
		},
		"missing name": {
			ConfigVariables: testutil.CloneWithout(configVariables, "name"),
			ExpectError:     regexp.MustCompile(`value for the name\s+attribute`),
		},
		"missing os type": {
			ConfigVariables: testutil.CloneWithout(configVariables, "os_type"),
			ExpectError:     regexp.MustCompile(`value for the os_type\s+attribute`),
		},
		// Smoke test that the generated enum validator is wired up, not exhaustive coverage of its values.
		"invalid os type": {
			ConfigVariables: testutil.CloneWith(configVariables, "os_type", config.StringVariable("INVALID_ENUM_VALUE")),
			ExpectError:     regexp.MustCompile(`value must\s+be one of`),
		},
	})
}
