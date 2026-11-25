package sweep_test

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	yellowdog "github.com/yellowdogai/yellowdog-sdk-go"
	"github.com/yellowdogai/yellowdog-sdk-go/images"
	"github.com/yellowdogai/yellowdog-sdk-go/model"
	"github.com/yellowdogai/yellowdog-sdk-go/optional"
)

// sweepImageFamilies deletes the machine image families the acceptance tests left behind, which deletes the groups
// and images in them too
func sweepImageFamilies() resource.SweeperFunc {
	return createSweeperFunc(
		func(ctx context.Context, client *yellowdog.PlatformClient) ([]images.MachineImageFamilySummary, error) {
			return client.Images.GetImageFamilies(&images.MachineImageFamilySearch{
				FamilyName: optional.Of(sweepPrefix),
			}).ListAll(ctx)
		},
		func(v images.MachineImageFamilySummary) string {
			return v.Name.Or("")
		},
		func(ctx context.Context, client *yellowdog.PlatformClient, v images.MachineImageFamilySummary) error {
			return client.Images.DeleteImageFamily(ctx, &model.MachineImageFamily{Id: v.Id})
		},
	)
}
