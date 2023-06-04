// resourceuser_role.go
package neuvector

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	goneuvector "github.com/theobori/go-neuvector/neuvector"
	"github.com/theobori/terraform-provider-neuvector/internal/helper"
)

var resourceUserRoleSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Description: "The name of the role.",
		Required:    true,
	},
	"comment": {
		Type:        schema.TypeString,
		Description: "The comment of the role.",
		Optional:    true,
	},
	"permission": {
		Type:        schema.TypeSet,
		Description: "The permissions associated to the role.",
		Required:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Represents the identifier or label for the criteria.",
				},
				"read": {
					Type:        schema.TypeBool,
					Required:    true,
					Description: "Flag indicating if the role has the read permission.",
				},
				"write": {
					Type:        schema.TypeBool,
					Required:    true,
					Description: "Flag indicating if the role has the write permission.",
				},
			},
		},
	},
}

func ResourceUserRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserRoleCreate,
		ReadContext:   resourceUserRoleRead,
		DeleteContext: resourceUserRoleDelete,
		UpdateContext: resourceUserRoleUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: resourceUserRoleSchema,
	}
}

func getPermissions(permissions []goneuvector.UserRolePermission) []map[string]any {
	var ret []map[string]any

	for _, p := range permissions {
		_map, err := helper.StructToMap(p)

		if err != nil {
			continue
		}

		ret = append(ret, _map)
	}

	return ret
}

func readUserRole(d *schema.ResourceData) goneuvector.CreateUserRoleBody {
	role := helper.FromSchemas[goneuvector.CreateUserRoleBody](resourceUserRoleSchema, d)

	permissionsRaw := d.Get("permission").(*schema.Set).List()
	permissions := helper.FromTypeSetDefault[goneuvector.UserRolePermission](permissionsRaw)

	role.Permissions = permissions

	return role
}

func resourceUserRoleCreate(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	APIClient := meta.(*goneuvector.Client)

	body := readUserRole(d)

	if err := APIClient.CreateUserRole(body); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(body.Name)

	return nil
}

func resourceUserRoleUpdate(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	if d.HasChange("name") {
		return diag.Errorf("You are not allowed to change the role name.")
	}
	
	APIClient := meta.(*goneuvector.Client)

	body := readUserRole(d)

	if err := APIClient.PatchUserRole(body.Name, body); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceUserRoleRead(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	APIClient := meta.(*goneuvector.Client)

	roleFull, err := APIClient.GetUserRole(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	role := roleFull.Role

	if err := helper.TfFromStruct(role, d, true); err != nil {
		return diag.FromErr(err)
	}

	d.Set("permission", getPermissions(role.Permissions))

	return nil
}

func resourceUserRoleDelete(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	APIClient := meta.(*goneuvector.Client)

	if err := APIClient.DeleteUserRole(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
