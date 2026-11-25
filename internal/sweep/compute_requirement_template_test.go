package sweep_test

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	yellowdog "github.com/yellowdogai/yellowdog-sdk-go"
	"github.com/yellowdogai/yellowdog-sdk-go/compute"
	"github.com/yellowdogai/yellowdog-sdk-go/optional"
)

// sweepComputeRequirementTemplates deletes the computeRequirementTemplates the acceptance tests left behind
func sweepComputeRequirementTemplates() resource.SweeperFunc {
	return createSweeperFunc(
		func(ctx context.Context, client *yellowdog.PlatformClient) ([]compute.ComputeRequirementTemplateSummary, error) {
			return client.Compute.GetComputeRequirementTemplates(&compute.ComputeRequirementTemplateSearch{
				Name: optional.Of(sweepPrefix),
			}).ListAll(ctx)
		},
		func(v compute.ComputeRequirementTemplateSummary) string {
			return v.Name.Or("")
		},
		func(ctx context.Context, client *yellowdog.PlatformClient, v compute.ComputeRequirementTemplateSummary) error {
			return client.Compute.DeleteComputeRequirementTemplateById(ctx, v.Id.Or(""))
		},
	)
}
