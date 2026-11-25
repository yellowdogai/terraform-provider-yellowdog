package image

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/yellowdogai/terraform-provider-yellowdog/internal/resourceerror"
	yellowdog "github.com/yellowdogai/yellowdog-sdk-go"
	"github.com/yellowdogai/yellowdog-sdk-go/images"
	"github.com/yellowdogai/yellowdog-sdk-go/model"
)

// A group or image does not say which family or group it belongs to, so after an import its family_id or group_id is
// unknown. Left null, the Required and RequiresReplace parent id would make the next plan replace the resource, so it is
// recovered by searching the families, each of which holds its groups and their images. This is only done when the id
// is missing, i.e. on import.

// findFamilyIdOfGroup returns the id of the family holding the group.
func findFamilyIdOfGroup(ctx context.Context, client *yellowdog.PlatformClient, groupId string) (string, error) {
	var familyId string
	err := forEachFamily(ctx, client, func(family *model.MachineImageFamily) bool {
		for _, group := range family.ImageGroups {
			if group.Id.Or("") == groupId {
				familyId = family.Id.Or("")
				return false
			}
		}
		return true
	})
	if err == nil && familyId == "" {
		err = fmt.Errorf("no machine image family holds machine image group %s", groupId)
	}
	return familyId, err
}

// findGroupIdOfImage returns the id of the group holding the image.
func findGroupIdOfImage(ctx context.Context, client *yellowdog.PlatformClient, imageId string) (string, error) {
	var groupId string
	err := forEachFamily(ctx, client, func(family *model.MachineImageFamily) bool {
		for _, group := range family.ImageGroups {
			for _, image := range group.Images {
				if image.Id.Or("") == imageId {
					groupId = group.Id.Or("")
					return false
				}
			}
		}
		return true
	})
	if err == nil && groupId == "" {
		err = fmt.Errorf("no machine image group holds machine image %s", imageId)
	}
	return groupId, err
}

// forEachFamily calls f with each of the account's families, with their groups and images, until f returns false.
func forEachFamily(ctx context.Context, client *yellowdog.PlatformClient, f func(*model.MachineImageFamily) bool) error {
	for summary, err := range client.Images.GetImageFamilies(&images.MachineImageFamilySearch{}).Iterate(ctx) {
		if err != nil {
			return fmt.Errorf("searching machine image families: %w", err)
		}
		family, err := client.Images.GetImageFamily(ctx, summary.Id.Or(""))
		if resourceerror.IsNotFound(err) {
			// deleted since it was listed
			// can happen if changes are being made concurrently, e.g. during parallel testing
			tflog.Info(ctx, "forEachFamily image family removed during for each")
			continue
		}
		if err != nil {
			return err
		}
		if !f(family) {
			return nil
		}
	}
	return nil
}
