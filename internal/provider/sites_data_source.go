package provider

import (
	"context"
	"fmt"

	"github.com/Tohaker/omada-go-sdk/omada"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &sitesDataSource{}
	_ datasource.DataSourceWithConfigure = &sitesDataSource{}
)

// NewSitesDataSource is a helper function to simplify the provider implementation.
func NewSitesDataSource() datasource.DataSource {
	return &sitesDataSource{}
}

// sitesDataSource is the data source implementation.
type sitesDataSource struct {
	client   *omada.APIClient
	omadacId string
}

// sitesDataSourceModel maps the data source schema data.
type sitesDataSourceModel struct {
	Sites []sitesModel `tfsdk:"sites"`
}

// sitesModel maps sites schema data.
type sitesModel struct {
	SiteID       types.String   `tfsdk:"site_id"`
	Name         types.String   `tfsdk:"name"`
	TagIDs       []types.String `tfsdk:"tag_ids"`
	Region       types.String   `tfsdk:"region"`
	TimeZone     types.String   `tfsdk:"time_zone"`
	Scenario     types.String   `tfsdk:"scenario"`
	Longitude    types.Float64  `tfsdk:"longitude"`
	Latitude     types.Float64  `tfsdk:"latitude"`
	Address      types.String   `tfsdk:"address"`
	Type         types.Int32    `tfsdk:"type"`
	SupportES    types.Bool     `tfsdk:"support_es"`
	SupportL2    types.Bool     `tfsdk:"support_l2"`
	SitePublicIP types.String   `tfsdk:"site_public_ip"`
}

// Configure adds the provider configured client to the data source.
func (d *sitesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*providerData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *omada.APIClient, got %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = data.Client
	d.omadacId = data.OmadacId
}

// Metadata returns the data source type name.
func (d *sitesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sites"
}

// Schema defines the schema for the data source.
func (d *sitesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"sites": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"site_id": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"tag_ids": schema.ListAttribute{
							ElementType: types.StringType,
							Computed:    true,
						},
						"region": schema.StringAttribute{
							Computed: true,
						},
						"time_zone": schema.StringAttribute{
							Computed: true,
						},
						"scenario": schema.StringAttribute{
							Computed: true,
						},
						"longitude": schema.Float64Attribute{
							Computed: true,
						},
						"latitude": schema.Float64Attribute{
							Computed: true,
						},
						"address": schema.StringAttribute{
							Computed: true,
						},
						"type": schema.Int32Attribute{
							Computed: true,
						},
						"support_es": schema.BoolAttribute{
							Computed: true,
						},
						"support_l2": schema.BoolAttribute{
							Computed: true,
						},
						"site_public_ip": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *sitesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state sitesDataSourceModel

	response, _, err := d.client.SiteAPI.GetSiteList(ctx, d.omadacId).Page(1).PageSize(1000).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Omada Sites",
			err.Error(),
		)
		return
	}

	sites := response.Result.Data

	for _, site := range sites {
		siteState := sitesModel{
			SiteID:       types.StringPointerValue(site.SiteId),
			Name:         types.StringPointerValue(site.Name),
			Region:       types.StringPointerValue(site.Region),
			TimeZone:     types.StringPointerValue(site.TimeZone),
			Scenario:     types.StringPointerValue(site.Scenario),
			Longitude:    types.Float64PointerValue(site.Longitude),
			Latitude:     types.Float64PointerValue(site.Latitude),
			Address:      types.StringPointerValue(site.Address),
			Type:         types.Int32PointerValue(site.Type),
			SupportES:    types.BoolPointerValue(site.SupportES),
			SupportL2:    types.BoolPointerValue(site.SupportL2),
			SitePublicIP: types.StringPointerValue(site.SitePublicIp),
		}

		for _, tagId := range site.TagIds {
			siteState.TagIDs = append(siteState.TagIDs, types.StringValue(tagId))
		}

		state.Sites = append(state.Sites, siteState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
