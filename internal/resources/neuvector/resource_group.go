// resource_group.go
package neuvector

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	goneuvector "github.com/theobori/go-neuvector/neuvector"
	"github.com/theobori/terraform-provider-neuvector/internal/helper"
)

var resourceGroupSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Name of the group.",
	},
	"criteria": {
		Type:        schema.TypeSet,
		Required:    true,
		Description: "Matching criteria applied associated with the group",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"key": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Represents the identifier or label for the criteria.",
				},
				"op": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Represents a comparison operator used to evaluate the rule criteria. It defines the relationship between the specified value and the actual value being checked against. Examples of common operators include equals (=), not equals (!=), greater than (>), less than (<), etc.",
				},
				"value": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Represents the reference value against which the actual value is compared using the specified operator.",
				},
			},
		},
	},
	"cfg_type": {
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "user_created",
		Description: "The type of configuration, its scope, for example whether the rule applies to the whole federation or just to the cluster.",
	},
}

func ResourceGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupCreate,
		ReadContext:   resourceGroupRead,
		DeleteContext: resourceGroupDelete,
		UpdateContext: resourceGroupUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: resourceGroupSchema,
	}
}

func readGroup(d *schema.ResourceData) *goneuvector.CreateGroupBody {
	group := helper.FromSchemas[goneuvector.CreateGroupBody](
		resourceGroupSchema,
		d,
	)

	criteriaRaw := d.Get("criteria").(*schema.Set).List()
	criteria := helper.FromTypeSetDefault[goneuvector.GroupCriteria](
		criteriaRaw,
	)

	group.Criteria = criteria

	return &group
}

func getGroupCriteria(criterias *[]goneuvector.GroupCriteria) []map[string]any {
	var ret []map[string]any

	for _, c := range *criterias {
		if s, _ := helper.StructToMap(c); s != nil {
			ret = append(ret, s)
		}
	}

	return ret
}

func resourceGroupCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	APIClient := meta.(*goneuvector.Client)

	group := readGroup(d)

	if err := APIClient.WithContext(ctx).CreateGroup(*group); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(group.Name)

	return nil
}

func resourceGroupUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	if d.HasChange("name") {
		return diag.Errorf("You are not allowed to change the group name.")
	}

	APIClient := meta.(*goneuvector.Client)

	group := readGroup(d)

	if err := APIClient.WithContext(ctx).PatchGroup(group.Name, *group); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceGroupRead(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	APIClient := meta.(*goneuvector.Client)

	groupData, err := APIClient.GetGroup(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	group := groupData.Group

	if err := helper.TfFromStruct(group, d, true); err != nil {
		return diag.FromErr(err)
	}

	d.Set(
		"criteria",
		getGroupCriteria(&group.Criteria),
	)

	return nil
}

func resourceGroupDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	APIClient := meta.(*goneuvector.Client)

	err := APIClient.
		WithContext(ctx).
		DeleteGroup(
			d.Get("name").(string),
		)

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
