// data_source_registry.go
package neuvector

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	goneuvector "github.com/theobori/go-neuvector/neuvector"
	"github.com/theobori/terraform-provider-neuvector/internal/helper"
)

var dataRegistrySchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Description: "Name of the registry.",
		Required:    true,
	},
	"registry_type": {
		Type:        schema.TypeString,
		Description: "Type of the registry.",
		Computed:    true,
	},
	"registry": {
		Type:        schema.TypeString,
		Description: "Registry URL.",
		Computed:    true,
	},
	"username": {
		Type:        schema.TypeString,
		Description: "Username for authentication.",
		Computed:    true,
		Sensitive:   true,
	},
	"password": {
		Type:        schema.TypeString,
		Description: "Password for authentication.",
		Computed:    true,
		Sensitive:   true,
	},
	"auth_token": {
		Type:        schema.TypeString,
		Description: "Authentication token.",
		Computed:    true,
		Sensitive:   true,
	},
	"auth_with_token": {
		Type:        schema.TypeBool,
		Description: "Flag indicating whether to authenticate using a token.",
		Computed:    true,
	},
	"filters": {
		Type:        schema.TypeList,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Description: "List of filters.",
		Computed:    true,
	},
	"rescan_after_db_update": {
		Type:        schema.TypeBool,
		Description: "Flag indicating whether to rescan after database update.",
		Computed:    true,
	},
	"scan_layers": {
		Type:        schema.TypeBool,
		Description: "Flag indicating whether to scan layers.",
		Computed:    true,
	},
	"repo_limit": {
		Type:        schema.TypeInt,
		Description: "Limit for the number of repositories.",
		Computed:    true,
	},
	"tag_limit": {
		Type:        schema.TypeInt,
		Description: "Limit for the number of tags.",
		Computed:    true,
	},
	// "schedule": {
	// 	Type:        schema.TypeList,
	// 	Elem:        &schema.Schema{Type: schema.TypeString},
	// 	Description: "Schedule configuration.",
	// 	Optional:    true,
	// },
	// "aws_key": {
	// 	Type:        schema.TypeList,
	// 	Elem:        &schema.Schema{Type: schema.TypeString},
	// 	Description: "AWS key configuration.",
	// 	Optional:    true,
	// },
	// "jfrog_xray": {
	// 	Type:        schema.TypeList,
	// 	Elem:        &schema.Schema{Type: schema.TypeString},
	// 	Description: "JFrog Xray configuration.",
	// 	Optional:    true,
	// },
	// "gcr_key": {
	// 	Type:        schema.TypeList,
	// 	Elem:        &schema.Schema{Type: schema.TypeString},
	// 	Description: "GCR key configuration.",
	// 	Optional:    true,
	// },
	// "jfrog_mode": {
	// 	Type:        schema.TypeString,
	// 	Description: "JFrog mode.",
	// 	Computed:    true,
	// },
	// "gitlab_external_url": {
	// 	Type:        schema.TypeString,
	// 	Description: "GitLab external URL.",
	// 	Computed:    true,
	// },
	// "gitlab_private_token": {
	// 	Type:        schema.TypeString,
	// 	Description: "GitLab private token.",
	// 	Computed:    true,
	// 	Sensitive: true,
	// },
	// "ibm_cloud_token_url": {
	// 	Type:        schema.TypeString,
	// 	Description: "IBM Cloud token URL.",
	// 	Computed:    true,
	// },
	// "ibm_cloud_account": {
	// 	Type:        schema.TypeString,
	// 	Description: "IBM Cloud account.",
	// 	Computed:   true,
	// },
	"status": {
		Type:        schema.TypeString,
		Description: "Status of the registry.",
		Computed:    true,
	},
	"error_message": {
		Type:        schema.TypeString,
		Description: "Error message associated with the registry.",
		Computed:    true,
	},
	"error_detail": {
		Type:        schema.TypeString,
		Description: "Detailed error information.",
		Computed:    true,
	},
	"started_at": {
		Type:        schema.TypeString,
		Description: "Start time of the registry.",
		Computed:    true,
	},
	"scanned": {
		Type:        schema.TypeInt,
		Description: "Number of items scanned in the registry.",
		Computed:    true,
	},
	"scheduled": {
		Type:        schema.TypeInt,
		Description: "Number of items scheduled for scanning.",
		Computed:    true,
	},
	"scanning": {
		Type:        schema.TypeInt,
		Description: "Number of items currently being scanned.",
		Computed:    true,
	},
	"failed": {
		Type:        schema.TypeInt,
		Description: "Number of items that failed scanning.",
		Computed:    true,
	},
	"cvedb_version": {
		Type:        schema.TypeString,
		Description: "CVE database version.",
		Computed:    true,
	},
	"cvedb_create_time": {
		Type:        schema.TypeString,
		Description: "Creation time of the CVE database.",
		Computed:    true,
	},
}

func DataSourceRegistry() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRegistryRead,
		Schema:      dataRegistrySchema,
	}
}

func dataSourceRegistryRead(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	APIClient := meta.(*goneuvector.Client)

	name := d.Get("name").(string)
	registrySummary, err := APIClient.GetRegistry(name)

	if err != nil {
		return diag.FromErr(err)
	}

	registry := registrySummary.Registry

	if err := helper.DataFromStruct(registry, d); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(registry.Name)
	d.Set("filters", registry.Filters)

	return nil
}
