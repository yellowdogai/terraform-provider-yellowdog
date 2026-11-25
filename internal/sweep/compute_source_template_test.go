package sweep_test

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	yellowdog "github.com/yellowdogai/yellowdog-sdk-go"
	"github.com/yellowdogai/yellowdog-sdk-go/compute"
	"github.com/yellowdogai/yellowdog-sdk-go/optional"
)

// sweepComputeSourceTemplates deletes the computeSourceTemplates the acceptance tests left behind
func sweepComputeSourceTemplates() resource.SweeperFunc {
	return createSweeperFunc(
		func(ctx context.Context, client *yellowdog.PlatformClient) ([]compute.ComputeSourceTemplateSummary, error) {
			return client.Compute.GetComputeSourceTemplates(&compute.ComputeSourceTemplateSearch{
				Name: optional.Of(sweepPrefix),
			}).ListAll(ctx)
		},
		func(v compute.ComputeSourceTemplateSummary) string {
			return v.Name.Or("")
		},
		func(ctx context.Context, client *yellowdog.PlatformClient, v compute.ComputeSourceTemplateSummary) error {
			return client.Compute.DeleteComputeSourceTemplateById(ctx, v.Id.Or(""))
		},
	)
}
