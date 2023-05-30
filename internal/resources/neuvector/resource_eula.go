// resource_EULAT.go
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

		Schema: resourceEULASchema,
	}
}

func resourceEULACreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	APIClient := meta.(*goneuvector.Client)

	eula := helper.FromSchemas[goneuvector.EULA](
		resourceEULASchema,
		d,
	)

	if err := APIClient.AcceptEULA(eula); err != nil {
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
	return nil
}

func resourceEULARead(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	APIClient := meta.(*goneuvector.Client)

	eula, err := APIClient.GetEULA()

	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("accepted", eula.EULA.Accepted)

	return nil
}

func resourceEULADelete(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	APIClient := meta.(*goneuvector.Client)

	eula := helper.FromSchemas[goneuvector.EULA](
		resourceEULASchema,
		d,
	)

	err := APIClient.AcceptEULA(
		goneuvector.EULA{
			Accepted: !eula.Accepted,
		},
	)

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
