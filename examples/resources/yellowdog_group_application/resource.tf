# Makes an application a member of a group, so that it is granted the group's
# roles. Changing either ID replaces the membership.
resource "yellowdog_group" "example" {
  name = "example-group"
}

resource "yellowdog_application" "example" {
  name = "example-application"
}

resource "yellowdog_group_application" "example" {
  group_id       = yellowdog_group.example.id
  application_id = yellowdog_application.example.id
}
