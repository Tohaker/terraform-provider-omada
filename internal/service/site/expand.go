package site

import (
	"github.com/Tohaker/omada-go-sdk/omada"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func expandTagIds(TagIDs *[]types.String) []string {
	var tagIds []string
	for _, tagId := range *TagIDs {
		tagIds = append(tagIds, tagId.ValueString())
	}

	return tagIds
}

func expandDeviceAccountSetting(deviceAccountSetting *siteDeviceAccountSettingModel) omada.DeviceAccountSettingOpenApiVO {
	return omada.DeviceAccountSettingOpenApiVO{
		Username: deviceAccountSetting.Username.ValueString(),
		Password: deviceAccountSetting.Password.ValueString(),
	}
}

func expandCreateSiteEntity(plan siteResourceModel) omada.CreateSiteEntity {
	var siteEntity omada.CreateSiteEntity
	siteEntity.Name = plan.Name.ValueString()
	siteEntity.Type = plan.Type.ValueInt32Pointer()
	siteEntity.Region = plan.Region.ValueString()
	siteEntity.TimeZone = plan.TimeZone.ValueString()
	siteEntity.Scenario = plan.Scenario.ValueString()
	siteEntity.Longitude = plan.Longitude.ValueFloat64Pointer()
	siteEntity.Latitude = plan.Latitude.ValueFloat64Pointer()
	siteEntity.Address = plan.Address.ValueStringPointer()
	siteEntity.SupportES = plan.SupportES.ValueBoolPointer()
	siteEntity.SupportL2 = plan.SupportL2.ValueBoolPointer()

	return siteEntity
}

func expandUpdateSiteEntity(plan siteResourceModel) omada.UpdateSiteEntity {
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

	return siteEntity
}
