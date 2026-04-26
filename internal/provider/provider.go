package provider

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"strconv"
	"terraform-provider-omada/internal/client"
	"terraform-provider-omada/internal/service/site"

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
	Host          types.String `tfsdk:"host"`
	ControllerId  types.String `tfsdk:"controller_id"`
	ClientId      types.String `tfsdk:"client_id"`
	ClientSecret  types.String `tfsdk:"client_secret"`
	TlsSkipVerify types.Bool   `tfsdk:"tls_skip_verify"`
}

// Metadata returns the provider type name.
func (p *omadaProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "omada"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *omadaProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Terraform Provider for managing a TP-Link Omada Software Controller. Supports v5.x and v6.x.",
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Description: "URI for the Omada Controller API. May also be provided via `OMADA_HOST` environment variable.",
				Optional:    true,
			},
			"controller_id": schema.StringAttribute{
				Description: "Unique ID assigned to the Omada Controller. May also be provided via `OMADA_CONTROLLER_ID` environment variable.",
				Optional:    true,
			},
			"client_id": schema.StringAttribute{
				Description: "Client ID for the Omada Controller Application. May also be provided via `OMADA_CLIENT_ID` environment variable.",
				Optional:    true,
			},
			"client_secret": schema.StringAttribute{
				Description: "Client Secret for the Omada Controller Application. May also be provided via `OMADA_CLIENT_SECRET` environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
			"tls_skip_verify": schema.BoolAttribute{
				Description: `When set to true, accepts any certificate presented by the server and any host name in that certificate.
				**It is unadvisable to use this in a production environment, as it makes the provider susceptible to man-in-the-middle attacks.**`,
				Optional: true,
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

	if config.ControllerId.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("controller_id"),
			"Unknown Omada Controller ID",
			"The provider cannot create the Omada API client as there is an unknown configuration value for the Omada Controller ID. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the OMADA_CONTROLLER_ID environment variable.",
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

	if config.TlsSkipVerify.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("tls_skip_verify"),
			"Unknown value for tls_skip_verify",
			"The provider cannot create the Omada API client as there is an unknown configuration value for the TLS verification. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the OMADA_TLS_SKIP_VERIFY environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	host := os.Getenv("OMADA_HOST")
	controller_id := os.Getenv("OMADA_CONTROLLER_ID")
	client_id := os.Getenv("OMADA_CLIENT_ID")
	client_secret := os.Getenv("OMADA_CLIENT_SECRET")

	// Read and parse the env variable as a boolean.
	tls_skip_verify, parse_err := strconv.ParseBool(os.Getenv("OMADA_TLS_SKIP_VERIFY"))

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.ControllerId.IsNull() {
		controller_id = config.ControllerId.ValueString()
	}

	if !config.ClientId.IsNull() {
		client_id = config.ClientId.ValueString()
	}

	if !config.ClientSecret.IsNull() {
		client_secret = config.ClientSecret.ValueString()
	}

	if !config.TlsSkipVerify.IsNull() {
		tls_skip_verify = config.TlsSkipVerify.ValueBool()
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

	if controller_id == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("controller_id"),
			"Missing Omada Controller ID",
			"The provider cannot create the Omada API client as there is a missing or empty value for the Omada Controller ID. "+
				"Set the controller_id value in the configuration or use the OMADA_CONTROLLER_ID environment variable. "+
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

	// We don't need to add an error, we will just consider the tls_skip_verify as false
	if parse_err != nil {
		tls_skip_verify = false
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "omada_host", host)
	ctx = tflog.SetField(ctx, "omada_controller_id", controller_id)
	ctx = tflog.SetField(ctx, "omada_client_id", client_id)
	ctx = tflog.SetField(ctx, "omada_client_secret", client_secret)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "omada_client_secret")

	tflog.Debug(ctx, "Creating Omada client")

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: tls_skip_verify},
		},
	}

	meta, err := client.New(ctx, client.Config{
		ClientID:     client_id,
		ClientSecret: client_secret,
		ControllerID: controller_id,
		Host:         host,
		HTTPClient:   httpClient,
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Omada API Client",
			"An unexpected error occurred when creating the Omada API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Omada Client Error: "+err.Error(),
		)
		return
	}

	resp.DataSourceData = meta
	resp.ResourceData = meta

	tflog.Info(ctx, "Configured Omada client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *omadaProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		site.NewDataSourceList,
	}
}

// Resources defines the resources implemented in the provider.
func (p *omadaProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		site.NewResource,
	}
}
