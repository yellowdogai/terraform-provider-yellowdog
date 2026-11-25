# An Azure service principal, used by compute source templates to provision
# virtual machines and scale sets in the subscription.
resource "yellowdog_keyring" "example" {
  name = "example-keyring"
}

resource "yellowdog_azure_client_credential" "example" {
  name         = "azure"
  description  = "Service principal for the production subscription"
  keyring_name = yellowdog_keyring.example.name

  client_id       = "00000000-0000-0000-0000-000000000000"
  tenant_id       = "00000000-0000-0000-0000-000000000000"
  subscription_id = "00000000-0000-0000-0000-000000000000"
  key             = "replace-with-your-client-secret"

  lifecycle {
    # Deleting a keyring deletes its credentials, so recreate the credential
    # whenever the keyring is replaced.
    replace_triggered_by = [yellowdog_keyring.example]
  }
}
