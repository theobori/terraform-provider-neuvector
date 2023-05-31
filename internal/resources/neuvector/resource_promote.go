// resource_promote.go
package neuvector

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	goneuvector "github.com/theobori/go-neuvector/neuvector"
	"github.com/theobori/terraform-provider-neuvector/internal/helper"
)

var resourcePromoteSchema = map[string]*schema.Schema{
	"server": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Server address.",
	},
	"port": {
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Controller federation manager port, usually `11443`.",
	},
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Cluster name.",
	},
	"user": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Username.",
	},
}

func ResourcePromote() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePromoteCreate,
		ReadContext:   resourcePromoteRead,
		DeleteContext: resourcePromoteDelete,
		UpdateContext: resourcePromoteUpdate,

		Schema: resourcePromoteSchema,
	}
}

func resourcePromoteCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	APIClient := meta.(*goneuvector.Client)

	masterRestInfo := helper.FromSchemas[goneuvector.MasterRestInfo](
		resourcePromoteSchema,
		d,
	)
	body := helper.FromSchemas[goneuvector.FederationMetadata](
		resourcePromoteSchema,
		d,
	)

	body.MasterRestInfo = masterRestInfo

	if err := APIClient.Promote(body); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(body.Name)

	return resourcePromoteRead(ctx, d, meta)
}

func resourcePromoteUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}

func resourcePromoteRead(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}

func resourcePromoteDelete(_ context.Context, _ *schema.ResourceData, meta any) diag.Diagnostics {
	APIClient := meta.(*goneuvector.Client)

	if err := APIClient.Demote(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
