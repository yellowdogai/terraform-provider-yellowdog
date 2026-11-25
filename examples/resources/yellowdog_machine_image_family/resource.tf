# An image family is the top level of the machine image hierarchy: a family
# holds groups, typically one per version, and a group holds the equivalent
# image in each cloud provider and region.
resource "yellowdog_machine_image_family" "example" {
  name      = "example-image-family"
  namespace = "example-namespace"
  os_type   = "LINUX"
  access    = "PRIVATE"
}
