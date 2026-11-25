package image

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/yellowdogai/terraform-provider-yellowdog/internal/convert"
	"github.com/yellowdogai/terraform-provider-yellowdog/internal/resourceerror"
	"github.com/yellowdogai/yellowdog-sdk-go/model"
)

func (r *machineImageResource) create(ctx context.Context, state *machineImageModel, diags *diag.Diagnostics) (bool, error) {
	imageGroup, err := r.client.Images.GetImageGroup(ctx, state.GroupId.ValueString())
	if err != nil {
		if resourceerror.IsNotFound(err) {
			return NotCreated, fmt.Errorf("create image in group %s: %w", state.GroupId.ValueString(), err)
		}
		return NotCreated, err
	}
	machineImage, err := r.client.Images.AddImage(ctx, imageGroup, state.toSDK(ctx, diags))
	if err != nil {
		return NotCreated, err
	}
	state.updateFromSDK(ctx, machineImage, diags)
	return Created, nil
}

func (r *machineImageResource) read(ctx context.Context, state *machineImageModel, diags *diag.Diagnostics) error {
	machineImage, err := r.client.Images.GetImage(ctx, state.Id.ValueString())
	if err != nil {
		return err
	}
	if state.GroupId.IsNull() {
		groupId, err := findGroupIdOfImage(ctx, r.client, state.Id.ValueString())
		if err != nil {
			return err
		}
		state.GroupId = types.StringValue(groupId)
	}
	state.updateFromSDK(ctx, machineImage, diags)
	return nil
}

func (r *machineImageResource) update(ctx context.Context, state *machineImageModel, diags *diag.Diagnostics) error {
	machineImage, err := r.client.Images.UpdateImage(ctx, state.toSDK(ctx, diags))
	if err != nil {
		return err
	}
	state.updateFromSDK(ctx, machineImage, diags)
	return nil
}

func (r *machineImageResource) delete(ctx context.Context, state machineImageModel, _ *diag.Diagnostics) error {
	machineImage, err := r.client.Images.GetImage(ctx, state.Id.ValueString())
	if err != nil {
		if resourceerror.IsNotFound(err) {
			// already deleted
			return nil
		}
		return err
	}
	return r.client.Images.DeleteImage(ctx, machineImage)
}

func (m *machineImageModel) toSDK(_ context.Context, _ *diag.Diagnostics) *model.MachineImage {
	return &model.MachineImage{
		Id:                     convert.TFToOptionalString(m.Id),
		Name:                   m.Name.ValueString(),
		Provider:               convert.TFToEnum[model.CloudProvider](m.CloudProvider),
		Regions:                convert.TFListToStringSlice(m.Regions),
		ProviderImageId:        m.ProviderImageId.ValueString(),
		OsType:                 convert.TFToEnum[model.ImageOsType](m.OsType),
		SupportedInstanceTypes: convert.TFListToStringSlice(m.SupportedInstanceTypes),
		Metadata:               convert.TFToStringMap(m.Metadata),
	}
}

func (m *machineImageModel) updateFromSDK(_ context.Context, image *model.MachineImage, _ *diag.Diagnostics) {
	m.Id = convert.OptionalStringToTF(image.Id)
	m.Name = types.StringValue(image.Name)
	m.CloudProvider = convert.EnumToTF[model.CloudProvider](image.Provider)
	m.Regions = convert.StringSliceToTFList(image.Regions)
	m.ProviderImageId = types.StringValue(image.ProviderImageId)
	m.OsType = convert.EnumToTF[model.ImageOsType](image.OsType)
	m.SupportedInstanceTypes = convert.StringSliceToTFList(image.SupportedInstanceTypes)
	m.Metadata = convert.StringMapToTF(image.Metadata)
}
