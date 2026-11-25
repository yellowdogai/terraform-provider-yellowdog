# A namespace groups compute, work and image resources, and is the unit that
# group roles can be scoped to. The platform will not delete a namespace once
# anything has been created in it.
resource "yellowdog_namespace" "example" {
  namespace = "example-namespace"
}
