package testutils

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/theobori/terraform-provider-neuvector/internal/provider"
)

var (
	// A factory containing provider names associated
	// with functions that return the provider
	ProviderFactories map[string]func() (*schema.Provider, error)

	// Main provider instance, used in acceptance tests
	Provider *schema.Provider
)

func init() {
	if Provider != nil {
		return
	}

	Provider = provider.Provider()

	ProviderFactories = map[string]func() (*schema.Provider, error){
		"neuvector": func() (*schema.Provider, error) {
			return provider.Provider(), nil
		},
	}

	if testAccEnabled("TF_ACC") {
		testAccProviderConfigure()
	}
}

// Out of the Terraform scope configuration
func testAccProviderConfigure() {
	err := Provider.Configure(
		context.Background(),
		terraform.NewResourceConfigRaw(nil),
	)

	if err != nil {
		panic(fmt.Sprintf("failed to configure provider: %v", err))
	}
}

// Check for environment variable flag value
// `name` --> env var name
func testAccEnabled(name string) bool {
	v, ok := os.LookupEnv(name)

	if !ok {
		return false
	}

	enabled, err := strconv.ParseBool(v)

	if err != nil {
		panic(fmt.Sprintf("%s must be set to a boolean value", name))
	}

	return enabled
}
