# An image group holds the equivalent machine images for each cloud provider
# and region, typically one version of an image family. Add images to it with
# yellowdog_machine_image.
resource "yellowdog_machine_image_family" "example" {
  name      = "example-image-family"
  namespace = "example-namespace"
  os_type   = "LINUX"
  access    = "PRIVATE"
}

resource "yellowdog_machine_image_group" "example" {
  family_id = yellowdog_machine_image_family.example.id
  name      = "v1"
  os_type   = "LINUX"
}
