package keyring

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/yellowdogai/terraform-provider-yellowdog/internal/convert"
	"github.com/yellowdogai/yellowdog-sdk-go/model"
)

func (r *azureClientCredentialResource) create(ctx context.Context, state *azureClientCredentialModel, diags *diag.Diagnostics) (bool, error) {
	err := putKeyringCredential(ctx, r.client, state.KeyringName, state.toSDK(ctx, diags))
	return err == nil, err
}

func (r *azureClientCredentialResource) read(ctx context.Context, state *azureClientCredentialModel, _ *diag.Diagnostics) error {
	return readKeyringCredential(ctx, r.client, state.KeyringName, "AzureClientCredential", &state.Name, &state.Description)
}

func (r *azureClientCredentialResource) update(ctx context.Context, state *azureClientCredentialModel, diags *diag.Diagnostics) error {
	return putKeyringCredential(ctx, r.client, state.KeyringName, state.toSDK(ctx, diags))
}

func (r *azureClientCredentialResource) delete(ctx context.Context, state azureClientCredentialModel, _ *diag.Diagnostics) error {
	return deleteKeyringCredential(ctx, r.client, state.KeyringName, state.Name)
}

func (m *azureClientCredentialModel) toSDK(_ context.Context, _ *diag.Diagnostics) *model.AzureClientCredential {
	credential := model.AzureClientCredential{}
	credential.Name = m.Name.ValueString()
	credential.Description = convert.TFToOptionalString(m.Description)
	credential.ClientId = m.ClientId.ValueString()
	credential.TenantId = m.TenantId.ValueString()
	credential.SubscriptionId = m.SubscriptionId.ValueString()
	credential.Key = m.Key.ValueString()
	return &credential
}
