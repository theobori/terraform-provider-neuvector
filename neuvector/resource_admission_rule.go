// resource_admission_rule.go
package neuvector

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/theobori/go-neuvector/client"
	"github.com/theobori/go-neuvector/controller/admission"
	"github.com/theobori/terraform-provider-neuvector/helper"
)

var resourceAdmissionRuleSchema = map[string]*schema.Schema{
	"rule_id": {
		Type:        schema.TypeInt,
		Optional:    true,
		Description: "Represents the admission rule unique ID.",
		Default:     0,
	},
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
					Type:     schema.TypeString,
					Required: true,
					Description: "Represents the identifier or label for the criteria.",
				},
				"op": {
					Type:     schema.TypeString,
					Required: true,
					Description: "Represents a comparison operator used to evaluate the rule criteria. It defines the relationship between the specified value and the actual value being checked against. Examples of common operators include equals (=), not equals (!=), greater than (>), less than (<), etc.",
				},
				"value": {
					Type:     schema.TypeString,
					Required: true,
					Description: "Represents the reference value against which the actual value is compared using the specified operator.",
				},
				"type": {
					Type:     schema.TypeString,
					Optional: true,
					Description: "The type field defines the category or nature of the admission rule criteria. It helps determine the context and behavior of the criteria.",
				},
				"template_kind": {
					Type:     schema.TypeString,
					Optional: true,
					Description: "identifies the type or category of the admission rule template associated with the criteria.",
				},
				"path": {
					Type:     schema.TypeString,
					Optional: true,
					Description: "Specifies the location or attribute within the relevant resource or object that the admission rule criteria should be applied to.",
				},
				"value_type": {
					Type:     schema.TypeString,
					Optional: true,
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

func resourceAdmissionRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAdmissionRuleCreate,
		ReadContext:   resourceAdmissionRuleRead,
		UpdateContext: resourceAdmissionRuleUpdate,
		DeleteContext: resourceAdmissionRuleDelete,

		Schema: resourceAdmissionRuleSchema,
	}
}

func resourceAdmissionRuleCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	APIClient := meta.(*client.Client)

	criteriasRaw := d.Get("criteria").(*schema.Set).List()
	criterias := helper.FromTypeSetDefault[admission.AdmissionRuleCriterion](criteriasRaw)

	// Injecting Terraform data into a struct used as HTTP request body
	body := helper.FromSchemas[admission.CreateAdmissionRuleBody](
		resourceAdmissionRuleSchema,
		d,
	)

	body.ID = d.Get("rule_id").(int)
	body.Criteria = criterias

	rule, err := admission.CreateAdmissionRule(
		APIClient,
		body,
	)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(rule.ID))

	return resourceAdmissionRuleRead(ctx, d, meta)
}

func resourceAdmissionRuleRead(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

func resourceAdmissionRuleDelete(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	APIClient := meta.(*client.Client)

	var err error

	ruleId, err := strconv.Atoi(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	err = admission.DeleteAdmissionRule(
		APIClient,
		ruleId,
	)

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
