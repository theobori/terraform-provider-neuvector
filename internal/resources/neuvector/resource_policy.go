// resource_policy.go
package neuvector

import (
	"context"
	"strconv"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	goneuvector "github.com/theobori/go-neuvector/neuvector"
	"github.com/theobori/go-neuvector/util"
	"github.com/theobori/terraform-provider-neuvector/internal/helper"
)

const (
	// Specified ID when we want to create a dynamic policy ID
	//
	// In the resource creation process, we will first get every policies IDs from NeuVector,
	// then we will use these IDs for find the available ones in the "user_created" range.
	// We only get the needed amount.
	DynamicPolicyID = -1
	// Default scope
	DefaultScope = "user_created"
)

var resourcePolicyRuleSchema = map[string]*schema.Schema{
	"policy_id": {
		Type:        schema.TypeInt,
		Optional:    true,
		Description: "Dont use this field if you want to generate a new ID.",
		Default:     DynamicPolicyID,
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "A comment from the user.",
	},
	"from": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specify the group from where the connection will originate.",
	},
	"to": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specify the destination GROUP where these connections are allowed or denied.",
	},
	"ports": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "If there are specific ports to limit this rule to, enter them here. For ICMP traffic, enter icmp. Sample: 80,tcp/8080,udp/6142-6150,tcp/any,udp/any,icmp,any",
	},
	"action": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Action when this policy is triggered.",
	},
	"applications": {
		Type:        schema.TypeList,
		Required:    true,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Description: "Enter applications for NeuVector to allow or deny. NeuVector understands deep application behavior and will analyze the payload to determine application protocols. Protocols include HTTP, HTTPS, SSL, SSH, DNS, DNCP, NTP, TFTP, ECHO, RTSP, SIP, MySQL, Redis, Zookeeper, Cassandra, MongoDB, PostgresSQL, Kafka, Couchbase, ActiveMQ, ElasticSearch, RabbitMQ, Radius, VoltDB, Consul, Syslog, Etcd, Spark, Apache, Nginx, Jetty, NodeJS, Oracle, MSSQL, Memcached and gRPC. To select everything enter \"any\"",
	},
	"learned": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Indicates if the rules has been learned.",
	},
	"disable": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Disable the policy.",
	},
	"priority": {
		Type:        schema.TypeInt,
		Optional:    true,
		Default:     0,
		Description: "The rule priority level.",
	},
	"cfg_type": {
		Type:        schema.TypeString,
		Optional:    true,
		Default:     DefaultScope,
		Description: "The type of configuration, its scope, for example whether the rule applies to the whole federation or just to the cluster.",
	},
}

// Read a policy rule
func readPolicyRule(_map map[string]any) (*goneuvector.PolicyRule, error) {
	policy := helper.FromMap[goneuvector.PolicyRule](_map)

	policyID := _map["policy_id"].(int)
	applicationsRaw := _map["applications"].([]any)

	applications, err := helper.FromSlice[string](applicationsRaw)

	if err != nil {
		return nil, err
	}

	policy.Applications = applications
	policy.ID = policyID

	return &policy, nil
}

// Only for `schema.TypeSet` with `schema.Resource`
func readPolicyRules(set []any) []goneuvector.PolicyRule {
	return helper.FromTypeSetCallback(set, readPolicyRule)
}

// Omitting the `delete` field by choice
// because as a Terraform resource, it is not revelant to delete policies
// at the creation
var resourcePolicySchema = map[string]*schema.Schema{
	"rule": {
		Type:        schema.TypeSet,
		Required:    true,
		Description: "Matching criteria applied associated with the rule.",
		Elem: &schema.Resource{
			Schema: resourcePolicyRuleSchema,
		},
	},
	"rules_scope": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Scope applied to every rules, it helps definin the url.",
		Default:     DefaultScope,
	},
	"policy_ids": {
		Type:        schema.TypeSet,
		Computed:    true,
		Description: "Contains every policy ID including the dynamic ones.",
		Elem:        &schema.Schema{Type: schema.TypeInt},
	},
}

func ResourcePolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyCreate,
		ReadContext:   resourcePolicyRead,
		DeleteContext: resourcePolicyDelete,
		UpdateContext: resourcePolicyUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: resourcePolicyImport,
		},

		Schema: resourcePolicySchema,
	}
}

// Return the indexes of the dynamic policies
func GetDynamicPolicyIndexes(policies *[]goneuvector.PolicyRule) []int {
	ret := []int{}

	for i, p := range *policies {
		if p.ID == DynamicPolicyID {
			ret = append(ret, i)
		}
	}

	return ret
}

// Stores the changing parameters between the differents scopes
type ScopeChanges struct {
	minID int
	maxID int
	IsFed bool
}

// The scopes changes associated with specific scope name
var scopes = map[string]*ScopeChanges{
	DefaultScope: {
		goneuvector.PolicyMinimumID + 1,
		goneuvector.PolicyMaximumID,
		false,
	},
	"federal": {
		goneuvector.FedPolicyMinimumID + 1,
		goneuvector.FedPolicyMaximumID,
		true,
	},
}

func GetScopeChanges(scopeName string) *ScopeChanges {
	ret, ok := scopes[scopeName]

	if !ok {
		return scopes[DefaultScope]
	}

	return ret
}

// Patch a policy rule, taking care of the scope
func patchPolicy(APIClient *goneuvector.Client, body *goneuvector.PatchPolicyBody, scopeName string) error {
	// Get the dynamic rules index in body.Rules
	// Used to determinate the amount of need available index
	indexes := GetDynamicPolicyIndexes(&body.Rules)
	params := GetScopeChanges(scopeName)

	// Get every available policy IDs
	policyIDs, err := APIClient.GetPolicyAvailableIDs(
		params.minID,
		params.maxID,
		len(indexes),
	)

	if err != nil {
		return err
	}

	// Updating the body with the new IDs if needed
	for i, index := range indexes {
		body.Rules[index].ID = policyIDs[i]
	}

	return APIClient.PatchPolicy(*body, params.IsFed)
}

func GetPolicyRuleMap(p *goneuvector.PolicyRule) *map[string]any {
	rule, err := helper.StructToMap(*p)

	if err != nil {
		return nil
	}

	delete(rule, "last_modified_timestamp")
	delete(rule, "created_timestamp")
	delete(rule, "id")

	rule["policy_id"] = p.ID
	rule["applications"] = p.Applications

	return &rule
}

func GetPolicyRulesSet(policies *[]goneuvector.PolicyRule) []map[string]any {
	var ret []map[string]any

	for _, p := range *policies {
		rule := GetPolicyRuleMap(&p)

		if rule == nil {
			return ret
		}

		ret = append(ret, *rule)
	}

	return ret
}

func resourcePolicyCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var err error

	APIClient := meta.(*goneuvector.Client)

	rulesRaw := d.Get("rule").(*schema.Set).List()
	body := goneuvector.PatchPolicyBody{
		Rules: readPolicyRules(rulesRaw),
	}

	APIClient.WithContext(ctx)

	// Patching policy handling the configuration scope
	err = patchPolicy(
		APIClient,
		&body,
		d.Get("rules_scope").(string),
	)

	if err != nil {
		return diag.FromErr(err)
	}

	// Random resource ID because this one doesnt have a specific ID
	// The only valid IDs are the ones in the `rule` field
	id, err := uuid.GenerateUUID()

	if err != nil {
		return diag.FromErr(err)
	}

	// Update policy_ids
	var policy_ids []int

	for _, rule := range body.Rules {
		policy_ids = append(policy_ids, rule.ID)
	}

	d.Set("rule", GetPolicyRulesSet(&body.Rules))
	d.Set("policy_ids", policy_ids)
	d.SetId(id)

	return resourcePolicyRead(ctx, d, meta)
}

func resourcePolicyUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}

func resourcePolicyRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var rules []map[string]any

	APIClient := meta.(*goneuvector.Client)

	policies, err := APIClient.
		WithContext(ctx).
		GetPolicies()

	if err != nil {
		return diag.FromErr(err)
	}

	policiesIDsRaw := d.Get("policy_ids").(*schema.Set).List()
	policiesIDs, err := helper.FromSlice[int](policiesIDsRaw)

	if err != nil {
		return diag.FromErr(err)
	}

	// Creating a set of rules to inject it into the resource
	for _, p := range policies.Rules {
		exists, _ := util.ItemExists(policiesIDs, p.ID)

		if !exists {
			continue
		}

		rule := GetPolicyRuleMap(&p)

		if rule == nil {
			continue
		}

		rules = append(rules, *rule)
	}

	d.Set("rule", rules)

	return nil
}

func resourcePolicyDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	APIClient := meta.(*goneuvector.Client)

	params := GetScopeChanges(d.Get("rules_scope").(string))

	deleteRaw := d.Get("policy_ids").(*schema.Set).List()
	delete, err := helper.FromSlice[int](deleteRaw)

	if err != nil {
		return diag.FromErr(err)
	}

	APIClient.
		WithContext(ctx).
		PatchPolicy(
			goneuvector.PatchPolicyBody{
				Delete: delete,
			},
			params.IsFed,
		)

	return nil
}

func resourcePolicyImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	APIClient := meta.(*goneuvector.Client)

	ruleID, err := strconv.Atoi(d.Id())

	if err != nil {
		return nil, err
	}

	p, err := APIClient.
		WithContext(ctx).
		GetPolicy(ruleID)

	if err != nil {
		return nil, err
	}

	rule := p.Rule
	id, err := uuid.GenerateUUID()

	if err != nil {
		return nil, err
	}

	d.Set("rules_scope", rule.CfgType)
	d.Set("policy_ids", []int{rule.ID})
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}
