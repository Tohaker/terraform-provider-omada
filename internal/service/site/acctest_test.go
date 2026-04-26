package site_test

import (
	"terraform-provider-omada/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var ProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"omada": providerserver.NewProtocol6WithError(provider.New("test")()),
}
