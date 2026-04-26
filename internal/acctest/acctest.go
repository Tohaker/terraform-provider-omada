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

// TestServer is a configurable httptest-backed Omada API stand-in.
//
// TokenHandler is invoked for POST /openapi/authorize/token requests. Tests may
// reassign it after construction to simulate alternate token-endpoint behavior
// (e.g. error responses) without re-registering on the underlying mux.
type TestServer struct {
	Mux            *http.ServeMux
	URL            string
	ProviderConfig string
	TokenHandler   http.HandlerFunc
}

func defaultTokenHandler(w http.ResponseWriter, _ *http.Request) {
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
}

func NewTestServer(t *testing.T) *TestServer {
	t.Helper()
	ts := &TestServer{
		Mux:          http.NewServeMux(),
		TokenHandler: defaultTokenHandler,
	}
	ts.Mux.HandleFunc("POST /openapi/authorize/token", func(w http.ResponseWriter, r *http.Request) {
		ts.TokenHandler(w, r)
	})
	server := httptest.NewServer(ts.Mux)
	t.Cleanup(server.Close)
	ts.URL = server.URL
	ts.ProviderConfig = providerConfig(server.URL)
	return ts
}
