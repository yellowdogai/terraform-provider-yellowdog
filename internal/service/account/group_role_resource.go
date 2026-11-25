package account

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/yellowdogai/terraform-provider-yellowdog/internal/convert"
	"github.com/yellowdogai/terraform-provider-yellowdog/internal/resourceerror"
	"github.com/yellowdogai/yellowdog-sdk-go/account"
	"github.com/yellowdogai/yellowdog-sdk-go/model"
)

func (r *groupRoleResource) create(ctx context.Context, state *groupRoleModel, _ *diag.Diagnostics) (bool, error) {
	if err := r.client.Account.AddRoleToGroup(ctx, state.GroupId.ValueString(), state.RoleId.ValueString(), state.toSDKRoleScope()); err != nil {
		return NotCreated, err
	}
	return Created, nil
}

func (r *groupRoleResource) read(ctx context.Context, state *groupRoleModel, _ *diag.Diagnostics) error {
	group, err := r.client.Account.GetGroup(ctx, state.GroupId.ValueString())
	if err != nil {
		return err
	}
	return state.updateFromSDK(group)
}

func getGroupRole(group *model.Group, state *groupRoleModel) *model.GroupRole {
	for _, role := range group.Roles {
		if role.Role.Id.Or("") == state.RoleId.ValueString() {
			return &role
		}
	}
	return nil
}

func (r *groupRoleResource) update(_ context.Context, _ *groupRoleModel, diagnostics *diag.Diagnostics) error {
	diagnostics.AddError("group roles cannot be updated", "there may be a bug in the provider")
	return nil
}

func (r *groupRoleResource) delete(ctx context.Context, state groupRoleModel, _ *diag.Diagnostics) error {
	return r.client.Account.RemoveRoleFromGroup(ctx, state.GroupId.ValueString(), state.RoleId.ValueString())
}

func (m *groupRoleModel) toSDKRoleScope() *account.RoleScope {
	if len(m.Namespaces.Elements()) == 0 {
		return &account.RoleScope{
			Global: true,
		}
	}
	return &account.RoleScope{
		Namespaces: convert.TFSetToStringSlice(m.Namespaces),
	}
}

func (m *groupRoleModel) updateFromSDK(group *model.Group) error {
	groupRole := getGroupRole(group, m)
	if groupRole == nil {
		return resourceerror.NewNotFoundError(m.GroupId.ValueString()+"/"+m.RoleId.ValueString(), "GroupRole")
	}
	if !groupRole.Scope.Global {
		var namespaces []string
		for _, namespaceRef := range groupRole.Scope.Namespaces {
			namespaces = append(namespaces, namespaceRef.Namespace)
		}
		m.Namespaces = convert.StringSliceToTFSet(namespaces)
	} else {
		m.Namespaces = types.SetNull(types.StringType)
	}
	return nil
}
