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

func providerConfigWithTLSSkipVerify(host string, tlsSkipVerify bool) string {
	return fmt.Sprintf(`
provider "omada" {
  host            = %q
  controller_id   = %q
  client_id       = "test-client-id"
  client_secret   = "test-client-secret"
  tls_skip_verify = %t
}
`, host, TestControllerID, tlsSkipVerify)
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

// NewTLSTestServer is like NewTestServer but serves HTTPS using httptest's
// auto-generated self-signed certificate. Use it to exercise the provider's
// tls_skip_verify behavior. The default ProviderConfig sets
// tls_skip_verify = true so the test server is reachable; use
// ProviderConfigWithTLSSkipVerify(false) to assert verification failures.
func NewTLSTestServer(t *testing.T) *TestServer {
	t.Helper()
	ts := &TestServer{
		Mux:          http.NewServeMux(),
		TokenHandler: defaultTokenHandler,
	}
	ts.Mux.HandleFunc("POST /openapi/authorize/token", func(w http.ResponseWriter, r *http.Request) {
		ts.TokenHandler(w, r)
	})
	server := httptest.NewTLSServer(ts.Mux)
	t.Cleanup(server.Close)
	ts.URL = server.URL
	ts.ProviderConfig = providerConfigWithTLSSkipVerify(server.URL, true)
	return ts
}

// ProviderConfigWithTLSSkipVerify returns a provider configuration block for
// this test server with tls_skip_verify set to the given value.
func (ts *TestServer) ProviderConfigWithTLSSkipVerify(tlsSkipVerify bool) string {
	return providerConfigWithTLSSkipVerify(ts.URL, tlsSkipVerify)
}
