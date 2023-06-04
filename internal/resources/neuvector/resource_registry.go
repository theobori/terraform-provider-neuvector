// resource_registry.go
package neuvector

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	goneuvector "github.com/theobori/go-neuvector/neuvector"
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
	"scan_after_add": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Indicates if the registry must be scanned immediatly after beeing added.",
		Default:     false,
	}}

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

func readRegistry(d *schema.ResourceData) (*goneuvector.CreateRegistryBody, error) {
	var ret goneuvector.CreateRegistryBody

	filtersRaw := d.Get("filters").([]any)
	filters, err := helper.FromSlice[string](filtersRaw)

	if err != nil {
		return &ret, err
	}

	ret = helper.FromSchemas[goneuvector.CreateRegistryBody](
		resourceRegistrySchema,
		d,
	)

	ret.Filters = filters

	return &ret, nil
}

func resourceRegistryCreate(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	APIClient := meta.(*goneuvector.Client)

	body, err := readRegistry(d)

	if err != nil {
		return diag.FromErr(err)
	}

	if err := APIClient.CreateRegistry(*body); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(body.Name)

	scan := d.Get("scan_after_add").(bool)

	if scan {
		// Needed pause because NeuVector is sending the HTTP status code
		// before the registry is fully added.
		time.Sleep(3 * time.Second)

		err = APIClient.Post(
			fmt.Sprintf("/scan/registry/%s/scan", body.Name),
			nil,
			nil,
		)

		if err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func resourceRegistryUpdate(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	APIClient := meta.(*goneuvector.Client)

	if d.HasChanges("name", "registry_type") {
		return diag.Errorf("You are not allowed to change the registry name and type.")
	}

	body, err := readRegistry(d)

	if err != nil {
		return diag.FromErr(err)
	}

	if err := APIClient.PatchRegistry(*body, body.Name); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceRegistryRead(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var err error

	APIClient := meta.(*goneuvector.Client)

	r, err := APIClient.GetRegistry(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	password := d.Get("password").(string)

	// Forcing field overriding to permits terraform import
	if err = helper.TfFromStruct(r.Registry, d, true); err != nil {
		return diag.FromErr(err)
	}

	d.Set("filters", r.Registry.Filters)
	d.Set("password", password)

	return nil
}

func resourceRegistryDelete(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	APIClient := meta.(*goneuvector.Client)

	if err := APIClient.DeleteRegistry(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
