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

func (r *computeSourceTemplateResource) create(ctx context.Context, state *computeSourceTemplateModel, diags *diag.Diagnostics) (bool, error) {
	template := state.toSDK(ctx, diags)
	if diags.HasError() {
		return NotCreated, nil
	}
	computeSourceTemplate, err := r.client.Compute.AddComputeSourceTemplate(ctx, template)
	if err != nil {
		return NotCreated, err
	}
	state.updateFromSDK(ctx, computeSourceTemplate, diags)
	return Created, nil
}

func (r *computeSourceTemplateResource) read(ctx context.Context, state *computeSourceTemplateModel, diagnostics *diag.Diagnostics) error {
	computeSourceTemplate, err := r.client.Compute.GetComputeSourceTemplate(ctx, state.Id.ValueString())
	if err != nil {
		return err
	}
	state.updateFromSDK(ctx, computeSourceTemplate, diagnostics)
	return nil
}

func (r *computeSourceTemplateResource) update(ctx context.Context, state *computeSourceTemplateModel, diags *diag.Diagnostics) error {
	computeSourceTemplate := state.toSDK(ctx, diags)
	if diags.HasError() {
		return nil
	}
	updatedComputeSourceTemplate, err := r.client.Compute.UpdateComputeSourceTemplate(ctx, computeSourceTemplate)
	if err != nil {
		return err
	}
	state.updateFromSDK(ctx, updatedComputeSourceTemplate, diags)
	return nil
}

func (r *computeSourceTemplateResource) delete(ctx context.Context, state computeSourceTemplateModel, _ *diag.Diagnostics) error {
	return r.client.Compute.DeleteComputeSourceTemplateById(ctx, state.Id.ValueString())
}

func (m *computeSourceTemplateModel) toSDK(ctx context.Context, diags *diag.Diagnostics) *model.ComputeSourceTemplate {
	return &model.ComputeSourceTemplate{
		Id:          convert.TFToOptionalString(m.Id),
		Namespace:   convert.TFToOptionalString(m.Namespace),
		Description: convert.TFToOptionalString(m.Description),
		Attributes:  convert.TFListToSlice(ctx, m.Attributes, (*AttributeValueModel).toSDK, diags),
		Source:      m.Source.toSDK(ctx, diags),
	}
}

// toSDK converts whichever attribute value is set; the schema validators ensure it is exactly one.
func (m *AttributeValueModel) toSDK() model.AttributeValue {
	switch {
	case m.StringAttributeValue != nil:
		sdkAttribute := model.StringAttributeValue{}
		sdkAttribute.Attribute = m.StringAttributeValue.Attribute.ValueString()
		sdkAttribute.Value = convert.TFToOptionalString(m.StringAttributeValue.Value)
		return &sdkAttribute
	case m.NumericAttributeValue != nil:
		sdkAttribute := model.NumericAttributeValue{}
		sdkAttribute.Attribute = m.NumericAttributeValue.Attribute.ValueString()
		sdkAttribute.Value = m.NumericAttributeValue.Value.ValueFloat64()
		return &sdkAttribute
	default:
		return nil
	}
}

// toSDK converts whichever compute source is set; the schema validators ensure it is exactly one.
func (m *ComputeSourceModel) toSDK(ctx context.Context, diags *diag.Diagnostics) model.ComputeSource {
	switch {
	case m == nil:
		diags.AddError("Missing compute source", "The compute source template has no compute source")
		return nil
	case m.AwsFleetComputeSource != nil:
		return m.AwsFleetComputeSource.toSDK(ctx, diags)
	case m.AwsInstancesComputeSource != nil:
		return m.AwsInstancesComputeSource.toSDK(ctx, diags)
	case m.AzureInstancesComputeSource != nil:
		return m.AzureInstancesComputeSource.toSDK()
	case m.AzureScaleSetComputeSource != nil:
		return m.AzureScaleSetComputeSource.toSDK()
	case m.GceInstancesComputeSource != nil:
		return m.GceInstancesComputeSource.toSDK()
	case m.GceInstanceGroupComputeSource != nil:
		return m.GceInstanceGroupComputeSource.toSDK()
	case m.OciInstancePoolComputeSource != nil:
		return m.OciInstancePoolComputeSource.toSDK()
	case m.OciInstancesComputeSource != nil:
		return m.OciInstancesComputeSource.toSDK()
	case m.SimulatorComputeSource != nil:
		return m.SimulatorComputeSource.toSDK()
	default:
		diags.AddError("Unsupported compute source", "The compute source template has no compute source of a supported type")
		return nil
	}
}

func (m *AwsFleetComputeSourceModel) toSDK(ctx context.Context, diags *diag.Diagnostics) *model.AwsFleetComputeSource {
	source := model.AwsFleetComputeSource{}
	source.Id = convert.TFToOptionalString(m.Id)
	source.Name = m.Name.ValueString()
	source.Credential = m.Credential.ValueString()
	source.Region = m.Region.ValueString()
	source.AvailabilityZone = convert.TFToOptionalString(m.AvailabilityZone)
	source.SecurityGroupId = m.SecurityGroupId.ValueString()
	source.SubnetId = convert.TFToOptionalString(m.SubnetId)
	source.InstanceType = m.InstanceType.ValueString()
	source.ImageId = m.ImageId.ValueString()
	source.UserData = convert.TFToOptionalString(m.UserData)
	source.InstanceTags = convert.TFToStringMap(m.InstanceTags)
	source.BootVolumeSizeGb = convert.TFToOptionalInt32(m.BootVolumeSizeGb)
	source.IamInstanceProfileArn = convert.TFToOptionalString(m.IamInstanceProfileArn)
	source.KeyName = convert.TFToOptionalString(m.KeyName)
	source.EnableDetailedMonitoring = convert.TFToOptionalBool(m.EnableDetailedMonitoring)
	source.EnableInstanceMetadataTags = convert.TFToOptionalBool(m.EnableInstanceMetadataTags)
	source.InstanceMetadataHttpPutResponseHopLimit = convert.TFToOptionalInt32(m.InstanceMetadataHttpPutResponseHopLimit)
	source.AssignPublicIp = m.AssignPublicIp.ValueBool()
	source.CreateClusterPlacementGroup = convert.TFToOptionalBool(m.CreateClusterPlacementGroup)
	source.ExistingPlacementGroup = convert.TFToOptionalObjectByMapping(m.ExistingPlacementGroup, (*AwsPlacementGroupModel).toSDK)
	source.CreateElasticFabricAdapter = convert.TFToOptionalBool(m.CreateElasticFabricAdapter)
	source.SecondaryNetworkInterfaces = convert.TFListToSlice(ctx, m.SecondaryNetworkInterfaces, (*AwsSecondaryNetworkInterfaceModel).toSDK, diags)
	source.CapacityReservation = convert.TFToOptionalObjectByMapping(m.CapacityReservation, (*AwsCapacityReservationModel).toSDK)
	source.Limit = int(m.Limit.ValueInt32())
	source.MaintainCapacity = m.MaintainCapacity.ValueBool()
	source.InstanceOverrides = convert.TFListToSlice(ctx, m.InstanceOverrides, (*AwsFleetInstanceOverrideModel).toSDK, diags)
	source.PurchaseOption = convert.TFToEnum[model.AwsFleetPurchaseOption](m.PurchaseOption)
	source.OnDemandOptions = convert.TFToOptionalObjectByMapping(m.OnDemandOptions, (*AwsFleetOnDemandOptionsModel).toSDK)
	source.SpotOptions = convert.TFToOptionalObjectByMapping(m.SpotOptions, (*AwsFleetSpotOptionsModel).toSDK)
	return &source
}

func (m *AwsInstancesComputeSourceModel) toSDK(ctx context.Context, diags *diag.Diagnostics) *model.AwsInstancesComputeSource {
	source := model.AwsInstancesComputeSource{}
	source.Id = convert.TFToOptionalString(m.Id)
	source.Name = m.Name.ValueString()
	source.Credential = m.Credential.ValueString()
	source.Region = m.Region.ValueString()
	source.AvailabilityZone = convert.TFToOptionalString(m.AvailabilityZone)
	source.SecurityGroupId = m.SecurityGroupId.ValueString()
	source.SubnetId = convert.TFToOptionalString(m.SubnetId)
	source.InstanceType = m.InstanceType.ValueString()
	source.ImageId = m.ImageId.ValueString()
	source.UserData = convert.TFToOptionalString(m.UserData)
	source.InstanceTags = convert.TFToStringMap(m.InstanceTags)
	source.BootVolumeSizeGb = convert.TFToOptionalInt32(m.BootVolumeSizeGb)
	source.IamInstanceProfileArn = convert.TFToOptionalString(m.IamInstanceProfileArn)
	source.KeyName = convert.TFToOptionalString(m.KeyName)
	source.EnableDetailedMonitoring = convert.TFToOptionalBool(m.EnableDetailedMonitoring)
	source.EnableInstanceMetadataTags = convert.TFToOptionalBool(m.EnableInstanceMetadataTags)
	source.InstanceMetadataHttpPutResponseHopLimit = convert.TFToOptionalInt32(m.InstanceMetadataHttpPutResponseHopLimit)
	source.UseCapacityBlock = convert.TFToOptionalBool(m.UseCapacityBlock)
	source.AssignPublicIp = m.AssignPublicIp.ValueBool()
	source.CreateClusterPlacementGroup = convert.TFToOptionalBool(m.CreateClusterPlacementGroup)
	source.ExistingPlacementGroup = convert.TFToOptionalObjectByMapping(m.ExistingPlacementGroup, (*AwsPlacementGroupModel).toSDK)
	source.CreateElasticFabricAdapter = convert.TFToOptionalBool(m.CreateElasticFabricAdapter)
	source.SecondaryNetworkInterfaces = convert.TFListToSlice(ctx, m.SecondaryNetworkInterfaces, (*AwsSecondaryNetworkInterfaceModel).toSDK, diags)
	source.CapacityReservation = convert.TFToOptionalObjectByMapping(m.CapacityReservation, (*AwsCapacityReservationModel).toSDK)
	source.Limit = int(m.Limit.ValueInt32())
	source.SpecifyMinimum = m.SpecifyMinimum.ValueBool()
	source.Spot = m.Spot.ValueBool()
	source.SpotMaxPrice = convert.TFToOptionalFloat64(m.SpotMaxPrice)
	return &source
}

func (m *AzureInstancesComputeSourceModel) toSDK() *model.AzureInstancesComputeSource {
	source := model.AzureInstancesComputeSource{}
	source.Id = convert.TFToOptionalString(m.Id)
	source.Name = m.Name.ValueString()
	source.Credential = m.Credential.ValueString()
	source.Region = m.Region.ValueString()
	source.AvailabilityZone = convert.TFToOptionalString(m.AvailabilityZone)
	source.NetworkResourceGroupName = m.NetworkResourceGroupName.ValueString()
	source.NetworkName = m.NetworkName.ValueString()
	source.SubnetName = m.SubnetName.ValueString()
	source.Environment = convert.TFToOptionalString(m.Environment)
	source.VmSize = m.VmSize.ValueString()
	source.ImageId = m.ImageId.ValueString()
	source.UserData = convert.TFToOptionalString(m.UserData)
	source.InstanceTags = convert.TFToStringMap(m.InstanceTags)
	source.BootVolumeSizeGb = convert.TFToOptionalInt32(m.BootVolumeSizeGb)
	source.AdminUserCredential = convert.TFToOptionalString(m.AdminUserCredential)
	source.SshKey = convert.TFToOptionalString(m.SshKey)
	source.AssignPublicIp = m.AssignPublicIp.ValueBool()
	source.CreateProximityPlacementGroup = convert.TFToOptionalBool(m.CreateProximityPlacementGroup)
	source.UseAcceleratedNetworking = convert.TFToOptionalBool(m.UseAcceleratedNetworking)
	source.UseSpot = convert.TFToOptionalBool(m.UseSpot)
	source.SpotMaxPrice = convert.TFToOptionalFloat64(m.SpotMaxPrice)
	source.Limit = int(m.Limit.ValueInt32())
	return &source
}

func (m *AzureScaleSetComputeSourceModel) toSDK() *model.AzureScaleSetComputeSource {
	source := model.AzureScaleSetComputeSource{}
	source.Id = convert.TFToOptionalString(m.Id)
	source.Name = m.Name.ValueString()
	source.Credential = m.Credential.ValueString()
	source.Region = m.Region.ValueString()
	source.AvailabilityZone = convert.TFToOptionalString(m.AvailabilityZone)
	source.NetworkResourceGroupName = m.NetworkResourceGroupName.ValueString()
	source.NetworkName = m.NetworkName.ValueString()
	source.SubnetName = m.SubnetName.ValueString()
	source.Environment = convert.TFToOptionalString(m.Environment)
	source.VmSize = m.VmSize.ValueString()
	source.ImageId = m.ImageId.ValueString()
	source.UserData = convert.TFToOptionalString(m.UserData)
	source.InstanceTags = convert.TFToStringMap(m.InstanceTags)
	source.BootVolumeSizeGb = convert.TFToOptionalInt32(m.BootVolumeSizeGb)
	source.AdminUserCredential = convert.TFToOptionalString(m.AdminUserCredential)
	source.SshKey = convert.TFToOptionalString(m.SshKey)
	source.AssignPublicIp = m.AssignPublicIp.ValueBool()
	source.CreateProximityPlacementGroup = convert.TFToOptionalBool(m.CreateProximityPlacementGroup)
	source.UseAcceleratedNetworking = convert.TFToOptionalBool(m.UseAcceleratedNetworking)
	source.UseSpot = convert.TFToOptionalBool(m.UseSpot)
	source.SpotMaxPrice = convert.TFToOptionalFloat64(m.SpotMaxPrice)
	source.Limit = int(m.Limit.ValueInt32())
	return &source
}

func (m *GceInstancesComputeSourceModel) toSDK() *model.GceInstancesComputeSource {
	source := model.GceInstancesComputeSource{}
	source.Id = convert.TFToOptionalString(m.Id)
	source.Name = m.Name.ValueString()
	source.Credential = m.Credential.ValueString()
	source.Limit = int(m.Limit.ValueInt32())
	source.AssignPublicIp = m.AssignPublicIp.ValueBool()
	source.SpecifyMinimum = m.SpecifyMinimum.ValueBool()
	source.CreateCompactPlacementPolicy = convert.TFToOptionalBool(m.CreateCompactPlacementPolicy)
	source.UserData = convert.TFToOptionalString(m.UserData)
	source.InstanceTags = convert.TFToStringMap(m.InstanceTags)
	source.BootVolumeSizeGb = convert.TFToOptionalInt32(m.BootVolumeSizeGb)
	source.SshKeys = convert.TFToOptionalString(m.SshKeys)
	source.ServiceAccountEmail = convert.TFToOptionalString(m.ServiceAccountEmail)
	source.Reservation = convert.TFToOptionalString(m.Reservation)
	source.ReservationBound = m.ReservationBound.ValueBool()
	source.Project = m.Project.ValueString()
	source.Region = m.Region.ValueString()
	source.Zone = convert.TFToOptionalString(m.Zone)
	source.MachineType = m.MachineType.ValueString()
	source.Image = m.Image.ValueString()
	source.Network = convert.TFToOptionalString(m.Network)
	source.Subnetwork = convert.TFToOptionalString(m.Subnetwork)
	source.Preemptible = m.Preemptible.ValueBool()
	source.Spot = m.Spot.ValueBool()
	source.Confidential = m.Confidential.ValueBool()
	source.AcceleratorType = convert.TFToOptionalString(m.AcceleratorType)
	source.AcceleratorCount = int(m.AcceleratorCount.ValueInt32())
	source.HostMaintenanceBehaviour = convert.TFToOptionalEnum[model.GceHostMaintenanceBehaviour](m.HostMaintenanceBehaviour)
	return &source
}

func (m *GceInstanceGroupComputeSourceModel) toSDK() *model.GceInstanceGroupComputeSource {
	source := model.GceInstanceGroupComputeSource{}
	source.Id = convert.TFToOptionalString(m.Id)
	source.Name = m.Name.ValueString()
	source.Credential = m.Credential.ValueString()
	source.Limit = int(m.Limit.ValueInt32())
	source.AssignPublicIp = m.AssignPublicIp.ValueBool()
	source.UserData = convert.TFToOptionalString(m.UserData)
	source.InstanceTags = convert.TFToStringMap(m.InstanceTags)
	source.BootVolumeSizeGb = convert.TFToOptionalInt32(m.BootVolumeSizeGb)
	source.SshKeys = convert.TFToOptionalString(m.SshKeys)
	source.ServiceAccountEmail = convert.TFToOptionalString(m.ServiceAccountEmail)
	source.Reservation = convert.TFToOptionalString(m.Reservation)
	source.ReservationBound = m.ReservationBound.ValueBool()
	source.Project = m.Project.ValueString()
	source.Region = m.Region.ValueString()
	source.Zone = convert.TFToOptionalString(m.Zone)
	source.MachineType = m.MachineType.ValueString()
	source.Image = m.Image.ValueString()
	source.Network = convert.TFToOptionalString(m.Network)
	source.Subnetwork = convert.TFToOptionalString(m.Subnetwork)
	source.Preemptible = m.Preemptible.ValueBool()
	source.Spot = m.Spot.ValueBool()
	source.Confidential = m.Confidential.ValueBool()
	source.AcceleratorType = convert.TFToOptionalString(m.AcceleratorType)
	source.AcceleratorCount = int(m.AcceleratorCount.ValueInt32())
	source.HostMaintenanceBehaviour = convert.TFToOptionalEnum[model.GceHostMaintenanceBehaviour](m.HostMaintenanceBehaviour)
	source.TargetDistributionShape = convert.TFToOptionalEnum[model.GceTargetDistributionShape](m.TargetDistributionShape)
	return &source
}

func (m *OciInstancePoolComputeSourceModel) toSDK() *model.OciInstancePoolComputeSource {
	source := model.OciInstancePoolComputeSource{}
	source.Id = convert.TFToOptionalString(m.Id)
	source.Name = m.Name.ValueString()
	source.Credential = m.Credential.ValueString()
	source.SshKey = convert.TFToOptionalString(m.SshKey)
	source.Region = m.Region.ValueString()
	source.AvailabilityDomain = convert.TFToOptionalString(m.AvailabilityDomain)
	source.CompartmentId = m.CompartmentId.ValueString()
	source.ImageId = m.ImageId.ValueString()
	source.Shape = m.Shape.ValueString()
	source.FlexOcpus = convert.TFToOptionalFloat32(m.FlexOcpus)
	source.FlexRam = convert.TFToOptionalFloat32(m.FlexRam)
	source.SubnetId = m.SubnetId.ValueString()
	source.Limit = int(m.Limit.ValueInt32())
	source.AssignPublicIp = m.AssignPublicIp.ValueBool()
	source.CreateClusterNetwork = convert.TFToOptionalBool(m.CreateClusterNetwork)
	source.UserData = convert.TFToOptionalString(m.UserData)
	source.InstanceTags = convert.TFToStringMap(m.InstanceTags)
	source.BootVolumeSizeGb = convert.TFToOptionalInt32(m.BootVolumeSizeGb)
	return &source
}

func (m *OciInstancesComputeSourceModel) toSDK() *model.OciInstancesComputeSource {
	source := model.OciInstancesComputeSource{}
	source.Id = convert.TFToOptionalString(m.Id)
	source.Name = m.Name.ValueString()
	source.Credential = m.Credential.ValueString()
	source.SshKey = convert.TFToOptionalString(m.SshKey)
	source.Region = m.Region.ValueString()
	source.AvailabilityDomain = convert.TFToOptionalString(m.AvailabilityDomain)
	source.CompartmentId = m.CompartmentId.ValueString()
	source.ImageId = m.ImageId.ValueString()
	source.Shape = m.Shape.ValueString()
	source.FlexOcpus = convert.TFToOptionalFloat32(m.FlexOcpus)
	source.FlexRam = convert.TFToOptionalFloat32(m.FlexRam)
	source.Preemptible = m.Preemptible.ValueBool()
	source.SubnetId = m.SubnetId.ValueString()
	source.Limit = int(m.Limit.ValueInt32())
	source.AssignPublicIp = m.AssignPublicIp.ValueBool()
	source.UserData = convert.TFToOptionalString(m.UserData)
	source.InstanceTags = convert.TFToStringMap(m.InstanceTags)
	source.BootVolumeSizeGb = convert.TFToOptionalInt32(m.BootVolumeSizeGb)
	return &source
}

func (m *SimulatorComputeSourceModel) toSDK() *model.SimulatorComputeSource {
	source := model.SimulatorComputeSource{}
	source.Id = convert.TFToOptionalString(m.Id)
	source.Name = m.Name.ValueString()
	source.Credential = convert.TFToOptionalString(m.Credential)
	// The Go SDK cannot leave these out, so an unset one is sent as the platform's default rather than ""
	source.Region = convert.TFToOptionalString(m.Region).Or("sim-region")
	source.InstanceType = convert.TFToOptionalString(m.InstanceType).Or("sim-instance")
	source.ImageId = convert.TFToOptionalString(m.ImageId).Or("sim-image")
	source.InstanceTags = convert.TFToStringMap(m.InstanceTags)
	source.BootVolumeSizeGb = convert.TFToOptionalInt32(m.BootVolumeSizeGb)
	source.ImplicitCapacity = convert.TFToOptionalInt32(m.ImplicitCapacity)
	source.InstanceStartupTimeSeconds = int(m.InstanceStartupTimeSeconds.ValueInt32())
	source.InstanceStartupTimeVariance = m.InstanceStartupTimeVariance.ValueFloat32()
	source.InstanceShutdownTimeSeconds = int(m.InstanceShutdownTimeSeconds.ValueInt32())
	source.InstanceShutdownTimeVariance = m.InstanceShutdownTimeVariance.ValueFloat32()
	source.UnexpectedInstanceTerminationProbability = m.UnexpectedInstanceTerminationProbability.ValueFloat32()
	source.FailOnRequestAllInstances = m.FailOnRequestAllInstances.ValueBool()
	source.FailOnProvision = m.FailOnProvision.ValueBool()
	source.ReportSupportingResourcesExist = m.ReportSupportingResourcesExist.ValueBool()
	source.Limit = int(m.Limit.ValueInt32())
	return &source
}

func (m *AwsPlacementGroupModel) toSDK() model.AwsPlacementGroup {
	return model.AwsPlacementGroup{
		GroupName: convert.TFToOptionalString(m.GroupName),
		GroupId:   convert.TFToOptionalString(m.GroupId),
	}
}

func (m *AwsSecondaryNetworkInterfaceModel) toSDK() model.AwsSecondaryNetworkInterface {
	return model.AwsSecondaryNetworkInterface{
		NetworkInterfaceType: convert.TFToEnum[model.AwsNetworkInterfaceType](m.NetworkInterfaceType),
		SubnetId:             convert.TFToOptionalString(m.SubnetId),
		SecurityGroupId:      convert.TFToOptionalString(m.SecurityGroupId),
	}
}

func (m *AwsCapacityReservationModel) toSDK() model.AwsCapacityReservation {
	return model.AwsCapacityReservation{
		Preference: convert.TFToEnum[model.AwsCapacityReservationPreference](m.Preference),
		GroupArn:   convert.TFToOptionalString(m.GroupArn),
		Id:         convert.TFToOptionalString(m.Id),
	}
}

func (m *AwsFleetInstanceOverrideModel) toSDK() model.AwsFleetInstanceOverride {
	return model.AwsFleetInstanceOverride{
		AvailabilityZone: convert.TFToOptionalString(m.AvailabilityZone),
		InstanceType:     m.InstanceType.ValueString(),
		SpotMaxPrice:     convert.TFToOptionalFloat64(m.SpotMaxPrice),
		SubnetId:         convert.TFToOptionalString(m.SubnetId),
	}
}

func (m *AwsFleetOnDemandOptionsModel) toSDK() model.AwsFleetOnDemandOptions {
	return model.AwsFleetOnDemandOptions{
		AllocationStrategy:           convert.TFToOptionalEnum[model.AwsFleetOnDemandAllocationStrategy](m.AllocationStrategy),
		UseCapacityReservationsFirst: convert.TFToOptionalBool(m.UseCapacityReservationsFirst),
		MaxTotalPrice:                convert.TFToOptionalFloat64(m.MaxTotalPrice),
		MinInstanceCount:             convert.TFToOptionalInt32(m.MinInstanceCount),
		SingleAvailabilityZone:       convert.TFToOptionalBool(m.SingleAvailabilityZone),
		SingleInstanceType:           convert.TFToOptionalBool(m.SingleInstanceType),
	}
}

func (m *AwsFleetSpotOptionsModel) toSDK() model.AwsFleetSpotOptions {
	return model.AwsFleetSpotOptions{
		AllocationStrategy:                   convert.TFToOptionalEnum[model.AwsFleetSpotAllocationStrategy](m.AllocationStrategy),
		InstancePoolsToUseCount:              convert.TFToOptionalInt32(m.InstancePoolsToUseCount),
		LaunchReplacementInstanceOnRebalance: convert.TFToOptionalBool(m.LaunchReplacementInstanceOnRebalance),
		MaxTotalPrice:                        convert.TFToOptionalFloat64(m.MaxTotalPrice),
		MinInstanceCount:                     convert.TFToOptionalInt32(m.MinInstanceCount),
		SingleAvailabilityZone:               convert.TFToOptionalBool(m.SingleAvailabilityZone),
		SingleInstanceType:                   convert.TFToOptionalBool(m.SingleInstanceType),
	}
}

func (m *computeSourceTemplateModel) updateFromSDK(ctx context.Context, template *model.ComputeSourceTemplate, diagnostics *diag.Diagnostics) {
	m.Id = convert.OptionalStringToTF(template.Id)
	m.Namespace = convert.OptionalStringToTF(template.Namespace)
	m.Description = convert.OptionalStringToTF(template.Description)
	m.Attributes = convert.SliceToTFList(ctx, template.Attributes, createAttributeValueModelFromSDK, attributeValueObjectType(), diagnostics)
	m.updateSourceFromSDK(ctx, template.Source, diagnostics)
}

func createAttributeValueModelFromSDK(attribute *model.AttributeValue, diags *diag.Diagnostics) AttributeValueModel {
	switch a := (*attribute).(type) {
	case *model.StringAttributeValue:
		return AttributeValueModel{
			StringAttributeValue: &StringAttributeValueModel{
				Attribute: types.StringValue(a.Attribute),
				Value:     convert.OptionalStringToTF(a.Value),
			},
		}
	case *model.NumericAttributeValue:
		return AttributeValueModel{
			NumericAttributeValue: &NumericAttributeValueModel{
				Attribute: types.StringValue(a.Attribute),
				Value:     types.Float64Value(a.Value),
			},
		}
	default:
		diags.AddError("Unsupported attribute value", fmt.Sprintf("The compute source template has an attribute value of unsupported type %T", attribute))
		return AttributeValueModel{}
	}
}

func attributeValueObjectType() types.ObjectType {
	return schema.NestedAttributeObject{Attributes: attributeValueSchema()}.Type().(types.ObjectType)
}

// updateSourceFromSDK sets the compute source of the source's type, and unsets the others.
func (m *computeSourceTemplateModel) updateSourceFromSDK(ctx context.Context, source model.ComputeSource, diags *diag.Diagnostics) {
	m.Source = &ComputeSourceModel{}
	switch s := source.(type) {
	case *model.AwsFleetComputeSource:
		m.Source.AwsFleetComputeSource = new(AwsFleetComputeSourceModel)
		m.Source.AwsFleetComputeSource.updateFromSDK(ctx, s, diags)
	case *model.AwsInstancesComputeSource:
		m.Source.AwsInstancesComputeSource = new(AwsInstancesComputeSourceModel)
		m.Source.AwsInstancesComputeSource.updateFromSDK(ctx, s, diags)
	case *model.AzureInstancesComputeSource:
		m.Source.AzureInstancesComputeSource = new(AzureInstancesComputeSourceModel)
		m.Source.AzureInstancesComputeSource.updateFromSDK(s)
	case *model.AzureScaleSetComputeSource:
		m.Source.AzureScaleSetComputeSource = new(AzureScaleSetComputeSourceModel)
		m.Source.AzureScaleSetComputeSource.updateFromSDK(s)
	case *model.GceInstancesComputeSource:
		m.Source.GceInstancesComputeSource = new(GceInstancesComputeSourceModel)
		m.Source.GceInstancesComputeSource.updateFromSDK(s)
	case *model.GceInstanceGroupComputeSource:
		m.Source.GceInstanceGroupComputeSource = new(GceInstanceGroupComputeSourceModel)
		m.Source.GceInstanceGroupComputeSource.updateFromSDK(s)
	case *model.OciInstancePoolComputeSource:
		m.Source.OciInstancePoolComputeSource = new(OciInstancePoolComputeSourceModel)
		m.Source.OciInstancePoolComputeSource.updateFromSDK(s)
	case *model.OciInstancesComputeSource:
		m.Source.OciInstancesComputeSource = new(OciInstancesComputeSourceModel)
		m.Source.OciInstancesComputeSource.updateFromSDK(s)
	case *model.SimulatorComputeSource:
		m.Source.SimulatorComputeSource = new(SimulatorComputeSourceModel)
		m.Source.SimulatorComputeSource.updateFromSDK(s)
	default:
		diags.AddError("Unsupported compute source", fmt.Sprintf("The compute source template has a compute source of unsupported type %T", source))
	}
}

func (m *AwsFleetComputeSourceModel) updateFromSDK(ctx context.Context, s *model.AwsFleetComputeSource, diags *diag.Diagnostics) {
	m.Id = convert.OptionalStringToTF(s.Id)
	m.Name = types.StringValue(s.Name)
	m.Credential = types.StringValue(s.Credential)
	m.Region = types.StringValue(s.Region)
	m.AvailabilityZone = convert.OptionalStringToTF(s.AvailabilityZone)
	m.SecurityGroupId = types.StringValue(s.SecurityGroupId)
	m.SubnetId = convert.OptionalStringToTF(s.SubnetId)
	m.InstanceType = types.StringValue(s.InstanceType)
	m.ImageId = types.StringValue(s.ImageId)
	m.UserData = convert.OptionalStringToTF(s.UserData)
	m.InstanceTags = convert.StringMapToTF(s.InstanceTags)
	m.BootVolumeSizeGb = convert.OptionalInt32ToTF(s.BootVolumeSizeGb)
	m.IamInstanceProfileArn = convert.OptionalStringToTF(s.IamInstanceProfileArn)
	m.KeyName = convert.OptionalStringToTF(s.KeyName)
	m.EnableDetailedMonitoring = convert.OptionalBoolToTF(s.EnableDetailedMonitoring)
	m.EnableInstanceMetadataTags = convert.OptionalBoolToTF(s.EnableInstanceMetadataTags)
	m.InstanceMetadataHttpPutResponseHopLimit = convert.OptionalInt32ToTF(s.InstanceMetadataHttpPutResponseHopLimit)
	m.AssignPublicIp = types.BoolValue(s.AssignPublicIp)
	m.CreateClusterPlacementGroup = convert.OptionalBoolToTF(s.CreateClusterPlacementGroup)
	m.ExistingPlacementGroup = convert.OptionalObjectToTFByMapping(s.ExistingPlacementGroup, (*AwsPlacementGroupModel).updateFromSDK)
	m.CreateElasticFabricAdapter = convert.OptionalBoolToTF(s.CreateElasticFabricAdapter)
	m.SecondaryNetworkInterfaces = convert.SliceToTFList(ctx, s.SecondaryNetworkInterfaces, createAwsSecondaryNetworkInterfaceModelFromSDK, awsSecondaryNetworkInterfaceObjectType(), diags)
	m.CapacityReservation = convert.OptionalObjectToTFByMapping(s.CapacityReservation, (*AwsCapacityReservationModel).updateFromSDK)
	m.Limit = types.Int32Value(int32(s.Limit))
	m.MaintainCapacity = types.BoolValue(s.MaintainCapacity)
	m.InstanceOverrides = convert.SliceToTFList(ctx, s.InstanceOverrides, createAwsFleetInstanceOverrideModelFromSDK, awsFleetInstanceOverrideObjectType(), diags)
	m.PurchaseOption = convert.EnumToTF(s.PurchaseOption)
	m.OnDemandOptions = convert.OptionalObjectToTFByMapping(s.OnDemandOptions, (*AwsFleetOnDemandOptionsModel).updateFromSDK)
	m.SpotOptions = convert.OptionalObjectToTFByMapping(s.SpotOptions, (*AwsFleetSpotOptionsModel).updateFromSDK)
}

func awsSecondaryNetworkInterfaceObjectType() types.ObjectType {
	return schema.NestedAttributeObject{Attributes: awsSecondaryNetworkInterfaceSchema()}.Type().(types.ObjectType)
}

func awsFleetInstanceOverrideObjectType() types.ObjectType {
	return schema.NestedAttributeObject{Attributes: awsFleetInstanceOverrideSchema()}.Type().(types.ObjectType)
}

func (m *AwsInstancesComputeSourceModel) updateFromSDK(ctx context.Context, s *model.AwsInstancesComputeSource, diags *diag.Diagnostics) {
	m.Id = convert.OptionalStringToTF(s.Id)
	m.Name = types.StringValue(s.Name)
	m.Credential = types.StringValue(s.Credential)
	m.Region = types.StringValue(s.Region)
	m.AvailabilityZone = convert.OptionalStringToTF(s.AvailabilityZone)
	m.SecurityGroupId = types.StringValue(s.SecurityGroupId)
	m.SubnetId = convert.OptionalStringToTF(s.SubnetId)
	m.InstanceType = types.StringValue(s.InstanceType)
	m.ImageId = types.StringValue(s.ImageId)
	m.UserData = convert.OptionalStringToTF(s.UserData)
	m.InstanceTags = convert.StringMapToTF(s.InstanceTags)
	m.BootVolumeSizeGb = convert.OptionalInt32ToTF(s.BootVolumeSizeGb)
	m.IamInstanceProfileArn = convert.OptionalStringToTF(s.IamInstanceProfileArn)
	m.KeyName = convert.OptionalStringToTF(s.KeyName)
	m.EnableDetailedMonitoring = convert.OptionalBoolToTF(s.EnableDetailedMonitoring)
	m.EnableInstanceMetadataTags = convert.OptionalBoolToTF(s.EnableInstanceMetadataTags)
	m.InstanceMetadataHttpPutResponseHopLimit = convert.OptionalInt32ToTF(s.InstanceMetadataHttpPutResponseHopLimit)
	m.UseCapacityBlock = convert.OptionalBoolToTF(s.UseCapacityBlock)
	m.AssignPublicIp = types.BoolValue(s.AssignPublicIp)
	m.CreateClusterPlacementGroup = convert.OptionalBoolToTF(s.CreateClusterPlacementGroup)
	m.ExistingPlacementGroup = convert.OptionalObjectToTFByMapping(s.ExistingPlacementGroup, (*AwsPlacementGroupModel).updateFromSDK)
	m.CreateElasticFabricAdapter = convert.OptionalBoolToTF(s.CreateElasticFabricAdapter)
	m.SecondaryNetworkInterfaces = convert.SliceToTFList(ctx, s.SecondaryNetworkInterfaces, createAwsSecondaryNetworkInterfaceModelFromSDK, awsSecondaryNetworkInterfaceObjectType(), diags)
	m.CapacityReservation = convert.OptionalObjectToTFByMapping(s.CapacityReservation, (*AwsCapacityReservationModel).updateFromSDK)
	m.Limit = types.Int32Value(int32(s.Limit))
	m.SpecifyMinimum = types.BoolValue(s.SpecifyMinimum)
	m.Spot = types.BoolValue(s.Spot)
	m.SpotMaxPrice = convert.OptionalFloat64ToTF(s.SpotMaxPrice)
}

func (m *AzureInstancesComputeSourceModel) updateFromSDK(s *model.AzureInstancesComputeSource) {
	m.Id = convert.OptionalStringToTF(s.Id)
	m.Name = types.StringValue(s.Name)
	m.Credential = types.StringValue(s.Credential)
	m.Region = types.StringValue(s.Region)
	m.AvailabilityZone = convert.OptionalStringToTF(s.AvailabilityZone)
	m.NetworkResourceGroupName = types.StringValue(s.NetworkResourceGroupName)
	m.NetworkName = types.StringValue(s.NetworkName)
	m.SubnetName = types.StringValue(s.SubnetName)
	m.Environment = convert.OptionalStringToTF(s.Environment)
	m.VmSize = types.StringValue(s.VmSize)
	m.ImageId = types.StringValue(s.ImageId)
	m.UserData = convert.OptionalStringToTF(s.UserData)
	m.InstanceTags = convert.StringMapToTF(s.InstanceTags)
	m.BootVolumeSizeGb = convert.OptionalInt32ToTF(s.BootVolumeSizeGb)
	m.AdminUserCredential = convert.OptionalStringToTF(s.AdminUserCredential)
	m.SshKey = convert.OptionalStringToTF(s.SshKey)
	m.AssignPublicIp = types.BoolValue(s.AssignPublicIp)
	m.CreateProximityPlacementGroup = convert.OptionalBoolToTF(s.CreateProximityPlacementGroup)
	m.UseAcceleratedNetworking = convert.OptionalBoolToTF(s.UseAcceleratedNetworking)
	m.UseSpot = convert.OptionalBoolToTF(s.UseSpot)
	m.SpotMaxPrice = convert.OptionalFloat64ToTF(s.SpotMaxPrice)
	m.Limit = types.Int32Value(int32(s.Limit))
}

func (m *AzureScaleSetComputeSourceModel) updateFromSDK(s *model.AzureScaleSetComputeSource) {
	m.Id = convert.OptionalStringToTF(s.Id)
	m.Name = types.StringValue(s.Name)
	m.Credential = types.StringValue(s.Credential)
	m.Region = types.StringValue(s.Region)
	m.AvailabilityZone = convert.OptionalStringToTF(s.AvailabilityZone)
	m.NetworkResourceGroupName = types.StringValue(s.NetworkResourceGroupName)
	m.NetworkName = types.StringValue(s.NetworkName)
	m.SubnetName = types.StringValue(s.SubnetName)
	m.Environment = convert.OptionalStringToTF(s.Environment)
	m.VmSize = types.StringValue(s.VmSize)
	m.ImageId = types.StringValue(s.ImageId)
	m.UserData = convert.OptionalStringToTF(s.UserData)
	m.InstanceTags = convert.StringMapToTF(s.InstanceTags)
	m.BootVolumeSizeGb = convert.OptionalInt32ToTF(s.BootVolumeSizeGb)
	m.AdminUserCredential = convert.OptionalStringToTF(s.AdminUserCredential)
	m.SshKey = convert.OptionalStringToTF(s.SshKey)
	m.AssignPublicIp = types.BoolValue(s.AssignPublicIp)
	m.CreateProximityPlacementGroup = convert.OptionalBoolToTF(s.CreateProximityPlacementGroup)
	m.UseAcceleratedNetworking = convert.OptionalBoolToTF(s.UseAcceleratedNetworking)
	m.UseSpot = convert.OptionalBoolToTF(s.UseSpot)
	m.SpotMaxPrice = convert.OptionalFloat64ToTF(s.SpotMaxPrice)
	m.Limit = types.Int32Value(int32(s.Limit))
}

func (m *GceInstancesComputeSourceModel) updateFromSDK(s *model.GceInstancesComputeSource) {
	m.Id = convert.OptionalStringToTF(s.Id)
	m.Name = types.StringValue(s.Name)
	m.Credential = types.StringValue(s.Credential)
	m.Limit = types.Int32Value(int32(s.Limit))
	m.AssignPublicIp = types.BoolValue(s.AssignPublicIp)
	m.SpecifyMinimum = types.BoolValue(s.SpecifyMinimum)
	m.CreateCompactPlacementPolicy = convert.OptionalBoolToTF(s.CreateCompactPlacementPolicy)
	m.UserData = convert.OptionalStringToTF(s.UserData)
	m.InstanceTags = convert.StringMapToTF(s.InstanceTags)
	m.BootVolumeSizeGb = convert.OptionalInt32ToTF(s.BootVolumeSizeGb)
	m.SshKeys = convert.OptionalStringToTF(s.SshKeys)
	m.ServiceAccountEmail = convert.OptionalStringToTF(s.ServiceAccountEmail)
	m.Reservation = convert.OptionalStringToTF(s.Reservation)
	m.ReservationBound = types.BoolValue(s.ReservationBound)
	m.Project = types.StringValue(s.Project)
	m.Region = types.StringValue(s.Region)
	m.Zone = convert.OptionalStringToTF(s.Zone)
	m.MachineType = types.StringValue(s.MachineType)
	m.Image = types.StringValue(s.Image)
	m.Network = convert.OptionalStringToTF(s.Network)
	m.Subnetwork = convert.OptionalStringToTF(s.Subnetwork)
	m.Preemptible = types.BoolValue(s.Preemptible)
	m.Spot = types.BoolValue(s.Spot)
	m.Confidential = types.BoolValue(s.Confidential)
	m.AcceleratorType = convert.OptionalStringToTF(s.AcceleratorType)
	m.AcceleratorCount = types.Int32Value(int32(s.AcceleratorCount))
	m.HostMaintenanceBehaviour = convert.OptionalEnumToTF(s.HostMaintenanceBehaviour)
}

func (m *GceInstanceGroupComputeSourceModel) updateFromSDK(s *model.GceInstanceGroupComputeSource) {
	m.Id = convert.OptionalStringToTF(s.Id)
	m.Name = types.StringValue(s.Name)
	m.Credential = types.StringValue(s.Credential)
	m.Limit = types.Int32Value(int32(s.Limit))
	m.AssignPublicIp = types.BoolValue(s.AssignPublicIp)
	m.UserData = convert.OptionalStringToTF(s.UserData)
	m.InstanceTags = convert.StringMapToTF(s.InstanceTags)
	m.BootVolumeSizeGb = convert.OptionalInt32ToTF(s.BootVolumeSizeGb)
	m.SshKeys = convert.OptionalStringToTF(s.SshKeys)
	m.ServiceAccountEmail = convert.OptionalStringToTF(s.ServiceAccountEmail)
	m.Reservation = convert.OptionalStringToTF(s.Reservation)
	m.ReservationBound = types.BoolValue(s.ReservationBound)
	m.Project = types.StringValue(s.Project)
	m.Region = types.StringValue(s.Region)
	m.Zone = convert.OptionalStringToTF(s.Zone)
	m.MachineType = types.StringValue(s.MachineType)
	m.Image = types.StringValue(s.Image)
	m.Network = convert.OptionalStringToTF(s.Network)
	m.Subnetwork = convert.OptionalStringToTF(s.Subnetwork)
	m.Preemptible = types.BoolValue(s.Preemptible)
	m.Spot = types.BoolValue(s.Spot)
	m.Confidential = types.BoolValue(s.Confidential)
	m.AcceleratorType = convert.OptionalStringToTF(s.AcceleratorType)
	m.AcceleratorCount = types.Int32Value(int32(s.AcceleratorCount))
	m.HostMaintenanceBehaviour = convert.OptionalEnumToTF(s.HostMaintenanceBehaviour)
	m.TargetDistributionShape = convert.OptionalEnumToTF(s.TargetDistributionShape)
}

func (m *OciInstancePoolComputeSourceModel) updateFromSDK(s *model.OciInstancePoolComputeSource) {
	m.Id = convert.OptionalStringToTF(s.Id)
	m.Name = types.StringValue(s.Name)
	m.Credential = types.StringValue(s.Credential)
	m.SshKey = convert.OptionalStringToTF(s.SshKey)
	m.Region = types.StringValue(s.Region)
	m.AvailabilityDomain = convert.OptionalStringToTF(s.AvailabilityDomain)
	m.CompartmentId = types.StringValue(s.CompartmentId)
	m.ImageId = types.StringValue(s.ImageId)
	m.Shape = types.StringValue(s.Shape)
	m.FlexOcpus = convert.OptionalFloat32ToTF(s.FlexOcpus)
	m.FlexRam = convert.OptionalFloat32ToTF(s.FlexRam)
	m.SubnetId = types.StringValue(s.SubnetId)
	m.Limit = types.Int32Value(int32(s.Limit))
	m.AssignPublicIp = types.BoolValue(s.AssignPublicIp)
	m.CreateClusterNetwork = convert.OptionalBoolToTF(s.CreateClusterNetwork)
	m.UserData = convert.OptionalStringToTF(s.UserData)
	m.InstanceTags = convert.StringMapToTF(s.InstanceTags)
	m.BootVolumeSizeGb = convert.OptionalInt32ToTF(s.BootVolumeSizeGb)
}

func (m *OciInstancesComputeSourceModel) updateFromSDK(s *model.OciInstancesComputeSource) {
	m.Id = convert.OptionalStringToTF(s.Id)
	m.Name = types.StringValue(s.Name)
	m.Credential = types.StringValue(s.Credential)
	m.SshKey = convert.OptionalStringToTF(s.SshKey)
	m.Region = types.StringValue(s.Region)
	m.AvailabilityDomain = convert.OptionalStringToTF(s.AvailabilityDomain)
	m.CompartmentId = types.StringValue(s.CompartmentId)
	m.ImageId = types.StringValue(s.ImageId)
	m.Shape = types.StringValue(s.Shape)
	m.FlexOcpus = convert.OptionalFloat32ToTF(s.FlexOcpus)
	m.FlexRam = convert.OptionalFloat32ToTF(s.FlexRam)
	m.Preemptible = types.BoolValue(s.Preemptible)
	m.SubnetId = types.StringValue(s.SubnetId)
	m.Limit = types.Int32Value(int32(s.Limit))
	m.AssignPublicIp = types.BoolValue(s.AssignPublicIp)
	m.UserData = convert.OptionalStringToTF(s.UserData)
	m.InstanceTags = convert.StringMapToTF(s.InstanceTags)
	m.BootVolumeSizeGb = convert.OptionalInt32ToTF(s.BootVolumeSizeGb)
}

func (m *SimulatorComputeSourceModel) updateFromSDK(s *model.SimulatorComputeSource) {
	m.Id = convert.OptionalStringToTF(s.Id)
	m.Name = types.StringValue(s.Name)
	m.Credential = convert.OptionalStringToTF(s.Credential)
	m.Region = types.StringValue(s.Region)
	m.InstanceType = types.StringValue(s.InstanceType)
	m.ImageId = types.StringValue(s.ImageId)
	m.InstanceTags = convert.StringMapToTF(s.InstanceTags)
	m.BootVolumeSizeGb = convert.OptionalInt32ToTF(s.BootVolumeSizeGb)
	m.ImplicitCapacity = convert.OptionalInt32ToTF(s.ImplicitCapacity)
	m.InstanceStartupTimeSeconds = types.Int32Value(int32(s.InstanceStartupTimeSeconds))
	m.InstanceStartupTimeVariance = types.Float32Value(s.InstanceStartupTimeVariance)
	m.InstanceShutdownTimeSeconds = types.Int32Value(int32(s.InstanceShutdownTimeSeconds))
	m.InstanceShutdownTimeVariance = types.Float32Value(s.InstanceShutdownTimeVariance)
	m.UnexpectedInstanceTerminationProbability = types.Float32Value(s.UnexpectedInstanceTerminationProbability)
	m.FailOnRequestAllInstances = types.BoolValue(s.FailOnRequestAllInstances)
	m.FailOnProvision = types.BoolValue(s.FailOnProvision)
	m.ReportSupportingResourcesExist = types.BoolValue(s.ReportSupportingResourcesExist)
	m.Limit = types.Int32Value(int32(s.Limit))
}

func (m *AwsPlacementGroupModel) updateFromSDK(s *model.AwsPlacementGroup) {
	m.GroupName = convert.OptionalStringToTF(s.GroupName)
	m.GroupId = convert.OptionalStringToTF(s.GroupId)
}

func createAwsSecondaryNetworkInterfaceModelFromSDK(s *model.AwsSecondaryNetworkInterface, _ *diag.Diagnostics) AwsSecondaryNetworkInterfaceModel {
	return AwsSecondaryNetworkInterfaceModel{
		NetworkInterfaceType: convert.EnumToTF(s.NetworkInterfaceType),
		SubnetId:             convert.OptionalStringToTF(s.SubnetId),
		SecurityGroupId:      convert.OptionalStringToTF(s.SecurityGroupId),
	}
}

func (m *AwsCapacityReservationModel) updateFromSDK(s *model.AwsCapacityReservation) {
	m.Preference = convert.EnumToTF(s.Preference)
	m.GroupArn = convert.OptionalStringToTF(s.GroupArn)
	m.Id = convert.OptionalStringToTF(s.Id)
}

func createAwsFleetInstanceOverrideModelFromSDK(s *model.AwsFleetInstanceOverride, _ *diag.Diagnostics) AwsFleetInstanceOverrideModel {
	return AwsFleetInstanceOverrideModel{
		AvailabilityZone: convert.OptionalStringToTF(s.AvailabilityZone),
		InstanceType:     types.StringValue(s.InstanceType),
		SpotMaxPrice:     convert.OptionalFloat64ToTF(s.SpotMaxPrice),
		SubnetId:         convert.OptionalStringToTF(s.SubnetId),
	}
}

func (m *AwsFleetOnDemandOptionsModel) updateFromSDK(s *model.AwsFleetOnDemandOptions) {
	m.AllocationStrategy = convert.OptionalEnumToTF(s.AllocationStrategy)
	m.UseCapacityReservationsFirst = convert.OptionalBoolToTF(s.UseCapacityReservationsFirst)
	m.MaxTotalPrice = convert.OptionalFloat64ToTF(s.MaxTotalPrice)
	m.MinInstanceCount = convert.OptionalInt32ToTF(s.MinInstanceCount)
	m.SingleAvailabilityZone = convert.OptionalBoolToTF(s.SingleAvailabilityZone)
	m.SingleInstanceType = convert.OptionalBoolToTF(s.SingleInstanceType)
}

func (m *AwsFleetSpotOptionsModel) updateFromSDK(s *model.AwsFleetSpotOptions) {
	m.AllocationStrategy = convert.OptionalEnumToTF(s.AllocationStrategy)
	m.InstancePoolsToUseCount = convert.OptionalInt32ToTF(s.InstancePoolsToUseCount)
	m.LaunchReplacementInstanceOnRebalance = convert.OptionalBoolToTF(s.LaunchReplacementInstanceOnRebalance)
	m.MaxTotalPrice = convert.OptionalFloat64ToTF(s.MaxTotalPrice)
	m.MinInstanceCount = convert.OptionalInt32ToTF(s.MinInstanceCount)
	m.SingleAvailabilityZone = convert.OptionalBoolToTF(s.SingleAvailabilityZone)
	m.SingleInstanceType = convert.OptionalBoolToTF(s.SingleInstanceType)
}
