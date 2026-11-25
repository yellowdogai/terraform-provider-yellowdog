package keyring

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/yellowdogai/terraform-provider-yellowdog/internal/convert"
	"github.com/yellowdogai/yellowdog-sdk-go/model"
)

func (r *ociCredentialResource) create(ctx context.Context, state *ociCredentialModel, diags *diag.Diagnostics) (bool, error) {
	err := putKeyringCredential(ctx, r.client, state.KeyringName, state.toSDK(ctx, diags))
	return err == nil, err
}

func (r *ociCredentialResource) read(ctx context.Context, state *ociCredentialModel, _ *diag.Diagnostics) error {
	return readKeyringCredential(ctx, r.client, state.KeyringName, "OciCredential", &state.Name, &state.Description)
}

func (r *ociCredentialResource) update(ctx context.Context, state *ociCredentialModel, diags *diag.Diagnostics) error {
	return putKeyringCredential(ctx, r.client, state.KeyringName, state.toSDK(ctx, diags))
}

func (r *ociCredentialResource) delete(ctx context.Context, state ociCredentialModel, _ *diag.Diagnostics) error {
	return deleteKeyringCredential(ctx, r.client, state.KeyringName, state.Name)
}

func (m *ociCredentialModel) toSDK(_ context.Context, _ *diag.Diagnostics) *model.OciCredential {
	credential := model.OciCredential{}
	credential.Name = m.Name.ValueString()
	credential.Description = convert.TFToOptionalString(m.Description)
	credential.UserId = m.UserId.ValueString()
	credential.TenantId = m.TenantId.ValueString()
	credential.Fingerprint = m.Fingerprint.ValueString()
	credential.PrivateKey = m.PrivateKey.ValueString()
	credential.Passphrase = convert.TFToOptionalString(m.Passphrase)
	return &credential
}
