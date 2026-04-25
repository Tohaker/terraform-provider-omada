package site

import (
	"context"
	"fmt"
	"terraform-provider-omada/internal/client"

	"github.com/Tohaker/omada-go-sdk/omada"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &siteResource{}
	_ resource.ResourceWithConfigure   = &siteResource{}
	_ resource.ResourceWithImportState = &siteResource{}
)

// NewResource is a helper function to simplify the provider implementation.
func NewResource() resource.Resource {
	return &siteResource{}
}

// siteResource is the resource implementation.
type siteResource struct {
	siteClient
}

// Configure adds the provider configured client to the resource.
func (d *siteResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*client.Meta)
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
		Description: "Manages a site. Your credentials must have the `Global Dashboard Manager Modify` permission.",
		Attributes: map[string]schema.Attribute{
			"site_id": schema.StringAttribute{
				Description: "Site ID",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the site. This must contain 1 to 64 characters and cannot be the same as any existing site.",
				Required:    true,
			},
			"type": schema.Int32Attribute{
				Description: "Type of the site (0 or 1).\n 0 means a Basic site, 1 means a Pro site.",
				Optional:    true,
				Computed:    true,
			},
			"region": schema.StringAttribute{
				Description: "The Country/Region of the site; For the values of `region`, refer to the abbreviation of the ISO country code; For example, \"United States\" refers to the United States of America.",
				Required:    true,
			},
			"time_zone": schema.StringAttribute{
				Description: "Time zone of the site. For possible values, refer to section 5.1 of the [Open API Access Guide](https://use1-omada-northbound.tplinkcloud.com/doc.html#/home).",
				Required:    true,
			},
			"scenario": schema.StringAttribute{
				Description: "Scenario in which the site is deployed. For the values of the scenario of the site, refer to result of the interface for [Get scenario list](https://use1-omada-northbound.tplinkcloud.com/doc.html#/00%20All/Site/getScenarioList).",
				Required:    true,
			},
			"tag_ids": schema.ListAttribute{
				Description: "List of site tag ids.",
				ElementType: types.StringType,
				Optional:    true,
			},
			"longitude": schema.Float64Attribute{
				Description: "Longitude of the site. Must be within the range of -180 - 180.",
				Optional:    true,
			},
			"latitude": schema.Float64Attribute{
				Description: "Latitude of the site. Must be within the range of -90 - 90.",
				Optional:    true,
			},
			"address": schema.StringAttribute{
				Description: "Address of the site.",
				Optional:    true,
			},
			"device_account_setting": schema.SingleNestedAttribute{
				Description: "Login information for devices.",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"username": schema.StringAttribute{
						Description: "Device account username. Must contain 1 to 64 ASCII characters.",
						Required:    true,
					},
					"password": schema.StringAttribute{
						Description: "Device account password. Must contain 10 to 64 ASCII characters.\nPasswords must be a combination of uppercase letters, lowercase letters, numbers, and special symbols. Symbols such as `!`, `#`, `$`, `%`, `&`, `*`, `@` and `^` are supported.\nThe password should not contain consecutive identical characters.\nUsername and Password should not be the same.",
						Required:    true,
						Sensitive:   true,
					},
				},
			},
			"support_es": schema.BoolAttribute{
				Description: "Whether the site supports adopting Agile Series Switches.",
				Optional:    true,
				Computed:    true,
			},
			"support_l2": schema.BoolAttribute{
				Description: "Whether the site supports adopting Non-Agile Series Switches.",
				Optional:    true,
				Computed:    true,
			},
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

	if plan.DeviceAccountSetting == nil {
		resp.Diagnostics.AddError(
			"Missing device_account_setting",
			"device_account_setting must be provided.",
		)
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

	if *site.ErrorCode != 0 {
		resp.Diagnostics.AddError(
			"Error creating site",
			fmt.Sprintf("Could not create site, error code %d: %s", *site.ErrorCode, *site.Msg),
		)
		return
	}

	// Map response body to schema
	plan.SiteId = types.StringPointerValue(site.Result.SiteId)

	updatedSite, _, err := r.client.SiteAPI.GetSiteEntity(ctx, r.omadacId, plan.SiteId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Site Entry",
			"Could not read Omada site ID "+plan.SiteId.ValueString()+": "+err.Error(),
		)
		return
	}

	// Ensure non-nullable properties without are set to their computed value in the plan
	plan.Type = types.Int32PointerValue(updatedSite.Result.Type)
	plan.SupportES = types.BoolPointerValue(updatedSite.Result.SupportES)
	plan.SupportL2 = types.BoolPointerValue(updatedSite.Result.SupportL2)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *siteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state siteResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed site from API
	site, _, err := r.client.SiteAPI.GetSiteEntity(ctx, r.omadacId, state.SiteId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Site Entity",
			"Could not read Omada site ID "+state.SiteId.ValueString()+": "+err.Error(),
		)
		return
	}

	// Get refreshed site device account setting from API
	deviceAccount, _, deviceAccountErr := r.client.SiteAPI.GetSiteDeviceAccountSetting(ctx, r.omadacId, state.SiteId.ValueString()).Execute()
	if deviceAccountErr != nil {
		resp.Diagnostics.AddError(
			"Error reading Device Account",
			"Could not read Device Account for site ID "+state.SiteId.ValueString()+": "+deviceAccountErr.Error(),
		)
		return
	}

	// Overwrite site with refreshed state
	flattenSiteEntity(&state, site.Result)
	state.DeviceAccountSetting = flattenDeviceAccountSettings(deviceAccount.Result)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *siteResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to site_id attribute
	SiteId := path.Root("site_id")

	resource.ImportStatePassthroughID(ctx, SiteId, req, resp)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *siteResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan siteResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.DeviceAccountSetting == nil {
		resp.Diagnostics.AddError(
			"Missing device_account_setting",
			"device_account_setting must be provided.",
		)
		return
	}

	// Generate API request body from plan
	var siteEntity omada.UpdateSiteEntity
	siteEntity.Name = plan.Name.ValueStringPointer()
	siteEntity.Region = plan.Region.ValueString()
	siteEntity.TimeZone = plan.TimeZone.ValueString()
	siteEntity.Scenario = plan.Scenario.ValueString()
	siteEntity.Longitude = plan.Longitude.ValueFloat64Pointer()
	siteEntity.Latitude = plan.Latitude.ValueFloat64Pointer()
	siteEntity.Address = plan.Address.ValueStringPointer()
	siteEntity.SupportES = plan.SupportES.ValueBoolPointer()
	siteEntity.SupportL2 = plan.SupportL2.ValueBoolPointer()

	var tagIds []string
	for _, tagId := range plan.TagIDs {
		tagIds = append(tagIds, tagId.ValueString())
	}

	siteEntity.TagIds = tagIds

	var siteDeviceAccountSetting omada.DeviceAccountSettingOpenApiVO
	siteDeviceAccountSetting.Username = plan.DeviceAccountSetting.Username.ValueString()
	siteDeviceAccountSetting.Password = plan.DeviceAccountSetting.Password.ValueString()

	// Update existing site
	site, _, err := r.client.SiteAPI.UpdateSiteEntity(ctx, r.omadacId, plan.SiteId.ValueString()).UpdateSiteEntity(siteEntity).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating site",
			"Could not update site, unexpected error: "+err.Error(),
		)
		return
	}

	if *site.ErrorCode != 0 {
		resp.Diagnostics.AddError(
			"Error updating site",
			fmt.Sprintf("Could not update site, error code %d: %s", *site.ErrorCode, *site.Msg),
		)
		return
	}

	// Update device account
	deviceAccount, _, deviceUpdateErr := r.client.SiteAPI.UpdateSiteDeviceAccountSetting(ctx, r.omadacId, plan.SiteId.ValueString()).DeviceAccountSettingOpenApiVO(siteDeviceAccountSetting).Execute()

	if deviceUpdateErr != nil {
		resp.Diagnostics.AddError(
			"Error updating site device account",
			"Could not update site device account, unexpected error: "+deviceUpdateErr.Error(),
		)
		return
	}

	if *deviceAccount.ErrorCode != 0 {
		resp.Diagnostics.AddError(
			"Error updating site device account",
			fmt.Sprintf("Could not update site device account, error code %d: %s", *deviceAccount.ErrorCode, *deviceAccount.Msg),
		)
		return
	}

	updatedSite, _, err := r.client.SiteAPI.GetSiteEntity(ctx, r.omadacId, plan.SiteId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Site Entry",
			"Could not read Omada site ID "+plan.SiteId.ValueString()+": "+err.Error(),
		)
		return
	}

	// Ensure non-nullable properties without are set to their computed value in the plan
	plan.Type = types.Int32PointerValue(updatedSite.Result.Type)
	plan.SupportES = types.BoolPointerValue(updatedSite.Result.SupportES)
	plan.SupportL2 = types.BoolPointerValue(updatedSite.Result.SupportL2)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *siteResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state siteResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing site
	_, _, err := r.client.SiteAPI.DeleteSite(ctx, r.omadacId, state.SiteId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting site",
			"Could not delete site, unexpected error: "+err.Error(),
		)
		return
	}
}
