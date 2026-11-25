package testutil

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// RunAllStepsPlanOnly run all steps as separate tests, in PlanOnly mode. This has the advantage that all steps run
// regardless of a failure in an earlier step.
func RunAllStepsPlanOnly(t *testing.T, hcl string, steps map[string]resource.TestStep) {
	t.Helper()
	for name, step := range steps {
		t.Run(name, func(t *testing.T) {
			step.Config = FakeYellowdogProviderHcl + hcl
			step.PlanOnly = true
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: ProtoV6ProviderFactories,
				Steps:                    []resource.TestStep{step},
			})
		})
	}
}
