package provider

import (
	"crypto/tls"
	"net/http"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/jarcoal/httpmock"
)

const (
	// providerConfig is a shared configuration to combine with the actual
	// test configuration so the Omada client is properly configured.
	// It is also possible to use the OMADA_ environment variables instead,
	// such as updating the Makefile and running the testing through that tool.
	providerConfig = `
	provider "omada" {
		host     		= "http://test.omada/api"
		customer_id 	= "test-customer-id"
		client_id 		= "test-client-id"
		client_secret 	= "test-client-secret"
	}
	`
)

var (
	testHTTPClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	// testAccProtoV6ProviderFactories are used to instantiate a provider during
	// acceptance testing. The factory function will be invoked for every Terraform
	// CLI command executed to create a provider server to which the CLI can
	// reattach.
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"omada": providerserver.NewProtocol6WithError(New("test")()),
	}
)

func TestMain(m *testing.M) {
	newOmadaHTTPClient = func() *http.Client {
		return testHTTPClient
	}

	httpmock.ActivateNonDefault(testHTTPClient)
	code := m.Run()
	httpmock.DeactivateNonDefault(testHTTPClient)
	os.Exit(code)
}

func activateHTTPMock(t *testing.T) {
	t.Helper()
	httpmock.Reset()
	registerAuthResponder()
	t.Cleanup(httpmock.Reset)
}

func registerAuthResponder() {
	httpmock.RegisterResponder(
		"POST",
		`=~^.*/openapi/authorize/token\?grant_type=client_credentials`,
		httpmock.NewStringResponder(
			http.StatusOK,
			`{"errorCode":0,"msg":"Open API Get Access Token successfully.","result":{"accessToken":"AT-bjaJkIMIiekZY6NBufoQO4hdmJTswlwU","tokenType":"bearer","expiresIn":7200,"refreshToken":"RT-3ZjJgcORJSh76UCh7pj0rs5VRISIpagV"}}`).HeaderSet(http.Header{
			"Content-Type": []string{"application/json"},
		}),
	)
}
