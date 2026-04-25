package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const testControllerID = "test-controller-id"

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"omada": providerserver.NewProtocol6WithError(New("test")()),
}

// providerConfig returns the provider HCL block configured to point at the
// supplied host (typically a httptest.Server URL).
func providerConfig(host string) string {
	return fmt.Sprintf(`
	provider "omada" {
		host          = %q
		controller_id = %q
		client_id     = "test-client-id"
		client_secret = "test-client-secret"
	}
	`, host, testControllerID)
}

// newTestServer spins up an httptest.Server backed by a fresh ServeMux. The
// authorize/token endpoint is pre-registered with a successful response so
// individual tests only need to wire up the endpoints they exercise. The
// returned string is the provider HCL block with `host` set to the server URL.
func newTestServer(t *testing.T) (*http.ServeMux, string) {
	t.Helper()

	mux := http.NewServeMux()

	mux.HandleFunc("POST /openapi/authorize/token", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"errorCode": 0,
			"msg":       "Open API Get Access Token successfully.",
			"result": map[string]any{
				"accessToken":  "AT-bjaJkIMIiekZY6NBufoQO4hdmJTswlwU",
				"tokenType":    "bearer",
				"expiresIn":    7200,
				"refreshToken": "RT-3ZjJgcORJSh76UCh7pj0rs5VRISIpagV",
			},
		})
	})

	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	return mux, providerConfig(server.URL)
}
