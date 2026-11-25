# A static compute requirement template names the compute source templates to
# provision from. strategy_type decides how instances are spread across them:
# SingleSourceProvisionStrategy, SplitProvisionStrategy or
# WaterfallProvisionStrategy.
resource "yellowdog_compute_source_template" "example" {
  namespace = "example-namespace"

  source = {
    simulator_compute_source = {
      name          = "simulator"
      region        = "sim-region-1"
      instance_type = "sim-small"
      image_id      = "sim-image-1"
    }
  }
}

resource "yellowdog_compute_requirement_static_template" "example" {
  name          = "example-static-template"
  namespace     = "example-namespace"
  description   = "Managed by Terraform"
  strategy_type = "co.yellowdog.platform.model.WaterfallProvisionStrategy"

  sources = [
    {
      source_template_id = yellowdog_compute_source_template.example.id
      # Overrides the source template's instance type for this requirement.
      instance_type = "sim-large"
    }
  ]

  instance_tags = {
    team = "research"
  }
}
