package account

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/yellowdogai/terraform-provider-yellowdog/internal/convert"
	"github.com/yellowdogai/yellowdog-sdk-go/account"
	"github.com/yellowdogai/yellowdog-sdk-go/model"
	"github.com/yellowdogai/yellowdog-sdk-go/optional"
)

func (r *applicationResource) create(ctx context.Context, state *applicationModel, diagnostics *diag.Diagnostics) (bool, error) {

	addReq := &account.AddApplicationRequest{
		Name:        state.Name.ValueString(),
		Description: convert.TFToOptionalString(state.Description),
	}

	addResp, err := r.client.Account.AddApplication(ctx, addReq)
	if err != nil {
		return NotCreated, err
	}
	app, ok := addResp.Application.Get()
	if !ok {
		// may have been created, but we cannot get the id from the response
		return NotCreated, errors.New("the API did not return the created application: application may have been created, but we do not understand the response, check for orphans manually")
	}

	state.updateFromSDK(ctx, &app, diagnostics)
	if key, ok := addResp.ApiKey.Get(); ok {
		state.ApiKey = convert.OptionalStringToTF(key.Id)
		state.ApiSecret = convert.OptionalStringToTF(key.Secret)
	} else {
		// Should never happen, erroring completely would leave the Application in place unless there was a separate call to delete it.
		// Simple solution for now, could be improved later
		diagnostics.AddWarning("Application created without an API key", "The API did not return an API key or secret for the new application, so api_key and api_secret are null.")
	}
	return Created, nil
}

func (r *applicationResource) read(ctx context.Context, state *applicationModel, diags *diag.Diagnostics) error {
	app, err := r.client.Account.GetApplication(ctx, state.Id.ValueString())
	if err != nil {
		return err
	}
	state.updateFromSDK(ctx, app, diags)
	return nil
}

func (r *applicationResource) update(ctx context.Context, state *applicationModel, diags *diag.Diagnostics) error {
	updateReq := &account.UpdateApplicationRequest{
		Name:        optional.Of(state.Name.ValueString()),
		Description: convert.TFToOptionalString(state.Description),
	}

	app, err := r.client.Account.UpdateApplication(ctx, state.Id.ValueString(), updateReq)
	if err != nil {
		return err
	}

	state.updateFromSDK(ctx, app, diags)
	return nil
}

func (r *applicationResource) delete(ctx context.Context, state applicationModel, _ *diag.Diagnostics) error {
	return r.client.Account.DeleteApplication(ctx, state.Id.ValueString())
}

func (s *applicationModel) updateFromSDK(_ context.Context, app *model.Application, _ *diag.Diagnostics) {
	s.Id = types.StringValue(app.Id.Or(""))
	s.Name = types.StringValue(app.Name)
	s.Description = convert.OptionalStringToTF(app.Description)
	s.CreatedByUserId = types.StringValue(app.CreatedByUserId)
}
