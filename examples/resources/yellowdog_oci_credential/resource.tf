# An Oracle Cloud Infrastructure API signing key, used by compute source
# templates to provision instances and instance pools in the tenancy.
resource "yellowdog_keyring" "example" {
  name = "example-keyring"
}

resource "yellowdog_oci_credential" "example" {
  name         = "oci"
  description  = "API key for the production tenancy"
  keyring_name = yellowdog_keyring.example.name

  tenant_id   = "ocid1.tenancy.oc1..aaaaaaaaexampleexampleexampleexampleexampleexampleexample"
  user_id     = "ocid1.user.oc1..aaaaaaaaexampleexampleexampleexampleexampleexampleexample"
  fingerprint = "12:34:56:78:90:ab:cd:ef:12:34:56:78:90:ab:cd:ef"
  # The PEM private key, for example file("oci_api_key.pem").
  private_key = "replace-with-your-pem-private-key"
  passphrase  = "replace-with-your-private-key-passphrase"

  lifecycle {
    # Deleting a keyring deletes its credentials, so recreate the credential
    # whenever the keyring is replaced.
    replace_triggered_by = [yellowdog_keyring.example]
  }
}
