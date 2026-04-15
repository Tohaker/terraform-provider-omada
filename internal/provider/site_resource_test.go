package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func TestAccSiteResource(t *testing.T) {
	activateHTTPMock(t)

	createResponse := `{
		"errorCode": 0,
		"msg": "",
		"result": {
			"siteId": "test-site-id"
		}
	}`

	httpmock.RegisterResponder(
		"POST",
		`=~^.*/openapi/v1/.*/sites$`,
		httpmock.NewStringResponder(
			http.StatusOK,
			createResponse,
		).HeaderSet(http.Header{
			"Content-Type": []string{"application/json"},
		}),
	)

	readResponse := `{
		"errorCode": 0,
		"msg": "",
		"result": {
			"siteId": "test-site-id",
			"name": "Test Site",
			"type": 0,
			"tagIds": [],
			"region": "United Kingdom",
			"timeZone": "Europe/London",
			"ntpEnable": true,
			"ntpServers": [],
			"scenario": "Home",
			"longitude": null,
			"latitude": null,
			"address": null,
			"supportES": true,
			"supportL2": true
		}
	}`

	httpmock.RegisterResponder(
		"GET",
		`=~^.*/openapi/v1/.*/sites/test-site-id$`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusOK, readResponse)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		},
	)

	readDeviceAccountResponse := `{
		"errorCode": 0,
		"msg": "",
		"result": {
			"username": "admin",
			"password": "password"
		}
	}`

	httpmock.RegisterResponder(
		"GET",
		`=~^.*/openapi/v1/.*/sites/test-site-id/device-account$`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusOK, readDeviceAccountResponse)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		},
	)

	emptyStringResponse := `{ "errorCode": 0, "msg": "" }`

	emptyStringResponder := httpmock.NewStringResponder(
		http.StatusOK,
		emptyStringResponse,
	).HeaderSet(http.Header{
		"Content-Type": []string{"application/json"},
	})

	httpmock.RegisterResponder(
		"PUT",
		"=~^.*/openapi/v1/.*/sites/test-site-id$",
		func(req *http.Request) (*http.Response, error) {
			newSite := make(map[string]any)
			if err := json.NewDecoder(req.Body).Decode(&newSite); err != nil {
				return httpmock.NewStringResponse(400, ""), nil
			}

			readResponse = fmt.Sprintf(`{
				"errorCode": 0,
				"msg": "",
				"result": {
					"siteId": "test-site-id",
					"name": "%s",
					"type": 0,
					"tagIds": [],
					"region": "%s",
					"timeZone": "%s",
					"ntpEnable": true,
					"ntpServers": [],
					"scenario": "%s",
					"longitude": null,
					"latitude": null,
					"address": null,
					"supportES": true,
					"supportL2": true
				}
			}`,
				newSite["name"],
				newSite["region"],
				newSite["timeZone"],
				newSite["scenario"])

			resp := httpmock.NewStringResponse(http.StatusOK, emptyStringResponse)
			resp.Header.Add("Content-Type", "application/json")

			return resp, nil
		},
	)

	httpmock.RegisterResponder(
		"PUT",
		`=~^.*/openapi/v1/.*/sites/test-site-id/device-account$`,
		func(r *http.Request) (*http.Response, error) {
			deviceAccount := make(map[string]any)
			if err := json.NewDecoder(r.Body).Decode(&deviceAccount); err != nil {
				return httpmock.NewStringResponse(400, ""), nil
			}

			readDeviceAccountResponse = fmt.Sprintf(`{
				"errorCode": 0,
				"msg": "",
				"result": {
					"username": "%s",
					"password": "%s"
				}
			}`,
				deviceAccount["username"],
				deviceAccount["password"])

			resp := httpmock.NewStringResponse(http.StatusOK, emptyStringResponse)
			resp.Header.Add("Content-Type", "application/json")

			return resp, nil
		},
	)

	httpmock.RegisterResponder(
		"DELETE",
		"=~^.*/openapi/v1/.*/sites/.*$",
		emptyStringResponder,
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
				resource "omada_site" "test" {
					name      = "Test Site"
					region    = "United Kingdom"
					time_zone = "Europe/London"
					scenario  = "Home"

					device_account_setting = {
						username = "admin"
						password = "password"
					}
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify user specified values are set
					resource.TestCheckResourceAttr("omada_site.test", "name", "Test Site"),
					resource.TestCheckResourceAttr("omada_site.test", "region", "United Kingdom"),
					resource.TestCheckResourceAttr("omada_site.test", "time_zone", "Europe/London"),
					resource.TestCheckResourceAttr("omada_site.test", "scenario", "Home"),
					resource.TestCheckResourceAttr("omada_site.test", "device_account_setting.username", "admin"),
					resource.TestCheckResourceAttr("omada_site.test", "device_account_setting.password", "password"),

					// Verify computed values are set
					resource.TestCheckResourceAttr("omada_site.test", "site_id", "test-site-id"),
					resource.TestCheckResourceAttr("omada_site.test", "type", "0"),
					resource.TestCheckResourceAttr("omada_site.test", "support_es", "true"),
					resource.TestCheckResourceAttr("omada_site.test", "support_l2", "true"),
				),
			},
			// Import state testing
			{
				ResourceName:                         "omada_site.test",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "test-site-id",
				ImportStateVerifyIdentifierAttribute: "site_id",
			},
			// Update and Read testing
			{
				Config: providerConfig + `
				resource "omada_site" "test" {
					name      = "Updated Site"
					region    = "United States"
					time_zone = "UTC"
					scenario  = "Home"

					device_account_setting = {
						username = "admin"
						password = "password2"
					}
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify user specified values updated
					resource.TestCheckResourceAttr("omada_site.test", "name", "Updated Site"),
					resource.TestCheckResourceAttr("omada_site.test", "region", "United States"),
					resource.TestCheckResourceAttr("omada_site.test", "time_zone", "UTC"),
					resource.TestCheckResourceAttr("omada_site.test", "device_account_setting.password", "password2"),

					// Verify computed values are updated
					resource.TestCheckResourceAttr("omada_site.test", "site_id", "test-site-id"),
					resource.TestCheckResourceAttr("omada_site.test", "type", "0"),
					resource.TestCheckResourceAttr("omada_site.test", "support_es", "true"),
					resource.TestCheckResourceAttr("omada_site.test", "support_l2", "true"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
