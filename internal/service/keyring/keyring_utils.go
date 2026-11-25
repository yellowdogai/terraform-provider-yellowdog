package keyring

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/yellowdogai/terraform-provider-yellowdog/internal/convert"
	"github.com/yellowdogai/terraform-provider-yellowdog/internal/resourceerror"
	yellowdog "github.com/yellowdogai/yellowdog-sdk-go"
	"github.com/yellowdogai/yellowdog-sdk-go/model"
)

// putKeyringCredential creates or replaces the credential in the named keyring.
func putKeyringCredential(ctx context.Context, client *yellowdog.PlatformClient, keyringName types.String, credential model.Credential) error {
	_, err := client.Keyring.PutCredentialByName(ctx, keyringName.ValueString(), credential)
	return err
}

// readKeyringCredential reads name and description from the client into the name and description pointers.
func readKeyringCredential(ctx context.Context, client *yellowdog.PlatformClient, keyringName types.String, resourceType string, name *types.String, description *types.String) error {
	keyring, err := client.Keyring.GetKeyring(ctx, keyringName.ValueString())
	if err != nil {
		return err
	}
	credential, found := getKeyringCredential(name.ValueString(), keyring)
	if !found {
		return resourceerror.NewNotFoundError(name.ValueString(), resourceType)
	}
	*name = convert.OptionalStringToTF(credential.Name)
	*description = convert.OptionalStringToTF(credential.Description)
	return nil
}

// deleteKeyringCredential deletes the credential from the keyring.
func deleteKeyringCredential(ctx context.Context, client *yellowdog.PlatformClient, keyringName types.String, name types.String) error {
	_, err := client.Keyring.DeleteCredentialByName(ctx, keyringName.ValueString(), name.ValueString())
	return err
}

// TODO move this to the yd-go-sdk
func getKeyringCredential(credentialName string, keyring *model.Keyring) (model.KeyringCredential, bool) {
	for _, k := range keyring.Credentials {
		if k.Name.Or("") == credentialName {
			return k, true
		}
	}
	return model.KeyringCredential{}, false
}
