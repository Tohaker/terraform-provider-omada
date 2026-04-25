package acctest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

const TestControllerID = "test-controller-id"

func providerConfig(host string) string {
	return fmt.Sprintf(`
provider "omada" {
  host          = %q
  controller_id = %q
  client_id     = "test-client-id"
  client_secret = "test-client-secret"
}
`, host, TestControllerID)
}

func NewTestServer(t *testing.T) (*http.ServeMux, string) {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("POST /openapi/authorize/token", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"errorCode": 0,
			"msg":       "Open API Get Access Token successfully.",
			"result": map[string]any{
				"accessToken":  "AT-test",
				"tokenType":    "bearer",
				"expiresIn":    7200,
				"refreshToken": "RT-test",
			},
		})
	})
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)
	return mux, providerConfig(server.URL)
}
