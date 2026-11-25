# An application is a non-human identity with its own API key and secret, for
# programs that call the YellowDog API.
resource "yellowdog_application" "example" {
  name        = "example-application"
  description = "Managed by Terraform"
}

# api_key and api_secret are returned only once, when the application is
# created. Store them somewhere durable, such as a secret manager, if you need
# them later -- `terraform import` cannot recover them.
output "example_api_key" {
  value     = yellowdog_application.example.api_key
  sensitive = true
}

output "example_api_secret" {
  value     = yellowdog_application.example.api_secret
  sensitive = true
}
