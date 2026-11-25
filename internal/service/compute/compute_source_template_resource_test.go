package compute_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/yellowdogai/terraform-provider-yellowdog/internal/testutil"
)

func TestComputeSourceTemplateResource(t *testing.T) {
	const resourceName = "yellowdog_compute_source_template.test"
	const hcl = `
variable "namespace"      { default = null }
variable "compute_source" { default = null }

resource "yellowdog_compute_source_template" "test" {
  namespace = var.namespace
  source    = var.compute_source
}
`
	simulatorSource := config.ObjectVariable(map[string]config.Variable{
		"name":          config.StringVariable("test-simulator-source"),
		"region":        config.StringVariable("sim-region-1"),
		"instance_type": config.StringVariable("sim.small"),
		"image_id":      config.StringVariable("sim-image-1"),
	})
	awsInstancesSource := config.ObjectVariable(map[string]config.Variable{
		"name":              config.StringVariable("test-aws-source"),
		"credential":        config.StringVariable("test-keyring/test-credential"),
		"region":            config.StringVariable("eu-west-2"),
		"security_group_id": config.StringVariable("sg-00000000"),
		"instance_type":     config.StringVariable("t3.micro"),
		"image_id":          config.StringVariable("ami-00000000"),
	})
	gceInstancesSourceWithInvalidEnum := config.ObjectVariable(map[string]config.Variable{
		"name":                       config.StringVariable("test-gce-source"),
		"credential":                 config.StringVariable("test-keyring/test-credential"),
		"project":                    config.StringVariable("test-project"),
		"region":                     config.StringVariable("europe-west2"),
		"machine_type":               config.StringVariable("e2-micro"),
		"image":                      config.StringVariable("test-image"),
		"host_maintenance_behaviour": config.StringVariable("INVALID_ENUM_VALUE"),
	})
	configVariables := config.Variables{
		"namespace":      config.StringVariable("test-namespace"),
		"compute_source": config.ObjectVariable(map[string]config.Variable{"simulator_compute_source": simulatorSource}),
	}
	cloudConfigVariables := config.Variables{
		"namespace":      config.StringVariable("test-namespace"),
		"compute_source": config.ObjectVariable(map[string]config.Variable{"aws_instances_compute_source": awsInstancesSource}),
	}

	testutil.RunAllStepsPlanOnly(t, hcl, map[string]resource.TestStep{
		"plan create simulator": {
			ConfigVariables:    configVariables,
			ExpectNonEmptyPlan: true,
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PostApplyPreRefresh: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
				},
			},
		},
		"plan create cloud": {
			ConfigVariables:    cloudConfigVariables,
			ExpectNonEmptyPlan: true,
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PostApplyPreRefresh: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
				},
			},
		},
		// Smoke tests that the generated validators are wired up, not exhaustive coverage of every source or value.
		"no source": {
			ConfigVariables: testutil.CloneWith(configVariables, "compute_source", config.ObjectVariable(map[string]config.Variable{})),
			ExpectError:     regexp.MustCompile("No attribute specified"),
		},
		"two sources": {
			ConfigVariables: testutil.CloneWith(configVariables, "compute_source", config.ObjectVariable(map[string]config.Variable{
				"simulator_compute_source":     simulatorSource,
				"aws_instances_compute_source": awsInstancesSource,
			})),
			ExpectError: regexp.MustCompile("Invalid Attribute Combination"),
		},
		"invalid nested enum": {
			ConfigVariables: config.Variables{
				"namespace":      config.StringVariable("test-namespace"),
				"compute_source": config.ObjectVariable(map[string]config.Variable{"gce_instances_compute_source": gceInstancesSourceWithInvalidEnum}),
			},
			ExpectError: regexp.MustCompile(`value must\s+be one of`),
		},
	})
}
