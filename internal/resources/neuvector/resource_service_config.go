// resource_service_config.go
package neuvector

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	goneuvector "github.com/theobori/go-neuvector/neuvector"
	"github.com/theobori/terraform-provider-neuvector/internal/helper"
)

var resourceServiceConfigSchema = map[string]*schema.Schema{
	"services": {
		Type:        schema.TypeList,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Required:    true,
		Description: "The services to update.",
	},
	"policy_mode": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"baseline_profile": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"not_scored": {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  true,
	},
}

func ResourceServiceConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServiceConfigCreateOrUpdate,
		ReadContext:   resourceServiceConfigRead,
		DeleteContext: resourceServiceConfigDelete,
		UpdateContext: resourceServiceConfigCreateOrUpdate,

		Schema: resourceServiceConfigSchema,
	}
}

func resourceServiceConfigCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	APIClient := meta.(*goneuvector.Client)

	body := helper.FromSchemas[goneuvector.PatchServiceConfigBody](
		resourceServiceConfigSchema,
		d,
	)

	servicesRaw := d.Get("services").([]any)
	services, err := helper.FromSlice[string](servicesRaw)

	if err != nil {
		return diag.FromErr(err)
	}

	body.Services = services

	if err := APIClient.PatchServiceConfig(body); err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)

	return resourceServiceConfigRead(ctx, d, meta)
}

func resourceServiceConfigRead(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}

func resourceServiceConfigDelete(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	APIClient := meta.(*goneuvector.Client)

	servicesRaw := d.Get("services").([]any)
	services, err := helper.FromSlice[string](servicesRaw)

	if err != nil {
		return diag.FromErr(err)
	}

	body := goneuvector.PatchServiceConfigBody{
		Services:  services,
		NotScored: new(bool),
	}

	*body.NotScored = false

	if err := APIClient.PatchServiceConfig(body); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
