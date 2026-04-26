package site

import (
	"github.com/Tohaker/omada-go-sdk/omada"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type siteClient struct {
	client   *omada.APIClient
	omadacId string
}

type siteCommonModel struct {
	SiteId    types.String   `tfsdk:"site_id"`
	Name      types.String   `tfsdk:"name"`
	Type      types.Int32    `tfsdk:"type"`
	Region    types.String   `tfsdk:"region"`
	TimeZone  types.String   `tfsdk:"time_zone"`
	Scenario  types.String   `tfsdk:"scenario"`
	TagIDs    []types.String `tfsdk:"tag_ids"`
	Longitude types.Float64  `tfsdk:"longitude"`
	Latitude  types.Float64  `tfsdk:"latitude"`
	Address   types.String   `tfsdk:"address"`
	SupportES types.Bool     `tfsdk:"support_es"`
	SupportL2 types.Bool     `tfsdk:"support_l2"`
}

// siteDeviceAccountSettingModel maps device account settings data.
type siteDeviceAccountSettingModel struct {
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

// siteResourceModel maps the resource schema data.
type siteResourceModel struct {
	siteCommonModel
	DeviceAccountSetting *siteDeviceAccountSettingModel `tfsdk:"device_account_setting"`
}

// siteModel maps sites schema data.
type siteModel struct {
	siteCommonModel
	SitePublicIP types.String `tfsdk:"site_public_ip"`
}

// sitesDataSourceModel maps the data source schema data.
type sitesDataSourceModel struct {
	Sites []siteModel `tfsdk:"sites"`
}
