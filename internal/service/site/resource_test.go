package site_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"terraform-provider-omada/internal/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_SiteResource(t *testing.T) {
	mux, providerCfg := acctest.NewTestServer(t)

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

	readDeviceAccountResponse := `{
		"errorCode": 0,
		"msg": "",
		"result": {
			"username": "admin",
			"password": "password"
		}
	}`

	const emptyResponse = `{ "errorCode": 0, "msg": "" }`

	writeJSON := func(w http.ResponseWriter, body string) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}

	// Create
	mux.HandleFunc("POST /openapi/v1/{omadacId}/sites", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, `{
			"errorCode": 0,
			"msg": "",
			"result": { "siteId": "test-site-id" }
		}`)
	})

	// Read site
	mux.HandleFunc("GET /openapi/v1/{omadacId}/sites/{siteId}", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, readResponse)
	})

	// Read device account
	mux.HandleFunc("GET /openapi/v1/{omadacId}/sites/{siteId}/device-account", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, readDeviceAccountResponse)
	})

	// Update site
	mux.HandleFunc("PUT /openapi/v1/{omadacId}/sites/{siteId}", func(w http.ResponseWriter, r *http.Request) {
		newSite := make(map[string]any)
		if err := json.NewDecoder(r.Body).Decode(&newSite); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
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

		writeJSON(w, emptyResponse)
	})

	// Update device account
	mux.HandleFunc("PUT /openapi/v1/{omadacId}/sites/{siteId}/device-account", func(w http.ResponseWriter, r *http.Request) {
		deviceAccount := make(map[string]any)
		if err := json.NewDecoder(r.Body).Decode(&deviceAccount); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
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

		writeJSON(w, emptyResponse)
	})

	// Delete
	mux.HandleFunc("DELETE /openapi/v1/{omadacId}/sites/{siteId}", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, emptyResponse)
	})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerCfg + `
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
				Config: providerCfg + `
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
