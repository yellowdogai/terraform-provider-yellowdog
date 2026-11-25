# A keyring holds the credentials -- cloud provider accounts, instance logins,
# storage keys -- that compute source templates refer to as
# "<keyring name>/<credential name>".
resource "yellowdog_keyring" "example" {
  name        = "example-keyring"
  description = "Cloud provider credentials managed by Terraform"
}

# The keyring password is returned only once, when the keyring is created.
output "example_keyring_password" {
  value     = yellowdog_keyring.example.keyring_password
  sensitive = true
}
