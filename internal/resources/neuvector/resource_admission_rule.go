// resource_admission_rule.go
package neuvector

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	goneuvector "github.com/theobori/go-neuvector/neuvector"
	"github.com/theobori/terraform-provider-neuvector/internal/helper"
)

var resourceAdmissionRuleSchema = map[string]*schema.Schema{
	"category": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Represents an orchestration platform category, could be Kubernetes, OpenShift, Rancher, Docker, or other platforms.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "A comment from the user.",
	},
	"criteria": {
		Type:        schema.TypeSet,
		Required:    true,
		Description: "Matching criteria applied associated with the rule",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
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
				"type": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The type field defines the category or nature of the admission rule criteria. It helps determine the context and behavior of the criteria.",
				},
				"template_kind": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "identifies the type or category of the admission rule template associated with the criteria.",
				},
				"path": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Specifies the location or attribute within the relevant resource or object that the admission rule criteria should be applied to.",
				},
				"value_type": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Indicates the data type of the actual value being checked against the rule criteria.",
				},
			},
		},
	},
	"disable": {
		Type:        schema.TypeBool,
		Required:    true,
		Description: "Disable the rule.",
	},
	"cfg_type": {
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "user_created",
		Description: "The type of configuration, its scope, for example whether the rule applies to the whole federation or just to the cluster.",
	},
	"rule_type": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Indicate whether this rule is to \"allow\" this type of connection, or \"deny\" it.",
	},
	"rule_mode": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Defines the rigor with which rules are assessed and applied.",
	},
}

func ResourceAdmissionRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAdmissionRuleCreate,
		ReadContext:   resourceAdmissionRuleRead,
		UpdateContext: resourceAdmissionRuleUpdate,
		DeleteContext: resourceAdmissionRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: resourceAdmissionRuleSchema,
	}
}

func resourceAdmissionRuleCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	APIClient := meta.(*goneuvector.Client)

	criteriasRaw := d.Get("criteria").(*schema.Set).List()
	criterias := helper.FromTypeSetDefault[goneuvector.AdmissionRuleCriterion](criteriasRaw)

	body := helper.FromSchemas[goneuvector.CreateAdmissionRuleBody](
		resourceAdmissionRuleSchema,
		d,
	)

	body.Criteria = criterias

	rule, err := APIClient.
		WithContext(ctx).
		CreateAdmissionRule(body)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(rule.ID))

	return resourceAdmissionRuleRead(ctx, d, meta)
}

// Get criteria type set as a map from []admission.AdmissionRuleCriterion
func getCriteria(criterias []goneuvector.AdmissionRuleCriterion) []map[string]any {
	var ret []map[string]any

	for _, criteria := range criterias {
		_map, err := helper.StructToMap(criteria)

		if err != nil {
			continue
		}

		ret = append(ret, _map)
	}

	return ret
}

func resourceAdmissionRuleRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	APIClient := meta.(*goneuvector.Client)

	id, err := strconv.Atoi(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	adm, err := APIClient.
		WithContext(ctx).
		GetAdmissionRule(id)

	if err != nil {
		d.SetId("")
		return nil
	}

	rule := adm.Rule

	if err := helper.TfFromStruct(rule, d, true); err != nil {
		return diag.FromErr(err)
	}

	d.Set("criteria", getCriteria(rule.Criteria))

	return nil
}

func resourceAdmissionRuleUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	if !d.HasChanges(
		"criteria",
		"rule_type",
		"rule_mode",
		"disable",
		"category",
		"cfg_type",
	) {
		return nil
	}

	if diag := resourceAdmissionRuleDelete(ctx, d, meta); diag != nil {
		return diag
	}

	return resourceAdmissionRuleCreate(ctx, d, meta)
}

func resourceAdmissionRuleDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	APIClient := meta.(*goneuvector.Client)

	var err error

	id, err := strconv.Atoi(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	if err := APIClient.WithContext(ctx).DeleteAdmissionRule(id); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
