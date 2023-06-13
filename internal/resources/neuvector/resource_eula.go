// resource_eula.go
package neuvector

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	goneuvector "github.com/theobori/go-neuvector/neuvector"
	"github.com/theobori/terraform-provider-neuvector/internal/helper"
)

var resourceEULASchema = map[string]*schema.Schema{
	"accepted": {
		Type:        schema.TypeBool,
		Required:    true,
		Description: "Accept the EULA.",
	},
}

func ResourceEULA() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEULACreate,
		ReadContext:   resourceEULARead,
		DeleteContext: resourceEULADelete,
		UpdateContext: resourceEULAUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: resourceEULASchema,
	}
}

func resourceEULACreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	APIClient := meta.(*goneuvector.Client)

	eula := helper.FromSchemas[goneuvector.EULA](
		resourceEULASchema,
		d,
	)

	if err := APIClient.WithContext(ctx).AcceptEULA(eula); err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)

	return resourceEULARead(ctx, d, meta)
}

func resourceEULAUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	APIClient := meta.(*goneuvector.Client)

	eula := helper.FromSchemas[goneuvector.EULA](
		resourceEULASchema,
		d,
	)

	if err := APIClient.WithContext(ctx).AcceptEULA(eula); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceEULARead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	APIClient := meta.(*goneuvector.Client)

	eula, err := APIClient.
		WithContext(ctx).	
		GetEULA()

	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("accepted", eula.EULA.Accepted)

	return nil
}

func resourceEULADelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	APIClient := meta.(*goneuvector.Client)

	eula := helper.FromSchemas[goneuvector.EULA](
		resourceEULASchema,
		d,
	)

	err := APIClient.
		WithContext(ctx).
		AcceptEULA(
			goneuvector.EULA{
				Accepted: !eula.Accepted,
			},
	)

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
