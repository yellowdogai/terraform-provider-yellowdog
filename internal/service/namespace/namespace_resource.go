package namespace

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/yellowdogai/yellowdog-sdk-go/account"
	"github.com/yellowdogai/yellowdog-sdk-go/model"
)

func (r *namespaceResource) create(ctx context.Context, state *namespaceModel, _ *diag.Diagnostics) (bool, error) {
	id, err := r.client.Namespaces.CreateNamespace(ctx, &account.CreateNamespaceRequest{
		Namespace: state.Namespace.ValueString(),
	})
	if err != nil {
		return NotCreated, err
	}
	state.Id = types.StringValue(id)
	return Created, nil
}

func (r *namespaceResource) read(ctx context.Context, state *namespaceModel, diags *diag.Diagnostics) error {
	namespace, err := r.client.Namespaces.GetNamespace(ctx, state.Id.ValueString())
	if err != nil {
		return err
	}
	state.updateFromSDK(ctx, namespace, diags)
	return nil
}

func (r *namespaceResource) update(_ context.Context, _ *namespaceModel, _ *diag.Diagnostics) error {
	return fmt.Errorf("update called on Namespace resource, should never happen because all changes require replace")
}

func (r *namespaceResource) delete(ctx context.Context, state namespaceModel, _ *diag.Diagnostics) error {
	return r.client.Namespaces.DeleteNamespace(ctx, state.Id.ValueString())
}

func (n *namespaceModel) updateFromSDK(_ context.Context, namespace *model.Namespace, _ *diag.Diagnostics) {
	n.Id = types.StringValue(namespace.Id)
	n.Namespace = types.StringValue(namespace.Namespace)
}
