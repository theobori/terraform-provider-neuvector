// resource_registry.go
package neuvector

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/theobori/go-neuvector/client"
	"github.com/theobori/go-neuvector/controller/scan"
	"github.com/theobori/terraform-provider-neuvector/internal/helper"
)

var resourceRegistrySchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Registry name",
	},
	"registry_type": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Registry type",
	},
	"registry": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Registry address",
	},
	"filters": {
		Type:        schema.TypeList,
		Required:    true,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Description: "Filters string list",
	},
	"username": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Registry username",
		Sensitive:   true,
	},
	"password": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Registry password",
		Sensitive:   true,
	},
	"auth_token": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Authentication token",
		Sensitive:   true,
	},
	"auth_with_token": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "That said if you are going to authenticate to the registry with a token",
		Default:     false,
	},
	"rescan_after_db_update": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Rescan after the CVE database update",
		Default:     true,
	},
	"scan_layers": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Scan the layers",
		Default:     false,
	},
	"repo_limit": {
		Type:        schema.TypeInt,
		Optional:    true,
		Description: "Repositories max amount on the registry",
	},
	"tag_limit": {
		Type:        schema.TypeInt,
		Optional:    true,
		Description: "Max images tag to scan",
	},
	"cfg_type": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Configuration type",
	},
}

func ResourceRegistry() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRegistryCreate,
		ReadContext:   resourceRegistryRead,
		DeleteContext: resourceRegistryDelete,
		UpdateContext: resourceRegistryUpdate,

		Schema: resourceRegistrySchema,
	}
}

func readRegistry(d *schema.ResourceData) (*scan.CreateRegistryBody, error) {
	var ret scan.CreateRegistryBody

	filtersRaw := d.Get("filters").([]any)
	filters, err := helper.FromSlice[string](filtersRaw)

	if err != nil {
		return &ret, err
	}

	ret = helper.FromSchemas[scan.CreateRegistryBody](
		resourceRegistrySchema,
		d,
	)

	ret.Filters = filters

	return &ret, nil
}

func resourceRegistryCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	APIClient := meta.(*client.Client)

	body, err := readRegistry(d)

	if err != nil {
		return diag.FromErr(err)
	}

	scan.CreateRegistry(
		APIClient,
		*body,
	)

	d.SetId(body.Name)

	return resourceRegistryRead(ctx, d, meta)
}

func resourceRegistryUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	APIClient := meta.(*client.Client)

	if d.HasChanges(
		"name",
		"registry_type",
	) {
		return diag.Errorf("You are not allowed to change the registry name and type.")
	}

	body, err := readRegistry(d)

	if err != nil {
		return diag.FromErr(err)
	}

	if err := scan.PatchRegistry(APIClient, *body, body.Name); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceRegistryRead(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}

func resourceRegistryDelete(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	APIClient := meta.(*client.Client)

	if err := scan.DeleteRegistry(APIClient, d.Id()); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
