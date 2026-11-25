// This resource is not generated, it does not have a Java model class

package account

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/yellowdogai/terraform-provider-yellowdog/internal/resourceerror"
	"github.com/yellowdogai/terraform-provider-yellowdog/internal/validators"
	yellowdog "github.com/yellowdogai/yellowdog-sdk-go"
)

var (
	// Ensure the implementation satisfies the expected interfaces.
	_ resource.ResourceWithConfigure            = &groupApplicationResource{}
	_ resource.Resource                         = &groupApplicationResource{}
	_ ResourceOperations[groupApplicationModel] = &groupApplicationResource{}

	// prevent import errors if imports are not used in generated code
	_ = stringvalidator.NoneOf()
	_ validator.Bool
)

// NewGroupApplicationResource is a helper function to simplify the provider implementation.
func NewGroupApplicationResource() resource.Resource {
	return &groupApplicationResource{}
}

// groupApplicationResource is the resource implementation.
type groupApplicationResource struct {
	client *yellowdog.PlatformClient
}

// groupApplicationModel maps the resource schema data.
type groupApplicationModel struct {
	ApplicationId types.String `tfsdk:"application_id"`
	GroupId       types.String `tfsdk:"group_id"`
}

func (r *groupApplicationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group_application"
}

func (r *groupApplicationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	attributes := map[string]schema.Attribute{
		"application_id": schema.StringAttribute{
			Optional:            false,
			Required:            true,
			MarkdownDescription: "Application that should be a member of the group",
			PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			Validators:          []validator.String{validators.Ydid("app")},
		},
		"group_id": schema.StringAttribute{
			Optional:            false,
			Required:            true,
			MarkdownDescription: "Group the Application should be a member of",
			PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			Validators:          []validator.String{validators.Ydid("group")},
		},
	}
	resp.Schema = schema.Schema{
		Description: "Manages a groupApplicationResource.",
		Attributes:  attributes,
	}
}

func (r *groupApplicationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// nil check when handling ProviderData because Terraform sets it after calling the ConfigureProvider RPC
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*yellowdog.PlatformClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *yellowdog.PlatformClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	r.client = client
}

func (r *groupApplicationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from context
	var state groupApplicationModel
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delegate creation
	created, err := r.create(ctx, &state, &resp.Diagnostics)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create group_application", err.Error())
	}
	if !created {
		return
	}

	// Update state
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *groupApplicationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state groupApplicationModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delegate read
	if err := r.read(ctx, &state, &resp.Diagnostics); err != nil {
		if resourceerror.IsNotFound(err) {
			// See https://developer.hashicorp.com/terraform/plugin/framework/resources/read#recommendations
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Unable to read group_application", err.Error())
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *groupApplicationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var state groupApplicationModel
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delegate update
	if err := r.update(ctx, &state, &resp.Diagnostics); err != nil {
		resp.Diagnostics.AddError("Unable to update group_application", err.Error())
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *groupApplicationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state groupApplicationModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// A resource that is already gone counts as deleted. Any other error keeps it under management.
	if err := r.delete(ctx, state, &resp.Diagnostics); err != nil && !resourceerror.IsNotFound(err) {
		resp.Diagnostics.AddError("Unable to delete group_application", err.Error())
	}
}

func (r *groupApplicationResource) create(ctx context.Context, state *groupApplicationModel, _ *diag.Diagnostics) (bool, error) {
	if err := r.client.Account.AddApplicationToGroup(ctx, state.GroupId.ValueString(), state.ApplicationId.ValueString()); err != nil {
		return NotCreated, err
	}
	return Created, nil
}

func (r *groupApplicationResource) read(ctx context.Context, state *groupApplicationModel, _ *diag.Diagnostics) error {
	for application, err := range r.client.Account.GetGroupApplications(state.GroupId.ValueString()).Iterate(ctx) {
		if err != nil {
			return err
		}
		applicationId, isPresent := application.Id.Get()
		if !isPresent {
			// should never happen
			tflog.Debug(ctx, "application with no id returned from get group applications")
			continue
		}
		if applicationId == state.ApplicationId.ValueString() {
			// the application is in the group, nothing to update
			return nil
		}
	}
	return resourceerror.NewNotFoundError(state.GroupId.ValueString()+"/"+state.ApplicationId.ValueString(), "GroupApplication")
}

func (r *groupApplicationResource) update(_ context.Context, _ *groupApplicationModel, diagnostics *diag.Diagnostics) error {
	diagnostics.AddError("group applications cannot be updated", "there may be a bug in the provider")
	return nil
}

func (r *groupApplicationResource) delete(ctx context.Context, state groupApplicationModel, _ *diag.Diagnostics) error {
	return r.client.Account.RemoveApplicationFromGroup(ctx, state.GroupId.ValueString(), state.ApplicationId.ValueString())
}
