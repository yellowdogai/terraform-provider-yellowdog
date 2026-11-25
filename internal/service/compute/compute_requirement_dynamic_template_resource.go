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

func (r *computeRequirementDynamicTemplateResource) create(ctx context.Context, state *computeRequirementDynamicTemplateModel, diags *diag.Diagnostics) (bool, error) {
	template := state.toSDK(ctx, diags)
	computeRequirementTemplate, err := r.client.Compute.AddComputeRequirementTemplate(ctx, template)
	if err != nil {
		return NotCreated, err
	}
	computeRequirementDynamicTemplate, ok := computeRequirementTemplate.(*model.ComputeRequirementDynamicTemplate)
	if !ok {
		state.Id = types.StringValue(computeRequirementTemplate.GetId())
		return Created, fmt.Errorf("expected a dynamic compute requirement template, got %T", computeRequirementTemplate)
	}
	state.updateFromSDK(ctx, computeRequirementDynamicTemplate, diags)
	return Created, nil
}

func (r *computeRequirementDynamicTemplateResource) read(ctx context.Context, state *computeRequirementDynamicTemplateModel, diagnostics *diag.Diagnostics) error {
	computeRequirementTemplate, err := r.client.Compute.GetComputeRequirementTemplate(ctx, state.Id.ValueString())
	if err != nil {
		return err
	}
	computeRequirementDynamicTemplate, ok := computeRequirementTemplate.(*model.ComputeRequirementDynamicTemplate)
	if !ok {
		return fmt.Errorf("expected a dynamic compute requirement template, got %T", computeRequirementTemplate)
	}
	state.updateFromSDK(ctx, computeRequirementDynamicTemplate, diagnostics)
	return nil
}

func (r *computeRequirementDynamicTemplateResource) update(ctx context.Context, state *computeRequirementDynamicTemplateModel, diags *diag.Diagnostics) error {

	template := state.toSDK(ctx, diags)
	computeRequirementTemplate, err := r.client.Compute.UpdateComputeRequirementTemplate(ctx, template)
	if err != nil {
		return err
	}
	computeRequirementDynamicTemplate, ok := computeRequirementTemplate.(*model.ComputeRequirementDynamicTemplate)
	if !ok {
		return fmt.Errorf("expected a dynamic compute requirement template, got %T", computeRequirementTemplate)
	}
	state.updateFromSDK(ctx, computeRequirementDynamicTemplate, diags)
	return nil
}

func (r *computeRequirementDynamicTemplateResource) delete(ctx context.Context, state computeRequirementDynamicTemplateModel, _ *diag.Diagnostics) error {
	return r.client.Compute.DeleteComputeRequirementTemplateById(ctx, state.Id.ValueString())
}

func (m *computeRequirementDynamicTemplateModel) toSDK(ctx context.Context, diags *diag.Diagnostics) *model.ComputeRequirementDynamicTemplate {
	template := model.ComputeRequirementDynamicTemplate{}
	template.Id = convert.TFToOptionalString(m.Id)
	template.Name = m.Name.ValueString()
	template.Namespace = convert.TFToOptionalString(m.Namespace)
	template.Description = convert.TFToOptionalString(m.Description)
	template.StrategyType = m.StrategyType.ValueString()
	template.MinimumSourceCount = convert.TFToOptionalInt32(m.MinimumSourceCount)
	template.MaximumSourceCount = convert.TFToOptionalInt32(m.MaximumSourceCount)
	template.ImagesId = convert.TFToOptionalString(m.ImagesId)
	template.UserData = convert.TFToOptionalString(m.UserData)
	template.InstanceTags = convert.TFToStringMap(m.InstanceTags)
	template.Constraints = convert.TFListToSlice(ctx, m.Constraints, (*AttributeConstraintModel).toSDK, diags)
	template.Preferences = convert.TFListToSlice(ctx, m.Preferences, (*AttributePreferenceModel).toSDK, diags)
	template.SourceTraits = convert.TFToOptionalObjectByMapping(m.SourceTraits, (*ComputeSourceTraitsFilterModel).toSDK)
	template.SourceNamespaces = convert.TFListToStringSlice(m.SourceNamespaces)
	return &template
}

// toSDK converts whichever constraint is set; the schema validators ensure it is exactly one.
func (m *AttributeConstraintModel) toSDK() model.AttributeConstraint {
	switch {
	case m.NumericAttributeConstraint != nil:
		return m.NumericAttributeConstraint.toSDK()
	case m.StringAttributeConstraint != nil:
		return m.StringAttributeConstraint.toSDK()
	default:
		return nil
	}
}

func (m *NumericAttributeConstraintModel) toSDK() *model.NumericAttributeConstraint {
	sdk := model.NumericAttributeConstraint{}
	sdk.Attribute = m.Attribute.ValueString()
	sdk.Min = convert.TFToOptionalFloat64(m.Min)
	sdk.Max = convert.TFToOptionalFloat64(m.Max)
	return &sdk
}

func (m *StringAttributeConstraintModel) toSDK() *model.StringAttributeConstraint {
	sdk := model.StringAttributeConstraint{}
	sdk.Attribute = m.Attribute.ValueString()
	sdk.AnyOf = convert.TFSetToStringSlice(m.AnyOf)
	sdk.NoneOf = convert.TFSetToStringSlice(m.NoneOf)
	sdk.Pattern = convert.TFToOptionalString(m.Pattern)
	return &sdk
}

// toSDK converts whichever preference is set; the schema validators ensure it is exactly one.
func (m *AttributePreferenceModel) toSDK() model.AttributePreference {
	switch {
	case m.NumericAttributePreference != nil:
		return m.NumericAttributePreference.toSDK()
	case m.StringAttributePreference != nil:
		return m.StringAttributePreference.toSDK()
	default:
		return nil
	}
}

func (m *NumericAttributePreferenceModel) toSDK() *model.NumericAttributePreference {
	sdk := model.NumericAttributePreference{}
	sdk.Attribute = m.Attribute.ValueString()
	sdk.Weight = m.Weight.ValueFloat64()
	sdk.RankOrder = convert.TFToOptionalEnum[model.NumericAttributeRankOrder](m.RankOrder)
	return &sdk
}

func (m *StringAttributePreferenceModel) toSDK() *model.StringAttributePreference {
	sdk := model.StringAttributePreference{}
	sdk.Attribute = m.Attribute.ValueString()
	sdk.Weight = m.Weight.ValueFloat64()
	sdk.PreferredValues = convert.TFListToStringSlice(m.PreferredValues)
	sdk.PreferredPatterns = convert.TFListToStringSlice(m.PreferredPatterns)
	return &sdk
}

func (m *ComputeSourceTraitsFilterModel) toSDK() model.ComputeSourceTraitsFilter {
	return model.ComputeSourceTraitsFilter{
		CanStopStart:     convert.TFToOptionalBool(m.CanStopStart),
		CanRestart:       convert.TFToOptionalBool(m.CanRestart),
		CanScaleOut:      convert.TFToOptionalBool(m.CanScaleOut),
		IsSelfMaintained: convert.TFToOptionalBool(m.IsSelfMaintained),
	}
}

func (m *computeRequirementDynamicTemplateModel) updateFromSDK(ctx context.Context, template *model.ComputeRequirementDynamicTemplate, diagnostics *diag.Diagnostics) {
	m.Id = convert.OptionalStringToTF(template.Id)
	m.Name = types.StringValue(template.Name)
	m.Namespace = convert.OptionalStringToTF(template.Namespace)
	m.Description = convert.OptionalStringToTF(template.Description)
	m.StrategyType = types.StringValue(template.StrategyType)
	m.MinimumSourceCount = convert.OptionalInt32ToTF(template.MinimumSourceCount)
	m.MaximumSourceCount = convert.OptionalInt32ToTF(template.MaximumSourceCount)
	m.ImagesId = convert.OptionalStringToTF(template.ImagesId)
	m.UserData = convert.OptionalStringToTF(template.UserData)
	m.InstanceTags = convert.StringMapToTF(template.InstanceTags)
	m.Constraints = convert.SliceToTFList(ctx, template.Constraints, createAttributeConstraintModelFromSDK, attributeConstraintObjectType(), diagnostics)
	m.Preferences = convert.SliceToTFList(ctx, template.Preferences, createAttributePreferenceModelFromSDK, attributePreferenceObjectType(), diagnostics)
	m.updateSourceTraitsFromSDK(template)
	m.SourceNamespaces = convert.StringSliceToTFList(template.SourceNamespaces)
}

func createAttributeConstraintModelFromSDK(constraint *model.AttributeConstraint, diags *diag.Diagnostics) AttributeConstraintModel {
	switch typedConstraint := (*constraint).(type) {
	case *model.NumericAttributeConstraint:
		return AttributeConstraintModel{
			NumericAttributeConstraint: &NumericAttributeConstraintModel{
				Attribute: types.StringValue(typedConstraint.Attribute),
				Min:       convert.OptionalFloat64ToTF(typedConstraint.Min),
				Max:       convert.OptionalFloat64ToTF(typedConstraint.Max),
			},
		}
	case *model.StringAttributeConstraint:
		return AttributeConstraintModel{
			StringAttributeConstraint: &StringAttributeConstraintModel{
				Attribute: types.StringValue(typedConstraint.Attribute),
				AnyOf:     convert.StringSliceToTFSet(typedConstraint.AnyOf),
				NoneOf:    convert.StringSliceToTFSet(typedConstraint.NoneOf),
				Pattern:   convert.OptionalStringToTF(typedConstraint.Pattern),
			},
		}
	default:
		diags.AddError(
			"Unsupported attribute constraint type",
			fmt.Sprintf("%T is not yet supported by this provider", constraint),
		)
		return AttributeConstraintModel{}
	}
}

func attributeConstraintObjectType() types.ObjectType {
	return schema.NestedAttributeObject{Attributes: attributeConstraintSchema()}.Type().(types.ObjectType)
}

func createAttributePreferenceModelFromSDK(preference *model.AttributePreference, diags *diag.Diagnostics) AttributePreferenceModel {
	switch typedPreference := (*preference).(type) {
	case *model.NumericAttributePreference:
		return AttributePreferenceModel{
			NumericAttributePreference: &NumericAttributePreferenceModel{
				Attribute: types.StringValue(typedPreference.Attribute),
				Weight:    types.Float64Value(typedPreference.Weight),
				RankOrder: convert.OptionalEnumToTF(typedPreference.RankOrder),
			},
		}
	case *model.StringAttributePreference:
		return AttributePreferenceModel{
			StringAttributePreference: &StringAttributePreferenceModel{
				Attribute:         types.StringValue(typedPreference.Attribute),
				Weight:            types.Float64Value(typedPreference.Weight),
				PreferredValues:   convert.StringSliceToTFList(typedPreference.PreferredValues),
				PreferredPatterns: convert.StringSliceToTFList(typedPreference.PreferredPatterns),
			},
		}
	default:
		diags.AddError(
			"Unsupported attribute preference type",
			fmt.Sprintf("%T is not yet supported by this provider", typedPreference),
		)
		return AttributePreferenceModel{}
	}
}

func attributePreferenceObjectType() types.ObjectType {
	return schema.NestedAttributeObject{Attributes: attributePreferenceSchema()}.Type().(types.ObjectType)
}

func (m *computeRequirementDynamicTemplateModel) updateSourceTraitsFromSDK(template *model.ComputeRequirementDynamicTemplate) {
	sourceTraits, isPresent := template.SourceTraits.Get()
	if !isPresent {
		m.SourceTraits = nil
		return
	}
	m.SourceTraits = &ComputeSourceTraitsFilterModel{
		CanStopStart:     convert.OptionalBoolToTF(sourceTraits.CanStopStart),
		CanRestart:       convert.OptionalBoolToTF(sourceTraits.CanRestart),
		CanScaleOut:      convert.OptionalBoolToTF(sourceTraits.CanScaleOut),
		IsSelfMaintained: convert.OptionalBoolToTF(sourceTraits.IsSelfMaintained),
	}
}
