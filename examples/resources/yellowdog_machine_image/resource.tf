# A machine image registers a cloud provider's image, such as an AWS AMI, in an
# image group, so that compute requirements can select it by family or group.
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

resource "yellowdog_machine_image" "example" {
  group_id          = yellowdog_machine_image_group.example.id
  name              = "ubuntu-22.04-eu-west-2"
  cloud_provider    = "AWS"
  os_type           = "LINUX"
  provider_image_id = "ami-09dbc7ce74870d573"
  regions           = ["eu-west-2"]

  supported_instance_types = ["t3a.micro", "t3a.small"]
  metadata = {
    "os.version" = "22.04"
  }
}
