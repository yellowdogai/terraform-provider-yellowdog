# An AWS access key, used by compute source templates to provision EC2
# instances in the account it belongs to. Credentials are written to a keyring
# and cannot be read back, so changes made outside Terraform are not detected.
resource "yellowdog_keyring" "example" {
  name = "example-keyring"
}

resource "yellowdog_aws_credential" "example" {
  name         = "aws"
  description  = "Access to the production AWS account"
  keyring_name = yellowdog_keyring.example.name

  access_key_id     = "AKIAIOSFODNN7EXAMPLE"
  secret_access_key = "replace-with-your-secret-access-key"

  lifecycle {
    # Deleting a keyring deletes its credentials, so recreate the credential
    # whenever the keyring is replaced.
    replace_triggered_by = [yellowdog_keyring.example]
  }
}
