# An Azure storage account key. The platform checks the account and key are
# valid when the credential is written.
resource "yellowdog_keyring" "example" {
  name = "example-keyring"
}

resource "yellowdog_azure_storage_credential" "example" {
  name         = "azure-storage"
  description  = "Storage account for job data"
  keyring_name = yellowdog_keyring.example.name

  account_name = "examplestorageaccount"
  account_key  = "replace-with-your-storage-account-key"

  lifecycle {
    # Deleting a keyring deletes its credentials, so recreate the credential
    # whenever the keyring is replaced.
    replace_triggered_by = [yellowdog_keyring.example]
  }
}
