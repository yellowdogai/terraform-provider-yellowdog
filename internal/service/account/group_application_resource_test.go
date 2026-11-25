package account_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/yellowdogai/terraform-provider-yellowdog/internal/testutil"
)

func TestGroupApplicationResource(t *testing.T) {
	const hcl = `
variable "group_id"       { default = null }
variable "application_id" { default = null }

resource "yellowdog_group_application" "test" {
  group_id       = var.group_id
  application_id = var.application_id
}
`
	groupId := "ydid:group:00000000-0000-0000-0000-000000000000"
	applicationId := "ydid:app:00000000-0000-0000-0000-000000000000"
	configVariables := config.Variables{
		"group_id":       config.StringVariable(groupId),
		"application_id": config.StringVariable(applicationId),
	}

	testutil.RunAllStepsPlanOnly(t, hcl, map[string]resource.TestStep{
		"invalid group id": {
			ConfigVariables: testutil.CloneWith(configVariables, "group_id", config.StringVariable("not a ydid")),
			ExpectError:     regexp.MustCompile("must be a group ydid"),
		},
		"invalid application id": {
			ConfigVariables: testutil.CloneWith(configVariables, "application_id", config.StringVariable("not a ydid")),
			ExpectError:     regexp.MustCompile("must be a app ydid"),
		},
		"missing group id": {
			ConfigVariables: testutil.CloneWithout(configVariables, "group_id"),
			ExpectError:     regexp.MustCompile("group_id"),
		},
		"missing application id": {
			ConfigVariables: testutil.CloneWithout(configVariables, "application_id"),
			ExpectError:     regexp.MustCompile("application_id"),
		},
	})
}
