# An IAM role in another AWS account for YellowDog to assume, rather than an
# access key. The role's trust policy must name the external ID.
resource "yellowdog_keyring" "example" {
  name = "example-keyring"
}

resource "yellowdog_aws_account_role_credential" "example" {
  name         = "aws-role"
  description  = "Cross-account role in the production AWS account"
  keyring_name = yellowdog_keyring.example.name

  external_role_arn = "arn:aws:iam::123456789012:role/yellowdog-provisioner"
  external_id       = "replace-with-your-external-id"

  lifecycle {
    # Deleting a keyring deletes its credentials, so recreate the credential
    # whenever the keyring is replaced.
    replace_triggered_by = [yellowdog_keyring.example]
  }
}
