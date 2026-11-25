package image

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/yellowdogai/terraform-provider-yellowdog/internal/convert"
	"github.com/yellowdogai/yellowdog-sdk-go/model"
)

func (r *machineImageFamilyResource) create(ctx context.Context, state *machineImageFamilyModel, diags *diag.Diagnostics) (bool, error) {
	imageFamily, err := r.client.Images.AddImageFamily(ctx, state.toSDK(ctx, diags))
	if err != nil {
		return NotCreated, err
	}
	state.updateFromSDK(ctx, imageFamily, diags)
	return Created, nil
}

func (r *machineImageFamilyResource) read(ctx context.Context, state *machineImageFamilyModel, diags *diag.Diagnostics) error {
	imageFamily, err := r.client.Images.GetImageFamily(ctx, state.Id.ValueString())
	if err != nil {
		return err
	}
	state.updateFromSDK(ctx, imageFamily, diags)
	return nil
}

func (r *machineImageFamilyResource) update(ctx context.Context, state *machineImageFamilyModel, diags *diag.Diagnostics) error {
	imageFamily, err := r.client.Images.UpdateImageFamily(ctx, state.toSDK(ctx, diags))
	if err != nil {
		return err
	}
	state.updateFromSDK(ctx, imageFamily, diags)
	return nil
}

func (r *machineImageFamilyResource) delete(ctx context.Context, state machineImageFamilyModel, diags *diag.Diagnostics) error {
	return r.client.Images.DeleteImageFamily(ctx, state.toSDK(ctx, diags))
}

func (m *machineImageFamilyModel) toSDK(_ context.Context, _ *diag.Diagnostics) *model.MachineImageFamily {
	return &model.MachineImageFamily{
		Id:                    convert.TFToOptionalString(m.Id),
		Namespace:             m.Namespace.ValueString(),
		Name:                  m.Name.ValueString(),
		OsType:                convert.TFToEnum[model.ImageOsType](m.OsType),
		Access:                convert.TFToEnum[model.ImageAccess](m.Access),
		MetadataSpecification: convert.TFToStringMap(m.MetadataSpecification),
	}
}

func (m *machineImageFamilyModel) updateFromSDK(_ context.Context, family *model.MachineImageFamily, _ *diag.Diagnostics) {
	m.Id = convert.OptionalStringToTF(family.Id)
	m.Namespace = types.StringValue(family.Namespace)
	m.Name = types.StringValue(family.Name)
	m.OsType = convert.EnumToTF(family.OsType)
	m.Access = convert.EnumToTF(family.Access)
	m.MetadataSpecification = convert.StringMapToTF(family.MetadataSpecification)
}
