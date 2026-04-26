package provider

import (
	"encoding/json"
	"net/http"
	"regexp"
	"terraform-provider-omada/internal/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var ProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"omada": providerserver.NewProtocol6WithError(New("test")()),
}

func Test_Provider_NoInlineConfig(t *testing.T) {
	// Ensure env-var fallbacks don't satisfy the provider.
	t.Setenv("OMADA_HOST", "")
	t.Setenv("OMADA_CONTROLLER_ID", "")
	t.Setenv("OMADA_CLIENT_ID", "")
	t.Setenv("OMADA_CLIENT_SECRET", "")

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					provider "omada" {}

					data "omada_sites" "test" {}
				`,
				ExpectError: regexp.MustCompile("Missing Omada API Host"),
			},
		},
	})
}

func Test_Provider_NoControllerId(t *testing.T) {
	// Ensure env-var fallbacks don't satisfy the provider.
	t.Setenv("OMADA_HOST", "")
	t.Setenv("OMADA_CONTROLLER_ID", "")
	t.Setenv("OMADA_CLIENT_ID", "")
	t.Setenv("OMADA_CLIENT_SECRET", "")

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					provider "omada" {
						host = "https://example.com"
					}

					data "omada_sites" "test" {}
				`,
				ExpectError: regexp.MustCompile("Missing Omada Controller ID"),
			},
		},
	})
}

func Test_Provider_NoClientId(t *testing.T) {
	// Ensure env-var fallbacks don't satisfy the provider.
	t.Setenv("OMADA_HOST", "")
	t.Setenv("OMADA_CONTROLLER_ID", "")
	t.Setenv("OMADA_CLIENT_ID", "")
	t.Setenv("OMADA_CLIENT_SECRET", "")

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					provider "omada" {
						host 		  = "https://example.com"
						controller_id = "test-controller-id"
					}

					data "omada_sites" "test" {}
				`,
				ExpectError: regexp.MustCompile("Missing Omada API Client ID"),
			},
		},
	})
}

func Test_Provider_NoClientSecret(t *testing.T) {
	// Ensure env-var fallbacks don't satisfy the provider.
	t.Setenv("OMADA_HOST", "")
	t.Setenv("OMADA_CONTROLLER_ID", "")
	t.Setenv("OMADA_CLIENT_ID", "")
	t.Setenv("OMADA_CLIENT_SECRET", "")

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					provider "omada" {
						host 		  = "https://example.com"
						controller_id = "test-controller-id"
						client_id	  = "test-client-id"
					}

					data "omada_sites" "test" {}
				`,
				ExpectError: regexp.MustCompile("Missing Omada API Client Secret"),
			},
		},
	})
}

func TestAcc_Provider_ClientCreationFailed(t *testing.T) {
	ts := acctest.NewTestServer(t)

	ts.TokenHandler = func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"errorCode": -1,
			"msg":       "No access token found",
		})
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ts.ProviderConfig + `data "omada_sites" "test" {}`,
				ExpectError: regexp.MustCompile("Unable to create Omada API Client"),
			},
		},
	})
}
