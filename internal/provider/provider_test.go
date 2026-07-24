package provider_test

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/TheWhale01/terraform-provider-jellyfin/internal/provider"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"jellyfin": providerserver.NewProtocol6WithError(provider.New("test")()),
}
