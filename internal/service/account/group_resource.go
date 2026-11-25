package account

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/yellowdogai/terraform-provider-yellowdog/internal/convert"
	"github.com/yellowdogai/yellowdog-sdk-go/account"
	"github.com/yellowdogai/yellowdog-sdk-go/model"
)

func (r *groupResource) create(ctx context.Context, state *groupModel, diags *diag.Diagnostics) (bool, error) {
	group, err := r.client.Account.AddGroup(ctx, &account.AddGroupRequest{
		Name:        state.Name.ValueString(),
		Description: convert.TFToOptionalString(state.Description),
	})
	if err != nil {
		return NotCreated, err
	}
	state.updateFromSDK(ctx, group, diags)
	return Created, nil
}

func (r *groupResource) read(ctx context.Context, state *groupModel, diags *diag.Diagnostics) error {
	group, err := r.client.Account.GetGroup(ctx, state.Id.ValueString())
	if err != nil {
		return err
	}
	state.updateFromSDK(ctx, group, diags)
	return nil
}

func (r *groupResource) update(ctx context.Context, state *groupModel, diags *diag.Diagnostics) error {
	group, err := r.client.Account.UpdateGroup(ctx, state.Id.ValueString(), &account.UpdateGroupRequest{
		Name:        state.Name.ValueString(),
		Description: convert.TFToOptionalString(state.Description),
	})
	if err != nil {
		return err
	}
	state.updateFromSDK(ctx, group, diags)
	return nil
}

func (r *groupResource) delete(ctx context.Context, state groupModel, _ *diag.Diagnostics) error {
	return r.client.Account.DeleteGroup(ctx, state.Id.ValueString())
}

func (g *groupModel) updateFromSDK(_ context.Context, group *model.Group, _ *diag.Diagnostics) {
	g.Id = types.StringValue(group.Id.Or(""))
	g.Name = types.StringValue(group.Name)
	g.Description = convert.OptionalStringToTF(group.Description)
}
