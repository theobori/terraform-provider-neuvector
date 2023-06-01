// resource_service.go
package neuvector

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	goneuvector "github.com/theobori/go-neuvector/neuvector"
	"github.com/theobori/terraform-provider-neuvector/internal/helper"
)

var resourceServiceSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Description: "The name of the service.",
		Required:    true,
	},
	"comment": {
		Type:        schema.TypeString,
		Description: "The comment of the service.",
		Optional:    true,
	},
	"domain": {
		Type:        schema.TypeString,
		Description: "Represents the namespace.",
		Optional:    true,
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

func ResourceService() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServiceCreate,
		ReadContext:   resourceServiceRead,
		DeleteContext: resourceServiceDelete,
		UpdateContext: resourceServiceUpdate,

		Schema: resourceServiceSchema,
	}
}

func resourceServiceCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	APIClient := meta.(*goneuvector.Client)

	body := helper.FromSchemas[goneuvector.CreateServiceBody](
		resourceServiceSchema,
		d,
	)

	if err := APIClient.CreateService(body); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(body.Name)

	return resourceServiceRead(ctx, d, meta)
}

func resourceServiceUpdate(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	if !d.HasChanges(
		"policy_mode",
		"baseline_profile",
		"not_scored",
	) {
		return nil
	}

	APIClient := meta.(*goneuvector.Client)

	body := helper.FromSchemas[goneuvector.PatchServiceConfigBody](
		resourceServiceSchema,
		d,
	)

	body.Services = []string{d.Id()}

	if err := APIClient.PatchServiceConfig(body); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceServiceRead(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}

func resourceServiceDelete(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	APIClient := meta.(*goneuvector.Client)

	domain := d.Get("domain").(string)
	name := "nv." + d.Id()

	if domain != "" {
		name += "." + domain
	}

	if err := APIClient.DeleteGroup(name); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
