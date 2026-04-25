package site

import (
	"context"
	"fmt"
	"terraform-provider-omada/internal/client"

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

// NewDataSourceList is a helper function to simplify the provider implementation.
func NewDataSourceList() datasource.DataSource {
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

	data, ok := req.ProviderData.(*client.Meta)
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
		Description: "Fetches a list of sites. Only the first 1000 sites are returned.",
		Attributes: map[string]schema.Attribute{
			"sites": schema.ListNestedAttribute{
				Description: "The list of site summaries. Up to 1000 will be returned.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"site_id": schema.StringAttribute{
							Description: "Site ID",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the site. Will contain 1 to 64 characters.",
							Computed:    true,
						},
						"tag_ids": schema.ListAttribute{
							Description: "List of site tag ids.",
							ElementType: types.StringType,
							Computed:    true,
						},
						"region": schema.StringAttribute{
							Description: "The Country/Region of the site. For the possible values of `region`, refer to the abbreviation of the ISO country code; For example, \"United States\" refers to the United States of America.",
							Computed:    true,
						},
						"time_zone": schema.StringAttribute{
							Description: "Time zone of the site. For possible values, refer to section 5.1 of the [Open API Access Guide](https://use1-omada-northbound.tplinkcloud.com/doc.html#/home).",
							Computed:    true,
						},
						"scenario": schema.StringAttribute{
							Description: "Scenario in which the site is deployed. For the values of the scenario of the site, refer to result of the interface for [Get scenario list](https://use1-omada-northbound.tplinkcloud.com/doc.html#/00%20All/Site/getScenarioList).",
							Computed:    true,
						},
						"longitude": schema.Float64Attribute{
							Description: "Longitude of the site. Will be within the range of -180 - 180.",
							Computed:    true,
						},
						"latitude": schema.Float64Attribute{
							Description: "Latitude of the site. Will be within the range of -90 - 90.",
							Computed:    true,
						},
						"address": schema.StringAttribute{
							Description: "Address of the site.",
							Computed:    true,
						},
						"type": schema.Int32Attribute{
							Description: "Type of the site (0 or 1).\n 0 means a Basic site, 1 means a Pro site.",
							Computed:    true,
						},
						"support_es": schema.BoolAttribute{
							Description: "Whether the site supports adopting Agile Series Switches.",
							Computed:    true,
						},
						"support_l2": schema.BoolAttribute{
							Description: "Whether the site supports adopting Non-Agile Series Switches.",
							Computed:    true,
						},
						"site_public_ip": schema.StringAttribute{
							Description: "Adopted gateway public IP of the site, only useful for cloud based controllers and remote management local Controllers.",
							Computed:    true,
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
