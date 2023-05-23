// resource_admission_rule.go
package neuvector

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/theobori/go-neuvector/client"
	"github.com/theobori/go-neuvector/controller/admission"
)

func resourceAdmissionRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAdmissionRuleCreate,
		ReadContext:   resourceAdmissionRuleRead,
		UpdateContext: resourceAdmissionRuleUpdate,
		DeleteContext: resourceAdmissionRuleDelete,

		Schema: map[string]*schema.Schema{
			"rule_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Admission rule ID",
			},
			"category": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Platform category, example `Kubernetes`",
			},
			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Comment",
			},
			"criteria": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"op": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"template_kind": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"path": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"value_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"disable": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Rule restriction",
			},
			"cfg_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default: "user_created",
			},
			"rule_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rule_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAdmissionRuleCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	APIClient := meta.(*client.Client)
	
	ruleMode := d.Get("rule_mode").(string)
	criteriasRaw := d.Get("criteria").(*schema.Set).List()
	
	// Collecting criterias
	var criterias []admission.AdmissionRuleCriterion
	
	for _, criteriaRaw := range criteriasRaw {
		criteriaRawMap := criteriaRaw.(map[string]any)

		criteria := admission.AdmissionRuleCriterion{
			Name: criteriaRawMap["name"].(string),
			Op: criteriaRawMap["op"].(string),
			Value: criteriaRawMap["value"].(string),
			Path: criteriaRawMap["path"].(string),
		}

		criterias = append(criterias, criteria)
	}

	body := admission.CreateAdmissionRuleBody{
		ID: d.Get("rule_id").(int),
		Category: d.Get("category").(string),
		Comment: d.Get("comment").(string),
		RuleType: d.Get("rule_type").(string),
		CfgType: d.Get("cfg_type").(string),
		RuleMode: &ruleMode,
		Disable: d.Get("disable").(bool),
		Criteria: criterias,
	}

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

func resourceAdmissionRuleUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	return nil
}

func resourceAdmissionRuleDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	d.SetId("")
	
	return nil
}