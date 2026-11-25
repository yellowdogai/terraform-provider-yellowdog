# Grants a role to a group, either in the listed namespaces or, with namespaces
# omitted, globally. Any change replaces the grant.
resource "yellowdog_group" "example" {
  name = "example-group"
}

resource "yellowdog_namespace" "example" {
  namespace = "example-namespace"
}

resource "yellowdog_group_role" "example" {
  group_id = yellowdog_group.example.id
  # A role's YDID; this is the built-in compute-viewer role.
  role_id = "ydid:role:843a35e5-45e5-4df9-a44a-b7306be213f3"
  # Namespace names, not IDs.
  namespaces = [yellowdog_namespace.example.namespace]
}
