package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/theobori/terraform-provider-neuvector/internal/provider"
)

func main() {
	opts := &plugin.ServeOpts{
		ProviderFunc: provider.Provider,
	}

	plugin.Serve(opts)
}
