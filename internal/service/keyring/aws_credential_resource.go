package keyring

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/yellowdogai/terraform-provider-yellowdog/internal/convert"
	"github.com/yellowdogai/yellowdog-sdk-go/model"
)

func (r *awsCredentialResource) create(ctx context.Context, state *awsCredentialModel, diags *diag.Diagnostics) (bool, error) {
	err := putKeyringCredential(ctx, r.client, state.KeyringName, state.toSDK(ctx, diags))
	return err == nil, err
}

func (r *awsCredentialResource) read(ctx context.Context, state *awsCredentialModel, _ *diag.Diagnostics) error {
	return readKeyringCredential(ctx, r.client, state.KeyringName, "AwsCredential", &state.Name, &state.Description)
}

func (r *awsCredentialResource) update(ctx context.Context, state *awsCredentialModel, diags *diag.Diagnostics) error {
	return putKeyringCredential(ctx, r.client, state.KeyringName, state.toSDK(ctx, diags))
}

func (r *awsCredentialResource) delete(ctx context.Context, state awsCredentialModel, _ *diag.Diagnostics) error {
	return deleteKeyringCredential(ctx, r.client, state.KeyringName, state.Name)
}

func (m *awsCredentialModel) toSDK(_ context.Context, _ *diag.Diagnostics) *model.AwsCredential {
	credential := model.AwsCredential{}
	credential.Name = m.Name.ValueString()
	credential.Description = convert.TFToOptionalString(m.Description)
	credential.AccessKeyId = m.AccessKeyId.ValueString()
	credential.SecretAccessKey = m.SecretAccessKey.ValueString()
	return &credential
}
