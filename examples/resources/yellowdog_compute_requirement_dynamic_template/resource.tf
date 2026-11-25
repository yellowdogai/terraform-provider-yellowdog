# A dynamic compute requirement template chooses its compute source templates
# when a requirement is created from it: constraints decide which sources are
# eligible, and preferences rank the ones that are. Attributes are the
# platform's own, such as yd.cpu, yd.ram, source.provider and source.region,
# or user attributes set on compute source templates.
resource "yellowdog_compute_requirement_dynamic_template" "example" {
  name          = "example-dynamic-template"
  namespace     = "example-namespace"
  description   = "Managed by Terraform"
  strategy_type = "co.yellowdog.platform.model.SplitProvisionStrategy"

  minimum_source_count = 1
  maximum_source_count = 4

  # Only consider source templates in these namespaces.
  source_namespaces = ["example-namespace"]

  # An image family, group or image YDID to provision with.
  images_id = "ydid:imgfam:000000:00000000-0000-0000-0000-000000000000"

  constraints = [
    {
      string_attribute_constraint = {
        attribute = "source.provider"
        any_of    = ["AWS", "GOOGLE"]
      }
    },
    {
      numeric_attribute_constraint = {
        attribute = "yd.cpu"
        min       = 2
        max       = 16
      }
    },
  ]

  preferences = [
    {
      string_attribute_preference = {
        attribute        = "source.region"
        weight           = 1
        preferred_values = ["eu-west-2", "europe-west2"]
      }
    },
    {
      numeric_attribute_preference = {
        attribute  = "yd.ram"
        weight     = 2
        rank_order = "PREFER_HIGHER"
      }
    },
  ]

  source_traits = {
    can_stop_start = true
  }

  instance_tags = {
    team = "research"
  }
}
