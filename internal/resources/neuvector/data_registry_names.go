package neuvector

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/theobori/go-neuvector/client"
	"github.com/theobori/go-neuvector/controller/scan"
)

var dataRegistryNamesSchema = map[string]*schema.Schema{
	"names": {
		Type:        schema.TypeSet,
		Description: "List of all register names.",
		Computed:    true,
		Elem:        &schema.Schema{Type: schema.TypeString},
	},
	"registry_type": {
		Type:        schema.TypeString,
		Description: "Type of the registry.",
		Optional:    true,
	},
}

func DataSourceRegistryNames() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRegistryNamesRead,
		Schema:      dataRegistryNamesSchema,
	}
}

func dataSourceRegistryNamesRead(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var names []string
	
	APIClient := meta.(*client.Client)

	registriesSummaries, err := scan.GetRegistries(APIClient)

	if err != nil {
		return diag.FromErr(err)
	}

	registryType := d.Get("registry_type").(string)

	// Add every registry name into the slice `names`
	for _, r := range registriesSummaries.Registries {
		if r.RegistryType == registryType || len(registryType) == 0 {
			names = append(names, r.Name)
		}
	}

	id, err := uuid.GenerateUUID()
	
	if err != nil {
		return diag.FromErr(err)
	}
	
	d.SetId(id)
	d.Set("names", names)

	return nil
}
