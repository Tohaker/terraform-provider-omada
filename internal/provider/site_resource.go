package provider

import (
	"context"
	"fmt"

	"github.com/Tohaker/omada-go-sdk/omada"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &siteResource{}
	_ resource.ResourceWithConfigure = &siteResource{}
)

// NewSiteResource is a helper function to simplify the provider implementation.
func NewSiteResource() resource.Resource {
	return &siteResource{}
}

// siteResource is the resource implementation.
type siteResource struct {
	client   *omada.APIClient
	omadacId string
}

// siteResourceModel maps the resource schema data.
type siteResourceModel struct {
	SiteId               types.String                  `tfsdk:"site_id"`
	Name                 types.String                  `tfsdk:"name"`
	Type                 types.Int32                   `tfsdk:"type"`
	Region               types.String                  `tfsdk:"region"`
	TimeZone             types.String                  `tfsdk:"time_zone"`
	Scenario             types.String                  `tfsdk:"scenario"`
	TagIDs               []types.String                `tfsdk:"tag_ids"`
	Longitude            types.Float64                 `tfsdk:"longitude"`
	Latitude             types.Float64                 `tfsdk:"latitude"`
	Address              types.String                  `tfsdk:"address"`
	DeviceAccountSetting siteDeviceAccountSettingModel `tfsdk:"device_account_setting"`
	SupportES            types.Bool                    `tfsdk:"support_es"`
	SupportL2            types.Bool                    `tfsdk:"support_l2"`
}

// siteDeviceAccountSettingModel maps device account settings data.
type siteDeviceAccountSettingModel struct {
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

// Configure adds the provider configured client to the resource.
func (d *siteResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*providerData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *omada.APIClient, got %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = data.Client
	d.omadacId = data.OmadacId
}

// Metadata returns the resource type name.
func (r *siteResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_site"
}

// Schema defines the schema for the resource.
func (r *siteResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"site_id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"type": schema.Int32Attribute{},
			"region": schema.StringAttribute{
				Required: true,
			},
			"time_zone": schema.StringAttribute{
				Required: true,
			},
			"scenario": schema.StringAttribute{
				Required: true,
			},
			"tag_ids": schema.ListAttribute{
				ElementType: types.StringType,
			},
			"longitude": schema.Float64Attribute{},
			"latitude":  schema.Float64Attribute{},
			"address":   schema.StringAttribute{},
			"device_account_setting": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"username": schema.StringAttribute{
						Required: true,
					},
					"password": schema.StringAttribute{
						Required:  true,
						Sensitive: true,
					},
				},
			},
			"support_es": schema.BoolAttribute{},
			"support_l2": schema.BoolAttribute{},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *siteResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan siteResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var siteEntity omada.CreateSiteEntity
	siteEntity.Name = plan.Name.ValueString()
	siteEntity.Type = plan.Type.ValueInt32Pointer()
	siteEntity.Region = plan.Region.ValueString()
	siteEntity.TimeZone = plan.TimeZone.ValueString()
	siteEntity.Scenario = plan.Scenario.ValueString()
	siteEntity.Longitude = plan.Longitude.ValueFloat64Pointer()
	siteEntity.Latitude = plan.Latitude.ValueFloat64Pointer()
	siteEntity.Address = plan.Address.ValueStringPointer()
	siteEntity.DeviceAccountSetting = omada.DeviceAccountSettingOpenApiVO{
		Username: plan.DeviceAccountSetting.Username.ValueString(),
		Password: plan.DeviceAccountSetting.Password.ValueString(),
	}
	siteEntity.SupportES = plan.SupportES.ValueBoolPointer()
	siteEntity.SupportL2 = plan.SupportL2.ValueBoolPointer()

	var tagIds []string
	for _, tagId := range plan.TagIDs {
		tagIds = append(tagIds, tagId.ValueString())
	}

	siteEntity.TagIds = tagIds

	// Create new site
	site, _, err := r.client.SiteAPI.CreateNewSite(ctx, r.omadacId).CreateSiteEntity(siteEntity).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating site",
			"Could not create site, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema
	if siteId, ok := site.Result["siteId"].(string); ok {
		plan.SiteId = types.StringValue(siteId)
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *siteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *siteResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *siteResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
