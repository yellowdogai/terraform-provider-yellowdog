# A group gives its members -- users and applications -- the roles assigned to
# it. See yellowdog_group_role and yellowdog_group_application.
resource "yellowdog_group" "example" {
  name        = "example-group"
  description = "Managed by Terraform"
}
