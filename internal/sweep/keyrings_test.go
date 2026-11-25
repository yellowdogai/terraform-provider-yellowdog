package sweep_test

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	yellowdog "github.com/yellowdogai/yellowdog-sdk-go"
	"github.com/yellowdogai/yellowdog-sdk-go/account"
	"github.com/yellowdogai/yellowdog-sdk-go/optional"
)

// sweepKeyrings deletes the keyrings acceptance tests left behind
func sweepKeyrings() resource.SweeperFunc {
	return createSweeperFunc(
		func(ctx context.Context, client *yellowdog.PlatformClient) ([]account.KeyringSummary, error) {
			return client.Keyring.GetKeyrings(&account.KeyringSearch{
				Name: optional.Of(sweepPrefix),
			}).ListAll(ctx)
		},
		func(v account.KeyringSummary) string {
			return v.Name.Or("")
		},
		func(ctx context.Context, client *yellowdog.PlatformClient, v account.KeyringSummary) error {
			return client.Keyring.DeleteKeyringByName(ctx, v.Name.Or(""))
		},
	)
}
