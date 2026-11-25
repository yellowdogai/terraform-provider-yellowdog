package sweep_test

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	yellowdog "github.com/yellowdogai/yellowdog-sdk-go"
	"github.com/yellowdogai/yellowdog-sdk-go/account"
	"github.com/yellowdogai/yellowdog-sdk-go/model"
	"github.com/yellowdogai/yellowdog-sdk-go/optional"
)

// sweepNamespaces deletes the namespaces acceptance tests left behind
func sweepNamespaces() resource.SweeperFunc {
	return createSweeperFunc(
		func(ctx context.Context, client *yellowdog.PlatformClient) ([]model.Namespace, error) {
			return client.Namespaces.GetNamespaces(&account.NamespaceSearch{
				NamespaceContains: optional.Of(sweepPrefix),
				Deletable:         optional.Of(true),
			}).ListAll(ctx)
		},
		func(v model.Namespace) string {
			return v.Namespace
		},
		func(ctx context.Context, client *yellowdog.PlatformClient, v model.Namespace) error {
			return client.Namespaces.DeleteNamespace(ctx, v.Id)
		},
	)
}
