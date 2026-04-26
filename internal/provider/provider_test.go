package provider

import (
	"encoding/json"
	"fmt"
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

func Test_Provider_NoHost(t *testing.T) {
	// Ensure env-var fallbacks don't satisfy the provider.
	t.Setenv("OMADA_HOST", "")
	t.Setenv("OMADA_CONTROLLER_ID", "")
	t.Setenv("OMADA_CLIENT_ID", "")
	t.Setenv("OMADA_CLIENT_SECRET", "")
	t.Setenv("OMADA_TLS_SKIP_VERIFY", "")

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					provider "omada" {
						tls_skip_verify = true
					}

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
	t.Setenv("OMADA_TLS_SKIP_VERIFY", "")

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					provider "omada" {
						host            = "https://example.com"
						tls_skip_verify = true
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
	t.Setenv("OMADA_TLS_SKIP_VERIFY", "")

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					provider "omada" {
						host 		    = "https://example.com"
						controller_id   = "test-controller-id"
						tls_skip_verify = true
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
	t.Setenv("OMADA_TLS_SKIP_VERIFY", "")

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					provider "omada" {
						host 		    = "https://example.com"
						controller_id   = "test-controller-id"
						client_id	    = "test-client-id"
						tls_skip_verify = true
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

// Ensures that when tls_skip_verify is false the provider refuses a self-signed certificate.
func TestAcc_Provider_TLSSkipVerify_False_RejectsSelfSigned(t *testing.T) {
	t.Setenv("OMADA_TLS_SKIP_VERIFY", "")

	ts := acctest.NewTLSTestServer(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ts.ProviderConfigWithTLSSkipVerify(false) + `data "omada_sites" "test" {}`,
				ExpectError: regexp.MustCompile(`(?i)x509|certificate signed by unknown authority|tls`),
			},
		},
	})
}

// TestAcc_Provider_TLSSkipVerify_True_AcceptsSelfSigned ensures that when
// tls_skip_verify is true the provider successfully completes the TLS
// handshake against a self-signed certificate. We force a known controller-
// side error so the test asserts the request reached the server (i.e. the
// handshake succeeded) without needing to stub the entire API surface.
func TestAcc_Provider_TLSSkipVerify_True_AcceptsSelfSigned(t *testing.T) {
	t.Setenv("OMADA_TLS_SKIP_VERIFY", "")

	ts := acctest.NewTLSTestServer(t)
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

// TestAcc_Provider_TLSSkipVerify_EnvVar ensures the OMADA_TLS_SKIP_VERIFY env
// var is honored when the attribute is omitted from configuration.
func TestAcc_Provider_TLSSkipVerify_EnvVar(t *testing.T) {
	t.Setenv("OMADA_TLS_SKIP_VERIFY", "true")

	ts := acctest.NewTLSTestServer(t)
	ts.TokenHandler = func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"errorCode": -1,
			"msg":       "No access token found",
		})
	}

	// Build a provider config that omits tls_skip_verify so the env var is used.
	cfg := fmt.Sprintf(`
provider "omada" {
  host          = %q
  controller_id = %q
  client_id     = "test-client-id"
  client_secret = "test-client-secret"
}
`, ts.URL, acctest.TestControllerID)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      cfg + `data "omada_sites" "test" {}`,
				ExpectError: regexp.MustCompile("Unable to create Omada API Client"),
			},
		},
	})
}
