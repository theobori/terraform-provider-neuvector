// resource_user.go
package neuvector

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	goneuvector "github.com/theobori/go-neuvector/neuvector"
	"github.com/theobori/terraform-provider-neuvector/internal/helper"
)

var resourceUserSchema = map[string]*schema.Schema{
	"fullname": {
		Type:        schema.TypeString,
		Description: "The full name of the user.",
		Required:    true,
	},
	"server": {
		Type:        schema.TypeString,
		Description: "The server associated with the user.",
		Optional:    true,
	},
	"username": {
		Type:        schema.TypeString,
		Description: "The username of the user.",
		Required:    true,
	},
	"password": {
		Type:        schema.TypeString,
		Description: "The password of the user.",
		Sensitive:   true,
		Optional:    true,
	},
	"email": {
		Type:        schema.TypeString,
		Description: "The email address of the user.",
		Optional:    true,
	},
	"role": {
		Type:        schema.TypeString,
		Description: "The role of the user.",
		Required:    true,
	},
	"timeout": {
		Type:        schema.TypeInt,
		Description: "The timeout value for the user session.",
		Optional:    true,
		Default:     300,
	},
	"locale": {
		Type:        schema.TypeString,
		Description: "The locale setting for the user.",
		Optional:    true,
		Default:     "en",
	},
	"default_password": {
		Type:        schema.TypeBool,
		Description: "Flag indicating if the user is using the default password.",
		Optional:    true,
		Default:     false,
	},
	"modify_password": {
		Type:        schema.TypeBool,
		Description: "Flag indicating if the user can modify the password.",
		Optional:    true,
		Default:     false,
	},
	// "role_domains_role": {
	// 	Type:        schema.TypeString,
	// 	Description: "The role associated with the domain.",
	// 	Optional:    true,
	// },
	// "role_domains": {
	// 	Type:        schema.TypeList,
	// 	Description: "List of domains associated with the role.",
	// 	Optional:    true,
	// 	Elem: &schema.Schema{
	// 		Type:        schema.TypeString,
	// 		Description: "A domain associated with the role.",
	// 	},
	// },
	"blocked_for_failed_login": {
		Type:        schema.TypeBool,
		Description: "Flag indicating if the user is blocked for failed login attempts.",
		Optional:    true,
		Default:     false,
	},
	"blocked_for_password_expired": {
		Type:        schema.TypeBool,
		Description: "Flag indicating if the user is blocked due to an expired password.",
		Optional:    true,
		Default:     false,
	},
}

func ResourceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		DeleteContext: resourceUserDelete,
		UpdateContext: resourceUserUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: resourceUserSchema,
	}
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	APIClient := meta.(*goneuvector.Client)

	body := helper.FromSchemas[goneuvector.User](resourceUserSchema, d)

	// TODO: Supports role domains

	if err := APIClient.CreateUser(body); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(body.Fullname)

	return resourceUserRead(ctx, d, meta)
}

func resourceUserUpdate(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}

func resourceUserRead(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	APIClient := meta.(*goneuvector.Client)

	user, err := APIClient.GetUser(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	// Because NeuVector hides the password and we do not want
	// changes in the Terraform state.
	password := d.Get("password").(string)

	if err := helper.TfFromStruct(user.User, d, true); err != nil {
		return diag.FromErr(err)
	}

	d.Set("password", password)

	return nil
}

func resourceUserDelete(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	APIClient := meta.(*goneuvector.Client)

	if err := APIClient.DeleteUser(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
