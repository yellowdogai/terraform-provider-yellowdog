package account_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/yellowdogai/terraform-provider-yellowdog/internal/testutil"
)

func TestGroupRoleResource(t *testing.T) {
	const resourceName = "yellowdog_group_role.test"
	const hcl = `
variable "role_id"    { default = null }
variable "group_id"   { default = null }
variable "namespaces" { default = null }

resource "yellowdog_group_role" "test" {
  role_id    = var.role_id
  group_id   = var.group_id
  namespaces = var.namespaces
}
`
	roleId := "ydid:role:00000000-0000-0000-0000-000000000000"
	groupId := "ydid:group:00000000-0000-0000-0000-000000000000"
	namespace := "test-namespace"
	configVariables := config.Variables{
		"role_id":    config.StringVariable(roleId),
		"group_id":   config.StringVariable(groupId),
		"namespaces": config.ListVariable(config.StringVariable(namespace)),
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
		"invalid role id": {
			ConfigVariables: testutil.CloneWith(configVariables, "role_id", config.StringVariable("not a ydid")),
			ExpectError:     regexp.MustCompile("must be a role ydid"),
		},
		"invalid group id": {
			ConfigVariables: testutil.CloneWith(configVariables, "group_id", config.StringVariable("not a ydid")),
			ExpectError:     regexp.MustCompile("must be a group ydid"),
		},
		"missing role id": {
			ConfigVariables: testutil.CloneWithout(configVariables, "role_id"),
			ExpectError:     regexp.MustCompile("role_id"),
		},
		"missing group id": {
			ConfigVariables: testutil.CloneWithout(configVariables, "group_id"),
			ExpectError:     regexp.MustCompile("group_id"),
		},
	})
}
