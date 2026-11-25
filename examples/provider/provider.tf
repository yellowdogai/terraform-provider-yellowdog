terraform {
  required_providers {
    yellowdog = {
      source = "yellowdogai/yellowdog"
    }
  }
}

# api_key_id, api_key_secret and api_url fall back to the YD_API_KEY_ID,
# YD_API_KEY_SECRET and YD_API_URL environment variables when unset
# here. Never write the API key secret as a literal value. Set it in the
# environment, or pass it from a secrets manager such as Vault.
provider "yellowdog" {
  api_url = "https://api.yellowdog.ai"
  timeout = "60s"
}
