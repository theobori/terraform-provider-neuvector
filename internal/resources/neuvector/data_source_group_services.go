// data_source_group_services.go
package neuvector

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	goneuvector "github.com/theobori/go-neuvector/neuvector"
	"github.com/theobori/go-neuvector/util"
)

var dataGroupServicesSchema = map[string]*schema.Schema{
	"services": {
		Type:        schema.TypeSet,
		Description: "List of every service in the group.",
		Computed:    true,
		Elem:        &schema.Schema{Type: schema.TypeString},
	},
	"name": {
		Type:        schema.TypeString,
		Description: "The name of the gorup.",
		Required:    true,
	},
}

func DataSourceGroupServices() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGroupServicesRead,
		Schema:      dataGroupServicesSchema,
	}
}

func dataSourceGroupServicesRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var services []string

	APIClient := meta.(*goneuvector.Client)

	name := d.Get("name").(string)
	group, err := APIClient.
		WithContext(ctx).
		GetGroup(name)

	if err != nil {
		return diag.FromErr(err)
	}

	members := group.Group.Members

	for _, m := range members {
		exists, err := util.ItemExists(services, m.Service)

		if !exists && err == nil {
			services = append(services, m.Service)
		}
	}

	id, err := uuid.GenerateUUID()

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)
	d.Set("services", services)

	return nil
}
