package keyring

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/yellowdogai/terraform-provider-yellowdog/internal/convert"
	"github.com/yellowdogai/yellowdog-sdk-go/model"
)

func (r *azureStorageCredentialResource) create(ctx context.Context, state *azureStorageCredentialModel, diags *diag.Diagnostics) (bool, error) {
	err := putKeyringCredential(ctx, r.client, state.KeyringName, state.toSDK(ctx, diags))
	return err == nil, err
}

func (r *azureStorageCredentialResource) read(ctx context.Context, state *azureStorageCredentialModel, _ *diag.Diagnostics) error {
	return readKeyringCredential(ctx, r.client, state.KeyringName, "AzureStorageCredential", &state.Name, &state.Description)
}

func (r *azureStorageCredentialResource) update(ctx context.Context, state *azureStorageCredentialModel, diags *diag.Diagnostics) error {
	return putKeyringCredential(ctx, r.client, state.KeyringName, state.toSDK(ctx, diags))
}

func (r *azureStorageCredentialResource) delete(ctx context.Context, state azureStorageCredentialModel, _ *diag.Diagnostics) error {
	return deleteKeyringCredential(ctx, r.client, state.KeyringName, state.Name)
}

func (m *azureStorageCredentialModel) toSDK(_ context.Context, _ *diag.Diagnostics) *model.AzureStorageCredential {
	credential := model.AzureStorageCredential{}
	credential.Name = m.Name.ValueString()
	credential.Description = convert.TFToOptionalString(m.Description)
	credential.AccountName = m.AccountName.ValueString()
	credential.AccountKey = m.AccountKey.ValueString()
	return &credential
}
