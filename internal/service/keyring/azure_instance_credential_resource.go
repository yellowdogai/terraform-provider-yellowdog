package keyring

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/yellowdogai/terraform-provider-yellowdog/internal/convert"
	"github.com/yellowdogai/yellowdog-sdk-go/model"
)

func (r *azureInstanceCredentialResource) create(ctx context.Context, state *azureInstanceCredentialModel, diags *diag.Diagnostics) (bool, error) {
	err := putKeyringCredential(ctx, r.client, state.KeyringName, state.toSDK(ctx, diags))
	return err == nil, err
}

func (r *azureInstanceCredentialResource) read(ctx context.Context, state *azureInstanceCredentialModel, _ *diag.Diagnostics) error {
	return readKeyringCredential(ctx, r.client, state.KeyringName, "AzureInstanceCredential", &state.Name, &state.Description)
}

func (r *azureInstanceCredentialResource) update(ctx context.Context, state *azureInstanceCredentialModel, diags *diag.Diagnostics) error {
	return putKeyringCredential(ctx, r.client, state.KeyringName, state.toSDK(ctx, diags))
}

func (r *azureInstanceCredentialResource) delete(ctx context.Context, state azureInstanceCredentialModel, _ *diag.Diagnostics) error {
	return deleteKeyringCredential(ctx, r.client, state.KeyringName, state.Name)
}

func (m *azureInstanceCredentialModel) toSDK(_ context.Context, _ *diag.Diagnostics) *model.AzureInstanceCredential {
	credential := model.AzureInstanceCredential{}
	credential.Name = m.Name.ValueString()
	credential.Description = convert.TFToOptionalString(m.Description)
	credential.AdminUsername = m.AdminUsername.ValueString()
	credential.AdminPassword = convert.TFToOptionalString(m.AdminPassword)
	return &credential
}
