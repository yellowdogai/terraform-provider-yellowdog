package sweep_test

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	yellowdog "github.com/yellowdogai/yellowdog-sdk-go"
	"github.com/yellowdogai/yellowdog-sdk-go/account"
	"github.com/yellowdogai/yellowdog-sdk-go/model"
	"github.com/yellowdogai/yellowdog-sdk-go/optional"
)

// sweepApplications deletes the applications acceptance tests left behind
func sweepApplications() resource.SweeperFunc {
	return createSweeperFunc(
		func(ctx context.Context, client *yellowdog.PlatformClient) ([]model.Application, error) {
			return client.Account.GetApplications(&account.ApplicationSearch{
				Name: optional.Of(sweepPrefix),
			}).ListAll(ctx)
		},
		func(v model.Application) string {
			return v.Name
		},
		func(ctx context.Context, client *yellowdog.PlatformClient, v model.Application) error {
			return client.Account.DeleteApplication(ctx, v.Id.Or(""))
		},
	)
}
