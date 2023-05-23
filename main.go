package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/theobori/terraform-provider-neuvector/neuvector"
)

func main() {
	opts := &plugin.ServeOpts{
		ProviderFunc: neuvector.Provider,
	}

	plugin.Serve(opts)
}
