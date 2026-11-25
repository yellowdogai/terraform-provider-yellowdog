// Terraform Provider for the YellowDog Platform.
package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/yellowdogai/terraform-provider-yellowdog/internal/provider"
)

// version is set by goreleaser at build time via -ldflags.
var version = "dev"

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "run the provider with support for debuggers like delve")
	flag.Parse()

	err := providerserver.Serve(context.Background(), provider.New(version), providerserver.ServeOpts{
		Address: "registry.terraform.io/yellowdogai/yellowdog",
		Debug:   debug,
	})
	if err != nil {
		log.Fatal(err.Error())
	}
}
