package image_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/yellowdogai/terraform-provider-yellowdog/internal/testutil"
)

func TestMachineImageResource(t *testing.T) {
	const resourceName = "yellowdog_machine_image.test"
	const hcl = `
variable "group_id"          { default = null }
variable "name"              { default = null }
variable "cloud_provider"    { default = null }
variable "provider_image_id" { default = null }
variable "os_type"           { default = null }

resource "yellowdog_machine_image" "test" {
  group_id          = var.group_id
  name              = var.name
  cloud_provider    = var.cloud_provider
  provider_image_id = var.provider_image_id
  os_type           = var.os_type
}
`
	configVariables := config.Variables{
		"group_id":          config.StringVariable("ydid:imggrp:00000000-0000-0000-0000-000000000000"),
		"name":              config.StringVariable("test-image"),
		"cloud_provider":    config.StringVariable("AWS"),
		"provider_image_id": config.StringVariable("ami-00000000"),
		"os_type":           config.StringVariable("LINUX"),
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
		"missing group id": {
			ConfigVariables: testutil.CloneWithout(configVariables, "group_id"),
			ExpectError:     regexp.MustCompile(`value for the group_id\s+attribute`),
		},
		"invalid group id": {
			ConfigVariables: testutil.CloneWith(configVariables, "group_id", config.StringVariable("not a ydid")),
			ExpectError:     regexp.MustCompile("must be a imggrp ydid"),
		},
		"missing name": {
			ConfigVariables: testutil.CloneWithout(configVariables, "name"),
			ExpectError:     regexp.MustCompile(`value for the name\s+attribute`),
		},
		"missing cloud provider": {
			ConfigVariables: testutil.CloneWithout(configVariables, "cloud_provider"),
			ExpectError:     regexp.MustCompile(`value for the cloud_provider\s+attribute`),
		},
		"missing provider image id": {
			ConfigVariables: testutil.CloneWithout(configVariables, "provider_image_id"),
			ExpectError:     regexp.MustCompile(`value for the provider_image_id\s+attribute`),
		},
		"missing os type": {
			ConfigVariables: testutil.CloneWithout(configVariables, "os_type"),
			ExpectError:     regexp.MustCompile(`value for the os_type\s+attribute`),
		},
		// Smoke tests that the generated enum validators are wired up, not exhaustive coverage of their values.
		"invalid cloud provider": {
			ConfigVariables: testutil.CloneWith(configVariables, "cloud_provider", config.StringVariable("INVALID_ENUM_VALUE")),
			ExpectError:     regexp.MustCompile(`value must\s+be one of`),
		},
		"invalid os type": {
			ConfigVariables: testutil.CloneWith(configVariables, "os_type", config.StringVariable("INVALID_ENUM_VALUE")),
			ExpectError:     regexp.MustCompile(`value must\s+be one of`),
		},
	})
}
