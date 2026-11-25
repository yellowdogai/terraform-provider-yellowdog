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

func (r *machineImageGroupResource) create(ctx context.Context, state *machineImageGroupModel, diags *diag.Diagnostics) (bool, error) {
	family, err := r.client.Images.GetImageFamily(ctx, state.FamilyId.ValueString())
	if err != nil {
		if resourceerror.IsNotFound(err) {
			return NotCreated, fmt.Errorf("create image group in family %s: %w", state.FamilyId.ValueString(), err)
		}
		return NotCreated, err
	}
	imageGroup, err := r.client.Images.AddImageGroup(ctx, family, state.toSDK(ctx, diags))
	if err != nil {
		return NotCreated, err
	}
	state.updateFromSDK(ctx, imageGroup, diags)
	return Created, nil
}

func (r *machineImageGroupResource) read(ctx context.Context, state *machineImageGroupModel, diags *diag.Diagnostics) error {
	group, err := r.client.Images.GetImageGroup(ctx, state.Id.ValueString())
	if err != nil {
		return err
	}
	if state.FamilyId.IsNull() {
		familyId, err := findFamilyIdOfGroup(ctx, r.client, state.Id.ValueString())
		if err != nil {
			return err
		}
		state.FamilyId = types.StringValue(familyId)
	}
	state.updateFromSDK(ctx, group, diags)
	return nil
}

func (r *machineImageGroupResource) update(ctx context.Context, state *machineImageGroupModel, diags *diag.Diagnostics) error {
	group, err := r.client.Images.UpdateImageGroup(ctx, state.toSDK(ctx, diags))
	if err != nil {
		return err
	}
	state.updateFromSDK(ctx, group, diags)
	return nil
}

func (r *machineImageGroupResource) delete(ctx context.Context, state machineImageGroupModel, diags *diag.Diagnostics) error {
	return r.client.Images.DeleteImageGroup(ctx, state.toSDK(ctx, diags))
}

func (m *machineImageGroupModel) toSDK(_ context.Context, _ *diag.Diagnostics) *model.MachineImageGroup {
	return &model.MachineImageGroup{
		Id:                    convert.TFToOptionalString(m.Id),
		Name:                  m.Name.ValueString(),
		OsType:                convert.TFToEnum[model.ImageOsType](m.OsType),
		MetadataSpecification: convert.TFToStringMap(m.MetadataSpecification),
	}
}

func (m *machineImageGroupModel) updateFromSDK(_ context.Context, group *model.MachineImageGroup, _ *diag.Diagnostics) {
	m.Id = convert.OptionalStringToTF(group.Id)
	m.Name = types.StringValue(group.Name)
	m.OsType = convert.EnumToTF(group.OsType)
	m.MetadataSpecification = convert.StringMapToTF(group.MetadataSpecification)
}
