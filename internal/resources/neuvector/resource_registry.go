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
		Description: "Name of the registry.",
	},
	"registry_type": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Type of the registry.",
	},
	"registry": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Registry URL.",
	},
	"filters": {
		Type:        schema.TypeList,
		Required:    true,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Description: "List of filters.",
	},
	"username": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Username for authenticate to the registry.",
		Sensitive:   true,
	},
	"password": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "password for authenticate to the registry.",
		Sensitive:   true,
	},
	"auth_token": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Authentication token.",
		Default:     "",
		Sensitive:   true,
	},
	"auth_with_token": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Flag indicating whether to authenticate using a token.",
		Default:     false,
	},
	"rescan_after_db_update": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Flag indicating whether to rescan after database update.",
		Default:     true,
	},
	"scan_layers": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Flag indicating whether to scan layers.",
		Default:     false,
	},
	"repo_limit": {
		Type:        schema.TypeInt,
		Optional:    true,
		Default:     200,
		Description: "Limit for the number of repositories.",
	},
	"tag_limit": {
		Type:        schema.TypeInt,
		Optional:    true,
		Default:     20,
		Description: "Max images tag to scan.",
	},
	"cfg_type": {
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "user_created",
		Description: "Configuration type",
	},
}

func ResourceRegistry() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRegistryCreate,
		ReadContext:   resourceRegistryRead,
		DeleteContext: resourceRegistryDelete,
		UpdateContext: resourceRegistryUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

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

func resourceRegistryCreate(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	APIClient := meta.(*client.Client)

	body, err := readRegistry(d)

	if err != nil {
		return diag.FromErr(err)
	}

	if err := scan.CreateRegistry(APIClient, *body); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(body.Name)

	return nil
}

func resourceRegistryUpdate(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	APIClient := meta.(*client.Client)

	if d.HasChanges("name", "registry_type") {
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
	var err error

	APIClient := meta.(*client.Client)

	r, err := scan.GetRegistry(
		APIClient,
		d.Id(),
	)

	if err != nil {
		return diag.FromErr(err)
	}

	// Forcing field overriding to permits terraform import
	if err = helper.TfFromStruct(r.Registry, d, true); err != nil {
		return diag.FromErr(err)
	}

	d.Set("filters", r.Registry.Filters)

	return nil
}

func resourceRegistryDelete(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	APIClient := meta.(*client.Client)

	if err := scan.DeleteRegistry(APIClient, d.Id()); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
