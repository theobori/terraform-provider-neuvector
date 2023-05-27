// resource_policy.go
package neuvector

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/theobori/go-neuvector/client"
	"github.com/theobori/go-neuvector/controller/policy"
	"github.com/theobori/terraform-provider-neuvector/internal/helper"
)

// Default scope
const DefaultScope = "user_created"

var resourcePolicyRuleSchema = map[string]*schema.Schema{
	"policy_id": {
		Type:        schema.TypeInt,
		Optional:    true,
		Description: "Dont use this field if you want to generate a new ID.",
		Default:     -1,
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

// Specified ID when we want to create a dynamic policy ID
//
// In the resource creation process, we will first get every policies IDs from NeuVector,
// then we will use these IDs for find the available ones in the "user_created" range.
// We only get the needed amount.
const DynamicPolicyID = -1

// Read a policy rule
func readPolicyRule(_map map[string]any) (*policy.PolicyRule, error) {
	policy := helper.FromMap[policy.PolicyRule](_map)

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
func readPolicyRules(set []any) []policy.PolicyRule {
	return helper.FromTypeSetCallback(set, readPolicyRule)
}

// Omitting the `delete` field by choice
// because as a Terraform resource, it is not revelant to delete policies
// at the creation
var resourcePolicySchema = map[string]*schema.Schema{
	"rule": {
		Type:        schema.TypeSet,
		Required:    true,
		Description: "Matching criteria applied associated with the rule",
		Elem: &schema.Resource{
			Schema: resourcePolicyRuleSchema,
		},
	},
	"rules_scope": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Scope applied to every rules, it helps definin the url.",
		Default: DefaultScope,
	},
}

func ResourcePolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyCreate,
		ReadContext:   resourcePolicyRead,
		DeleteContext: resourcePolicyDelete,
		UpdateContext: resourcePolicyUpdate,

		Schema: resourcePolicySchema,
	}
}

// Return the indexes of the dynamic policies
func getDynamicPolicyIndexes(policies *[]policy.PolicyRule) []int {
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
	patchFunc func(
		client *client.Client,
		body policy.PatchPolicyBody,
	) error
	
}

// The scopes changes associated with specific scope name
var scopes = map[string]*ScopeChanges{
	DefaultScope: {
		policy.PolicyMinimumID+1,
		policy.PolicyMaximumID,
		policy.PatchPolicy,
	},
	"federal": {
		policy.FedPolicyMinimumID+1,
		policy.FedPolicyMaximumID,
		policy.PatchFedPolicy,
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
func patchPolicy(body *policy.PatchPolicyBody, APIClient *client.Client, scopeName string) error {
	var err error
	var params *ScopeChanges
	
	// Get the dynamic rules index in body.Rules
	// Used to determinate the amount of need available index
	indexes := getDynamicPolicyIndexes(&body.Rules)
	
	params = GetScopeChanges(scopeName)

	// Get every available policy IDs
	policyIDs, err := policy.GetPolicyAvailableIDs(
		APIClient,
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

	return params.patchFunc(APIClient, *body)
}

func createPolicy(d *schema.ResourceData, meta any) error {
	var err error

	APIClient := meta.(*client.Client)

	rulesRaw := d.Get("rule").(*schema.Set).List()
	body := policy.PatchPolicyBody{
		Rules: readPolicyRules(rulesRaw),
	}

	err = patchPolicy(
		&body,
		APIClient,
		d.Get("rules_scope").(string),
	)

	if err != nil {
		return err
	}

	// Random resource ID because this one doesnt have a specific ID
	// The only valid IDs are the ones in the `rule` field
	id, err := uuid.GenerateUUID()

	if err != nil {
		return err
	}

	d.SetId(id)
	
	return nil
}

func resourcePolicyCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	if err := createPolicy(d, meta); err != nil {
		return diag.FromErr(err)
	}

	return resourcePolicyRead(ctx, d, meta)
}

func resourcePolicyUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}

func resourcePolicyRead(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}

func resourcePolicyDelete(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var delete []int

	APIClient := meta.(*client.Client)

	// Magic trick
	//
	// NeuVector itself can identifies policy rules
	// with other fields than the unique ID. When we create a policy, NeuVector will check if
	// another one has the same behavior, if yes if will patch it instead of creating a new one.
	// We cannot retrieve the dynamics IDS injected into the HTTP body at the creation.
	//
	// So we are going to get every policies, then we will get the
	// dynamically created IDs by comparing the unique fields:
	// "from", "to", "ports", "applications", "learned", etc..
	//
	// Yes, even the field "comment" :]

	rulesRaw := d.Get("rule").(*schema.Set).List()
	resourcePolicyRules := readPolicyRules(rulesRaw)

	policies, err := policy.GetPolicies(APIClient)

	if err != nil {
		return diag.FromErr(err)
	}

	scope := d.Get("rules_scope").(string)
	params := GetScopeChanges(scope)

	for _, policyRule := range policies.Rules {
		for _, resourcePolicyRule := range resourcePolicyRules {
			if policyRule.Equal(&resourcePolicyRule) {
				delete = append(delete, policyRule.ID)
			}
		}
	}

	params.patchFunc(
		APIClient,
		policy.PatchPolicyBody{
			Delete: delete,
		},
	)

	return nil
}
