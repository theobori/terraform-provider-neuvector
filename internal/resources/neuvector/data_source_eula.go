// data_source_eula.go
package neuvector

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	goneuvector "github.com/theobori/go-neuvector/neuvector"
)

var dataEULASchema = map[string]*schema.Schema{
	"accepted": {
		Type:        schema.TypeBool,
		Description: "Represents the EULA status.",
		Computed:    true,
	},
}

func DataSourceEULA() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEULARead,
		Schema:      dataEULASchema,
	}
}

func dataSourceEULARead(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	APIClient := meta.(*goneuvector.Client)

	eula, err := APIClient.GetEULA()

	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)
	d.Set("accepted", eula.EULA.Accepted)

	return nil
}
