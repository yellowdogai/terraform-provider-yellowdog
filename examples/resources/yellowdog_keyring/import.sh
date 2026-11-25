# The import ID is the keyring's name or its YDID; the id recorded in state is
# always the YDID. keyring_password is null after import, as the API returns it
# only when the keyring is created.
terraform import yellowdog_keyring.example example-keyring
