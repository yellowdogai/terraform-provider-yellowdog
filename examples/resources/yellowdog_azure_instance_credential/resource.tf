# The administrator login for Azure virtual machines, referred to by the
# admin_user_credential of an Azure compute source.
resource "yellowdog_keyring" "example" {
  name = "example-keyring"
}

resource "yellowdog_azure_instance_credential" "example" {
  name         = "azure-admin"
  description  = "Administrator login for Azure instances"
  keyring_name = yellowdog_keyring.example.name

  admin_username = "yellowdog"
  admin_password = "replace-with-your-admin-password"

  lifecycle {
    # Deleting a keyring deletes its credentials, so recreate the credential
    # whenever the keyring is replaced.
    replace_triggered_by = [yellowdog_keyring.example]
  }
}
