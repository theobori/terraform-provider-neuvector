package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	goneuvector "github.com/theobori/go-neuvector/neuvector"
	"github.com/theobori/terraform-provider-neuvector/internal/resources/neuvector"
)

// Provider returns a *schema.Provider.
func Provider() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NEUVECTOR_PUSERNAME", goneuvector.DefaultUsername),
				Description: "Represents the NeuVector username.",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NEUVECTOR_PASSWORD", goneuvector.DefaultPassword),
				Description: "Represents the NeuVector password.",
			},
			"base_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NEUVECTOR_BASE_URL", goneuvector.DefaultBaseUrl),
				Description: "Represents the NeuVector Controller REST API base url.",
			},
			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Skip the TLS verification. Default: `true`.",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			// neuvector
			"neuvector_admission_rule": neuvector.ResourceAdmissionRule(),
			"neuvector_promote":        neuvector.ResourcePromote(),
			"neuvector_registry":       neuvector.ResourceRegistry(),
			"neuvector_policy":         neuvector.ResourcePolicy(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			// neuvector
			"neuvector_registry":       neuvector.DataSourceRegistry(),
			"neuvector_registry_names": neuvector.DataSourceRegistryNames(),
			"neuvector_policy_ids":     neuvector.DataSourcePolicyIDs(),
		},
	}

	provider.ConfigureContextFunc = func(_ context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
		// Setup the authentication
		auth := goneuvector.NewClientAuth(
			d.Get("username").(string),
			d.Get("password").(string),
		)

		// Configure the API client
		config := goneuvector.NewClientConfig(auth)

		config.Insecure = d.Get("insecure").(bool)
		config.BaseUrl = d.Get("base_url").(string)

		// Get a new client
		APIClient, err := goneuvector.NewClient(config)

		if err != nil {
			return nil, diag.FromErr(err)
		}

		return APIClient, nil
	}

	return provider
}
