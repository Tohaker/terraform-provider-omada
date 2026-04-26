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

func flattenSites(m *sitesDataSourceModel, r *[]omada.SiteSummaryInfo) {
	for _, site := range *r {
		var siteState siteModel

		siteState.SiteId = types.StringPointerValue(site.SiteId)
		siteState.Name = types.StringPointerValue(site.Name)
		siteState.Region = types.StringPointerValue(site.Region)
		siteState.TimeZone = types.StringPointerValue(site.TimeZone)
		siteState.Scenario = types.StringPointerValue(site.Scenario)
		siteState.Longitude = types.Float64PointerValue(site.Longitude)
		siteState.Latitude = types.Float64PointerValue(site.Latitude)
		siteState.Address = types.StringPointerValue(site.Address)
		siteState.Type = types.Int32PointerValue(site.Type)
		siteState.SupportES = types.BoolPointerValue(site.SupportES)
		siteState.SupportL2 = types.BoolPointerValue(site.SupportL2)
		siteState.SitePublicIP = types.StringPointerValue(site.SitePublicIp)

		flattenTagIds(&siteState.TagIDs, &site.TagIds)

		m.Sites = append(m.Sites, siteState)
	}

}
