package compute

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/yellowdogai/terraform-provider-yellowdog/internal/convert"
	"github.com/yellowdogai/yellowdog-sdk-go/model"
)

func (r *computeRequirementStaticTemplateResource) create(ctx context.Context, state *computeRequirementStaticTemplateModel, diags *diag.Diagnostics) (bool, error) {
	template := state.toSDK(ctx, diags)
	computeRequirementTemplate, err := r.client.Compute.AddComputeRequirementTemplate(ctx, template)
	if err != nil {
		return NotCreated, err
	}
	computeRequirementStaticTemplate, ok := computeRequirementTemplate.(*model.ComputeRequirementStaticTemplate)
	if !ok {
		state.Id = types.StringValue(computeRequirementTemplate.GetId())
		return Created, fmt.Errorf("expected a static compute requirement template, got %T", computeRequirementTemplate)
	}
	state.updateFromSDK(ctx, computeRequirementStaticTemplate, diags)
	return Created, nil
}

func (m *computeRequirementStaticTemplateModel) toSDK(ctx context.Context, diags *diag.Diagnostics) *model.ComputeRequirementStaticTemplate {
	template := model.ComputeRequirementStaticTemplate{}
	template.Id = convert.TFToOptionalString(m.Id)
	template.Name = m.Name.ValueString()
	template.Namespace = convert.TFToOptionalString(m.Namespace)
	template.Description = convert.TFToOptionalString(m.Description)
	template.StrategyType = m.StrategyType.ValueString()
	template.ImagesId = convert.TFToOptionalString(m.ImagesId)
	template.UserData = convert.TFToOptionalString(m.UserData)
	template.InstanceTags = convert.TFToStringMap(m.InstanceTags)
	template.Sources = convert.TFListToSlice(ctx, m.Sources, (*ComputeSourceUsageModel).toSDK, diags)
	return &template
}

func (m *ComputeSourceUsageModel) toSDK() model.ComputeSourceUsage {
	return model.ComputeSourceUsage{
		SourceTemplateId: m.SourceTemplateId.ValueString(),
		InstanceType:     convert.TFToOptionalString(m.InstanceType),
		ImageId:          convert.TFToOptionalString(m.ImageId),
	}
}

func (r *computeRequirementStaticTemplateResource) read(ctx context.Context, state *computeRequirementStaticTemplateModel, diagnostics *diag.Diagnostics) error {
	computeRequirementTemplate, err := r.client.Compute.GetComputeRequirementTemplate(ctx, state.Id.ValueString())
	if err != nil {
		return err
	}
	computeRequirementStaticTemplate, ok := computeRequirementTemplate.(*model.ComputeRequirementStaticTemplate)
	if !ok {
		return fmt.Errorf("expected a static compute requirement template, got %T", computeRequirementTemplate)
	}
	state.updateFromSDK(ctx, computeRequirementStaticTemplate, diagnostics)
	return nil
}

func (r *computeRequirementStaticTemplateResource) update(ctx context.Context, state *computeRequirementStaticTemplateModel, diags *diag.Diagnostics) error {
	template := state.toSDK(ctx, diags)
	computeRequirementTemplate, err := r.client.Compute.UpdateComputeRequirementTemplate(ctx, template)
	if err != nil {
		return err
	}
	computeRequirementStaticTemplate, ok := computeRequirementTemplate.(*model.ComputeRequirementStaticTemplate)
	if !ok {
		return fmt.Errorf("expected a static compute requirement template, got %T", computeRequirementTemplate)
	}
	state.updateFromSDK(ctx, computeRequirementStaticTemplate, diags)
	return nil
}

func (r *computeRequirementStaticTemplateResource) delete(ctx context.Context, state computeRequirementStaticTemplateModel, _ *diag.Diagnostics) error {
	return r.client.Compute.DeleteComputeRequirementTemplateById(ctx, state.Id.ValueString())
}

func (m *computeRequirementStaticTemplateModel) updateFromSDK(ctx context.Context, template *model.ComputeRequirementStaticTemplate, diags *diag.Diagnostics) {
	m.Id = convert.OptionalStringToTF(template.Id)
	m.Name = types.StringValue(template.Name)
	m.Namespace = convert.OptionalStringToTF(template.Namespace)
	m.Description = convert.OptionalStringToTF(template.Description)
	m.StrategyType = types.StringValue(template.StrategyType)
	m.ImagesId = convert.OptionalStringToTF(template.ImagesId)
	m.UserData = convert.OptionalStringToTF(template.UserData)
	m.InstanceTags = convert.StringMapToTF(template.InstanceTags)
	m.Sources = convert.SliceToTFList(ctx, template.Sources, createComputeSourceUsageModelFromSDK, computeSourceUsageObjectType(), diags)
}

func createComputeSourceUsageModelFromSDK(source *model.ComputeSourceUsage, _ *diag.Diagnostics) ComputeSourceUsageModel {
	return ComputeSourceUsageModel{
		SourceTemplateId: types.StringValue(source.SourceTemplateId),
		InstanceType:     convert.OptionalStringToTF(source.InstanceType),
		ImageId:          convert.OptionalStringToTF(source.ImageId),
	}
}

func computeSourceUsageObjectType() types.ObjectType {
	return schema.NestedAttributeObject{Attributes: computeSourceUsageSchema()}.Type().(types.ObjectType)
}
