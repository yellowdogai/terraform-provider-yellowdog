package keyring_test

import (
	"fmt"
	"maps"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/yellowdogai/terraform-provider-yellowdog/internal/testutil"
)

// TestKeyringCredentialResources covers every keyring credential type, as they share their keyring_name handling.
func TestKeyringCredentialResources(t *testing.T) {
	for resourceType, requiredAttributes := range map[string][]string{
		"yellowdog_aws_account_role_credential": {"external_role_arn", "external_id"},
		"yellowdog_aws_credential":              {"access_key_id", "secret_access_key"},
		"yellowdog_azure_client_credential":     {"client_id", "tenant_id", "subscription_id", "key"},
		"yellowdog_azure_instance_credential":   {"admin_username"},
		"yellowdog_azure_storage_credential":    {"account_name", "account_key"},
		"yellowdog_google_cloud_credential":     {"service_account_key_json"},
		"yellowdog_oci_credential":              {"user_id", "tenant_id", "fingerprint", "private_key"},
	} {
		t.Run(strings.TrimPrefix(resourceType, "yellowdog_"), func(t *testing.T) {
			resourceName := resourceType + ".test"
			attributes := append([]string{"keyring_name", "name"}, requiredAttributes...)

			var hcl strings.Builder
			hcl.WriteString("\n")
			configVariables := config.Variables{}
			for _, attribute := range attributes {
				fmt.Fprintf(&hcl, "variable %q { default = null }\n", attribute)
				configVariables[attribute] = config.StringVariable("test-" + attribute)
			}
			fmt.Fprintf(&hcl, "\nresource %q \"test\" {\n", resourceType)
			for _, attribute := range attributes {
				fmt.Fprintf(&hcl, "  %s = var.%s\n", attribute, attribute)
			}
			hcl.WriteString("}\n")

			steps := map[string]resource.TestStep{
				"plan create": {
					ConfigVariables:    configVariables,
					ExpectNonEmptyPlan: true,
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PostApplyPreRefresh: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
						},
					},
				},
			}
			for _, attribute := range slices.Sorted(maps.Keys(configVariables)) {
				steps["missing "+attribute] = resource.TestStep{
					ConfigVariables: testutil.CloneWithout(configVariables, attribute),
					ExpectError:     regexp.MustCompile(`value for the ` + attribute + `\s+attribute`),
				}
			}
			testutil.RunAllStepsPlanOnly(t, hcl.String(), steps)
		})
	}
}
