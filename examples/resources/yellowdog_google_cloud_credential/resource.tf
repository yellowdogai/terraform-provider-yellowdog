# A Google Cloud service account key, used by compute source templates to
# provision Compute Engine instances in the project.
resource "yellowdog_keyring" "example" {
  name = "example-keyring"
}

resource "yellowdog_google_cloud_credential" "example" {
  name         = "gcp"
  description  = "Service account for the production project"
  keyring_name = yellowdog_keyring.example.name

  # The contents of the service account's JSON key file, for example
  # file("service-account-key.json").
  service_account_key_json = "replace-with-your-service-account-key-json"

  lifecycle {
    # Deleting a keyring deletes its credentials, so recreate the credential
    # whenever the keyring is replaced.
    replace_triggered_by = [yellowdog_keyring.example]
  }
}
