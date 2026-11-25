# Terraform Provider for YellowDog

Manage [YellowDog Platform](https://yellowdog.ai/) configuration as code: namespaces, keyrings and cloud
credentials, compute source and requirement templates, machine images, and groups, roles and applications.

The full documentation, with every resource and its arguments, is on the
[Terraform Registry](https://registry.terraform.io/providers/yellowdogai/yellowdog/latest/docs).

## Usage

Requires Terraform 1.0 or later.

```terraform
terraform {
  required_providers {
    yellowdog = {
      source  = "yellowdogai/yellowdog"
      version = "~> 0.1"
    }
  }
}

provider "yellowdog" {}

resource "yellowdog_namespace" "example" {
  namespace = "example"
}
```

The provider authenticates with a YellowDog [application](https://docs.yellowdog.ai/#/the-platform/applications)'s API key.
The simplest option is with environment variables:

```shell
export YD_API_KEY_ID=...
export YD_API_KEY_SECRET=...
```

Alternatively read it from a secrets manager, such as HashiCorp Vault, and pass it to the `provider` block:

```terraform
provider "yellowdog" {
  api_key_id     = ephemeral.vault_kv_secret_v2.yellowdog.data.api_key_id
  api_key_secret = ephemeral.vault_kv_secret_v2.yellowdog.data.api_key_secret
}
```

> [!WARNING]
> Never write your API key secret as a literal value in configuration files. Anyone who can read the file can use the
> key, and the file may be committed to source control.

| Argument         | Environment variable | Default                    |
|------------------|----------------------|----------------------------|
| `api_key_id`     | `YD_API_KEY_ID`      | Required                   |
| `api_key_secret` | `YD_API_KEY_SECRET`  | Required                   |
| `api_url`        | `YD_API_URL`         | `https://api.yellowdog.ai` |
| `timeout`        |                      | `30s`, per request         |

## License

[Apache License 2.0](LICENSE). Copyright 2026 YellowDog Limited.
