// resource_promoteT.go
package neuvector

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/theobori/go-neuvector/client"
	"github.com/theobori/go-neuvector/controller/federation"
	"github.com/theobori/terraform-provider-neuvector/helper"
)

var resourcePromoteSchema = map[string]*schema.Schema{
	"server": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Server address",
	},
	"port": {
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Controller federation manager port, usually `11443`",
	},
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Cluster name",
	},
	"user": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Username",
	},
}

func resourcePromote() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePromoteCreate,
		ReadContext:   resourcePromoteRead,
		DeleteContext: resourcePromoteDelete,
		UpdateContext: resourcePromoteUpdate,

		Schema: resourcePromoteSchema,
	}
}

func resourcePromoteCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	APIClient := meta.(*client.Client)

	masterRestInfo := helper.FromSchemas[federation.MasterRestInfo](
		resourcePromoteSchema,
		d,
	)
	body := helper.FromSchemas[federation.FederationMetadata](
		resourcePromoteSchema,
		d,
	)

	body.MasterRestInfo = masterRestInfo

	if err := federation.Promote(APIClient, body); err != nil {
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
	APIClient := meta.(*client.Client)

	if err := federation.Demote(APIClient); err != nil {
		return diag.FromErr(err)
	}

	return nil
}