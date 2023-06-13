// data_source_policy_ids.go
package neuvector

import (
	"context"
	"reflect"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	goneuvector "github.com/theobori/go-neuvector/neuvector"
	"github.com/theobori/terraform-provider-neuvector/internal/helper"
)

var dataPolicyIDsSchema = map[string]*schema.Schema{
	"ids": {
		Type:        schema.TypeSet,
		Description: "List of every policy ID.",
		Computed:    true,
		Elem:        &schema.Schema{Type: schema.TypeInt},
	},
	"from": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Used to filter. Specify the group from where the connection will originate.",
	},
	"to": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Used to filter. Specify the destination GROUP where these connections are allowed or denied.",
	},
	"ports": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Used to filter. If there are specific ports to limit this rule to, enter them here. For ICMP traffic, enter icmp. Sample: 80,tcp/8080,udp/6142-6150,tcp/any,udp/any,icmp,any",
	},
	"action": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Used to filter. Action when this policy is triggered.",
	},
	"applications": {
		Type:        schema.TypeList,
		Optional:    true,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Description: "Used to filter. Enter applications for NeuVector to allow or deny. NeuVector understands deep application behavior and will analyze the payload to determine application protocols. Protocols include HTTP, HTTPS, SSL, SSH, DNS, DNCP, NTP, TFTP, ECHO, RTSP, SIP, MySQL, Redis, Zookeeper, Cassandra, MongoDB, PostgresSQL, Kafka, Couchbase, ActiveMQ, ElasticSearch, RabbitMQ, Radius, VoltDB, Consul, Syslog, Etcd, Spark, Apache, Nginx, Jetty, NodeJS, Oracle, MSSQL, Memcached and gRPC. To select everything enter \"any\"",
	},
	"learned": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Used to filter. Indicates if the rules has been learned.",
	},
	"disable": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Used to filter. Disable the policy.",
	},
	"priority": {
		Type:        schema.TypeInt,
		Optional:    true,
		Default:     0,
		Description: "Used to filter. The rule priority level.",
	},
	"cfg_type": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Used to filter. The type of configuration, its scope, for example whether the rule applies to the whole federation or just to the cluster.",
	},
}

func DataSourcePolicyIDs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePolicyIDsRead,
		Schema:      dataPolicyIDsSchema,
	}
}

func dataSourcePolicyIDsRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	APIClient := meta.(*goneuvector.Client)

	var ids []int

	policies, err := APIClient.
		WithContext(ctx).
		GetPolicies()

	if err != nil {
		return diag.FromErr(err)
	}

	applicationsRaw := d.Get("applications").([]any)
	applications, _ := helper.FromSlice[string](applicationsRaw)

	for _, p := range policies.Rules {
		// Comparing the basic type fields
		has, _ := helper.StructHasResource[goneuvector.PolicyRule](
			p,
			dataPolicyIDsSchema,
			d,
		)

		// Comparing applications
		if has || reflect.DeepEqual(p.Applications, applications) {
			ids = append(ids, p.ID)
		}
	}

	id, err := uuid.GenerateUUID()

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)
	d.Set("ids", ids)

	return nil
}
