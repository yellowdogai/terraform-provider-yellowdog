package keyring

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/yellowdogai/terraform-provider-yellowdog/internal/convert"
	"github.com/yellowdogai/yellowdog-sdk-go/model"
)

func (r *googleCloudCredentialResource) create(ctx context.Context, state *googleCloudCredentialModel, diags *diag.Diagnostics) (bool, error) {
	err := putKeyringCredential(ctx, r.client, state.KeyringName, state.toSDK(ctx, diags))
	return err == nil, err
}

func (r *googleCloudCredentialResource) read(ctx context.Context, state *googleCloudCredentialModel, _ *diag.Diagnostics) error {
	return readKeyringCredential(ctx, r.client, state.KeyringName, "GoogleCloudCredential", &state.Name, &state.Description)
}

func (r *googleCloudCredentialResource) update(ctx context.Context, state *googleCloudCredentialModel, diags *diag.Diagnostics) error {
	return putKeyringCredential(ctx, r.client, state.KeyringName, state.toSDK(ctx, diags))
}

func (r *googleCloudCredentialResource) delete(ctx context.Context, state googleCloudCredentialModel, _ *diag.Diagnostics) error {
	return deleteKeyringCredential(ctx, r.client, state.KeyringName, state.Name)
}

func (m *googleCloudCredentialModel) toSDK(_ context.Context, _ *diag.Diagnostics) *model.GoogleCloudCredential {
	credential := model.GoogleCloudCredential{}
	credential.Name = m.Name.ValueString()
	credential.Description = convert.TFToOptionalString(m.Description)
	credential.ServiceAccountKeyJson = m.ServiceAccountKeyJson.ValueString()
	return &credential
}
