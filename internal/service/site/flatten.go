package site

import (
	"github.com/Tohaker/omada-go-sdk/omada"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func flattenTagIds(m *[]types.String, r *[]string) {
	for _, tagId := range *r {
		*m = append(*m, types.StringValue(tagId))
	}
}

func flattenDeviceAccountSettings(r *omada.DeviceAccountSettingOpenApiVO) *siteDeviceAccountSettingModel {
	if r == nil {
		return nil
	}

	return &siteDeviceAccountSettingModel{
		Username: types.StringValue(r.Username),
		Password: types.StringValue(r.Password),
	}
}

func flattenSiteEntity(m *siteResourceModel, r *omada.SiteEntity) {
	m.Name = types.StringPointerValue(r.Name)
	m.Type = types.Int32PointerValue(r.Type)
	m.Region = types.StringPointerValue(r.Region)
	m.TimeZone = types.StringPointerValue(r.TimeZone)
	m.Scenario = types.StringPointerValue(r.Scenario)
	m.Longitude = types.Float64PointerValue(r.Longitude)
	m.Latitude = types.Float64PointerValue(r.Latitude)
	m.Address = types.StringPointerValue(r.Address)
	m.SupportES = types.BoolPointerValue(r.SupportES)
	m.SupportL2 = types.BoolPointerValue(r.SupportL2)

	flattenTagIds(&m.TagIDs, &r.TagIds)
}

func flattenSiteSummaryInfo(m *siteModel, r *omada.SiteSummaryInfo) {
	m.SiteId = types.StringPointerValue(r.SiteId)
	m.Name = types.StringPointerValue(r.Name)
	m.Region = types.StringPointerValue(r.Region)
	m.TimeZone = types.StringPointerValue(r.TimeZone)
	m.Scenario = types.StringPointerValue(r.Scenario)
	m.Longitude = types.Float64PointerValue(r.Longitude)
	m.Latitude = types.Float64PointerValue(r.Latitude)
	m.Address = types.StringPointerValue(r.Address)
	m.Type = types.Int32PointerValue(r.Type)
	m.SupportES = types.BoolPointerValue(r.SupportES)
	m.SupportL2 = types.BoolPointerValue(r.SupportL2)
	m.SitePublicIP = types.StringPointerValue(r.SitePublicIp)

	flattenTagIds(&m.TagIDs, &r.TagIds)
}
