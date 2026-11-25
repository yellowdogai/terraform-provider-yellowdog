package keyring

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/yellowdogai/terraform-provider-yellowdog/internal/convert"
	"github.com/yellowdogai/yellowdog-sdk-go/account"
	"github.com/yellowdogai/yellowdog-sdk-go/model"
)

func (r *keyringResource) create(ctx context.Context, state *keyringModel, diags *diag.Diagnostics) (bool, error) {
	addResp, err := r.client.Keyring.AddKeyring(ctx, state.Name.ValueString(), state.Description.ValueString())
	if err != nil {
		return NotCreated, err
	}
	keyring, ok := addResp.Keyring.Get()
	if !ok {
		// Should never happen, we have to error because we cannot get the Id without this
		return NotCreated, errors.New("the API did not return the created keyring: keyring may have been created, but we do not understand the response, check for orphans manually")
	}
	state.updateFromSDK(ctx, &keyring, diags)
	state.KeyringPassword = convert.OptionalStringToTF(addResp.KeyringPassword)
	return Created, nil
}

func (r *keyringResource) read(ctx context.Context, state *keyringModel, diags *diag.Diagnostics) error {
	keyring, err := r.client.Keyring.GetKeyring(ctx, state.Id.ValueString())
	if err != nil {
		return err
	}
	state.updateFromSDK(ctx, keyring, diags)
	return nil
}

func (r *keyringResource) update(ctx context.Context, state *keyringModel, diags *diag.Diagnostics) error {
	keyring, err := r.client.Keyring.UpdateKeyring(ctx, state.Id.ValueString(), &account.UpdateKeyringRequest{
		Description: convert.TFToOptionalString(state.Description),
	})
	if err != nil {
		return err
	}

	state.updateFromSDK(ctx, keyring, diags)
	return nil
}

func (r *keyringResource) delete(ctx context.Context, state keyringModel, _ *diag.Diagnostics) error {
	return r.client.Keyring.DeleteKeyringByName(ctx, state.Name.ValueString())
}

func (k *keyringModel) updateFromSDK(_ context.Context, keyring *model.Keyring, _ *diag.Diagnostics) {
	k.Id = types.StringValue(keyring.Id.Or(""))
	k.Name = types.StringValue(keyring.Name.Or(""))
	k.Description = convert.OptionalStringToTF(keyring.Description)
}
