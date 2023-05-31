// resource_service.go
package neuvector

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	// goneuvector "github.com/theobori/go-neuvector/neuvector"
)

var resourceServiceSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Name of the Service.",
	},
}

func ResourceService() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServiceCreate,
		ReadContext:   resourceServiceRead,
		DeleteContext: resourceServiceDelete,
		UpdateContext: resourceServiceUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: resourceServiceSchema,
	}
}

func resourceServiceCreate(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}

func resourceServiceUpdate(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}

func resourceServiceRead(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {

	return nil
}

func resourceServiceDelete(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}
