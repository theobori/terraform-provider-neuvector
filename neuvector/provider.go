package neuvector

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/theobori/go-neuvector/client"
)

// Provider returns a *schema.Provider.
func Provider() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Represents the NeuVector username.",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Represents the NeuVector password.",
			},
			"base_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Represents the NeuVector Controller REST API base url.",
			},
			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Skip the TLS verification. Default: `true`.",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{},
		ResourcesMap: map[string]*schema.Resource{
			"neuvector_admission_rule": resourceAdmissionRule(),
			"neuvector_promote":        resourcePromote(),
			"neuvector_registry":       resourceRegistry(),
			"neuvector_policy":         resourcePolicy(),
		},
	}

	provider.ConfigureContextFunc = func(_ context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
		// Setup the authentication
		auth := client.NewClientAuth(
			d.Get("username").(string),
			d.Get("password").(string),
		)

		// Configure the API client
		config := client.NewClientConfig(auth)

		config.Insecure = d.Get("insecure").(bool)
		config.BaseUrl = d.Get("base_url").(string)

		// Get a new client
		APIClient, err := client.NewClient(config)

		if err != nil {
			return nil, diag.FromErr(err)
		}

		return APIClient, nil
	}

	return provider
}
