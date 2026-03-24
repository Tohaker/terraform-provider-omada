package provider

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"

	omada "github.com/Tohaker/omada-go-sdk/omada"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &omadaProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &omadaProvider{
			version: version,
		}
	}
}

// omadaProvider is the provider implementation.
type omadaProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// omadaProviderModel maps provider schema data to a Go type.
type omadaProviderModel struct {
	Host         types.String `tfsdk:"host"`
	CustomerId   types.String `tfsdk:"customer_id"`
	ClientId     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
}

type providerData struct {
	Client   *omada.APIClient
	OmadacId string
}

// Metadata returns the provider type name.
func (p *omadaProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "omada"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *omadaProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Optional: true,
			},
			"customer_id": schema.StringAttribute{
				Optional: true,
			},
			"client_id": schema.StringAttribute{
				Optional: true,
			},
			"client_secret": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

// Configure prepares a Omada API client for data sources and resources.
func (p *omadaProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Omada client")

	// Retrieve provider data from configuration
	var config omadaProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown Omada API Host",
			"The provider cannot create the Omada API client as there is an unknown configuration value for the Omada API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the OMADA_HOST environment variable.",
		)
	}

	if config.CustomerId.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("customer_id"),
			"Unknown Omada Customer or MSP ID",
			"The provider cannot create the Omada API client as there is an unknown configuration value for the Omada Customer or MSP ID. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the OMADA_CUSTOMER_ID environment variable.",
		)
	}

	if config.ClientId.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_id"),
			"Unknown Omada API Client ID",
			"The provider cannot create the Omada API client as there is an unknown configuration value for the Omada API Client ID. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the OMADA_CLIENT_ID environment variable.",
		)
	}

	if config.ClientSecret.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_secret"),
			"Unknown Omada API Client Secret",
			"The provider cannot create the Omada API client as there is an unknown configuration value for the Omada API Client Secret. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the OMADA_CLIENT_SECRET environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	host := os.Getenv("OMADA_HOST")
	customer_id := os.Getenv("OMADA_CUSTOMER_ID")
	client_id := os.Getenv("OMADA_CLIENT_ID")
	client_secret := os.Getenv("OMADA_CLIENT_SECRET")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.CustomerId.IsNull() {
		customer_id = config.CustomerId.ValueString()
	}

	if !config.ClientId.IsNull() {
		client_id = config.ClientId.ValueString()
	}

	if !config.ClientSecret.IsNull() {
		client_secret = config.ClientSecret.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing Omada API Host",
			"The provider cannot create the Omada API client as there is a missing or empty value for the Omada API host. "+
				"Set the host value in the configuration or use the OMADA_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if customer_id == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("customer_id"),
			"Missing Omada Customer or MSP ID",
			"The provider cannot create the Omada API client as there is a missing or empty value for the Omada Customer or MSP ID. "+
				"Set the customer_id value in the configuration or use the OMADA_CUSTOMER_ID environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if client_id == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_id"),
			"Missing Omada API Client ID",
			"The provider cannot create the Omada API client as there is a missing or empty value for the Omada API Client ID. "+
				"Set the client_id value in the configuration or use the OMADA_CLIENT_ID environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if client_secret == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_secret"),
			"Missing Omada API Client Secret",
			"The provider cannot create the Omada API client as there is a missing or empty value for the Omada API Client Secret. "+
				"Set the client_secret value in the configuration or use the OMADA_CLIENT_SECRET environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "omada_host", host)
	ctx = tflog.SetField(ctx, "omada_customer_id", customer_id)
	ctx = tflog.SetField(ctx, "omada_client_id", client_id)
	ctx = tflog.SetField(ctx, "omada_client_secret", client_secret)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "omada_client_secret")

	tflog.Debug(ctx, "Creating Omada client")

	// Create a new Omada client using the configuration values
	cfg := omada.NewConfiguration()
	cfg.Servers = omada.ServerConfigurations{
		{URL: host},
	}

	// TODO: Remove this when testing environment has SSL certification configured
	cfg.HTTPClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	client := omada.NewAPIClient(cfg)

	tokenResp, _, err := client.AuthorizeAPI.AuthorizeToken(context.Background()).GrantType("client_credentials").TokenRequest(omada.TokenRequest{
		ClientId:     client_id,
		ClientSecret: client_secret,
		OmadacId:     &customer_id,
	}).Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Omada API Client",
			"An unexpected error occurred when creating the Omada API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Omada Client Error: "+err.Error(),
		)
		return
	}

	cfg.DefaultHeader["Authorization"] = "AccessToken=" + *tokenResp.Result.AccessToken

	data := &providerData{
		Client:   client,
		OmadacId: customer_id,
	}

	resp.DataSourceData = data
	resp.ResourceData = data

	tflog.Info(ctx, "Configured Omada client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *omadaProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewSitesDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *omadaProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}
