package provider

import (
	"context"
	"os"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/yellowdogai/terraform-provider-yellowdog/internal/service/account"
	yellowdog "github.com/yellowdogai/yellowdog-sdk-go"
)

// Ensure YellowdogProvider satisfies the expected interfaces.
var _ provider.Provider = &YellowdogProvider{}

// Environment variables that configure the provider when its configuration leaves the matching attribute unset.
const (
	EnvAPIURL       = "YD_API_URL"
	EnvAPIKeyID     = "YD_API_KEY_ID"
	EnvAPIKeySecret = "YD_API_KEY_SECRET"
)

// YellowdogProvider is the provider implementation.
type YellowdogProvider struct {
	// version is set by main.go at build time, mostly for the User-Agent header.
	version string
}

// yellowdogProviderModel maps the provider configuration block to Terraform types.
type yellowdogProviderModel struct {
	ApiKeyId     types.String `tfsdk:"api_key_id"`
	ApiKeySecret types.String `tfsdk:"api_key_secret"`
	ApiUrl       types.String `tfsdk:"api_url"`
	Timeout      types.String `tfsdk:"timeout"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &YellowdogProvider{version: version}
	}
}

func (p *YellowdogProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "yellowdog"
	resp.Version = p.version
}

func (p *YellowdogProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manage [YellowDog](https://yellowdog.ai/) configuration as code.",
		Description:         "Manage YellowDog configuration as code.",
		Attributes: map[string]schema.Attribute{
			"api_key_id": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "The ID of a YellowDog application's API key. May also be set via the " + EnvAPIKeyID + " environment variable.",
			},
			"api_key_secret": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "The secret of a YellowDog application's API key. May also be set via the " + EnvAPIKeySecret + " environment variable.",
			},
			"api_url": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Base URL for the YellowDog API. May also be set via the " + EnvAPIURL + " environment variable. Defaults to the production API.",
			},
			"timeout": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Default per-request timeout, as a Go duration string (e.g. \"60s\"). Defaults to \"30s\".",
			},
		},
	}
}

func (p *YellowdogProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var model yellowdogProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	config := yellowdog.DefaultConfig(
		p.servicesSchema(model),
		p.apiKey(model),
		applicationProducts(p.version, req.TerraformVersion)...,
	)

	if !model.Timeout.IsNull() {
		d, err := time.ParseDuration(model.Timeout.ValueString())
		if err != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("timeout"),
				"Invalid timeout",
				"timeout must be a valid Go duration string, e.g. \"60s\": "+err.Error(),
			)
			return
		}
		config.Timeout = d
	}

	if config.APIKey.ID == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key_id"),
			"Missing API key",
			"Set api_key_id in the provider configuration or the "+EnvAPIKeyID+" environment variable.",
		)
	}
	if config.APIKey.Secret == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key_secret"),
			"Missing API secret",
			"Set api_key_secret in the provider configuration or the "+EnvAPIKeySecret+" environment variable.",
		)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := yellowdog.NewClient(config)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create YellowDog client", err.Error())
		return
	}

	yc := client
	resp.DataSourceData = yc
	resp.ResourceData = yc
}

func (p *YellowdogProvider) servicesSchema(model yellowdogProviderModel) yellowdog.ServicesSchema {
	defaultUrl := os.Getenv(EnvAPIURL)
	if !model.ApiUrl.IsNull() {
		defaultUrl = model.ApiUrl.ValueString()
	}
	if defaultUrl == "" {
		defaultUrl = "https://api.yellowdog.ai"
	}
	return yellowdog.ServicesSchema{DefaultURL: defaultUrl}
}

func (p *YellowdogProvider) apiKey(model yellowdogProviderModel) yellowdog.APIKey {
	id := os.Getenv(EnvAPIKeyID)
	if !model.ApiKeyId.IsNull() {
		id = model.ApiKeyId.ValueString()
	}
	secret := os.Getenv(EnvAPIKeySecret)
	if !model.ApiKeySecret.IsNull() {
		secret = model.ApiKeySecret.ValueString()
	}
	return yellowdog.APIKey{ID: id, Secret: yellowdog.Secret(secret)}
}

// applicationProducts lists the provider before Terraform, as each product follows the one it uses, and Terraform runs
// the provider. Terraform can leave its version empty, and then is left out rather than sent without a version.
func applicationProducts(providerVersion, terraformVersion string) []yellowdog.Product {
	products := []yellowdog.Product{{Name: "terraform-provider-yellowdog", Version: providerVersion}}
	if terraformVersion != "" {
		products = append(products, yellowdog.Product{Name: "terraform", Version: terraformVersion})
	}
	return products
}

func (p *YellowdogProvider) Resources(ctx context.Context) []func() resource.Resource {
	return append(
		p.generatedResources(ctx),
		account.NewGroupApplicationResource,
	)
}

func (p *YellowdogProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}
