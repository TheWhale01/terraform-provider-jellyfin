package main

import (
	"context"
	"flag"
	"log"

	"github.com/TheWhale01/terraform-provider-jellyfin/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

var (
	version string = "dev"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with debugger support")
	flag.Parse()

	opts := providerserver.ServeOpts {
		Address: "registry.terraform.io/TheWhale01/jellyfin",
		Debug: debug,
	}

	err := providerserver.Serve(context.Background(), provider.New(version), opts)
	if err != nil {
		log.Fatal(err.Error())
	}
}
