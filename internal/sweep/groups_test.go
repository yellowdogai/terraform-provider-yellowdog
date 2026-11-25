package sweep_test

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	yellowdog "github.com/yellowdogai/yellowdog-sdk-go"
	"github.com/yellowdogai/yellowdog-sdk-go/account"
	"github.com/yellowdogai/yellowdog-sdk-go/optional"
)

// sweepGroups deletes the groups acceptance tests left behind
func sweepGroups() resource.SweeperFunc {
	return createSweeperFunc(
		func(ctx context.Context, client *yellowdog.PlatformClient) ([]account.GroupSummary, error) {
			return client.Account.GetGroups(&account.GroupSearch{
				Name: optional.Of(sweepPrefix),
			}).ListAll(ctx)
		},
		func(v account.GroupSummary) string {
			return v.Name.Or("")
		},
		func(ctx context.Context, client *yellowdog.PlatformClient, v account.GroupSummary) error {
			return client.Account.DeleteGroup(ctx, v.Id.Or(""))
		},
	)
}
