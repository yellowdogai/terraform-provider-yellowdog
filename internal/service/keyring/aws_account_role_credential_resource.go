package keyring

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/yellowdogai/terraform-provider-yellowdog/internal/convert"
	"github.com/yellowdogai/yellowdog-sdk-go/model"
)

func (r *awsAccountRoleCredentialResource) create(ctx context.Context, state *awsAccountRoleCredentialModel, diags *diag.Diagnostics) (bool, error) {
	err := putKeyringCredential(ctx, r.client, state.KeyringName, state.toSDK(ctx, diags))
	return err == nil, err
}

func (r *awsAccountRoleCredentialResource) read(ctx context.Context, state *awsAccountRoleCredentialModel, _ *diag.Diagnostics) error {
	return readKeyringCredential(ctx, r.client, state.KeyringName, "AwsAccountRoleCredential", &state.Name, &state.Description)
}

func (r *awsAccountRoleCredentialResource) update(ctx context.Context, state *awsAccountRoleCredentialModel, diags *diag.Diagnostics) error {
	return putKeyringCredential(ctx, r.client, state.KeyringName, state.toSDK(ctx, diags))
}

func (r *awsAccountRoleCredentialResource) delete(ctx context.Context, state awsAccountRoleCredentialModel, _ *diag.Diagnostics) error {
	return deleteKeyringCredential(ctx, r.client, state.KeyringName, state.Name)
}

func (m *awsAccountRoleCredentialModel) toSDK(_ context.Context, _ *diag.Diagnostics) *model.AwsAccountRoleCredential {
	credential := model.AwsAccountRoleCredential{}
	credential.Name = m.Name.ValueString()
	credential.Description = convert.TFToOptionalString(m.Description)
	credential.ExternalRoleArn = m.ExternalRoleArn.ValueString()
	credential.ExternalId = m.ExternalId.ValueString()
	return &credential
}
