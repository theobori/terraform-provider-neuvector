---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "neuvector_policy_ids Data Source - terraform-provider-neuvector"
subcategory: ""
description: |-
  
---

# neuvector_policy_ids (Data Source)



## Example Usage

```terraform
data "neuvector_policy_ids" "test" {}

# data "neuvector_policy_ids" "from_containers" {
#   from         = "containers"
#   to           = "containers"
#   ports        = "any"
#   applications = ["HTTP", "MySQL"]
# }

# data "neuvector_policy_ids" "http" {
#   applications = ["HTTP"]
# }

# data "neuvector_policy_ids" "http" {
#   applications = ["HTTP"]
# }

# data "neuvector_policy_ids" "federation_rules" {
#   cfg_type = "federal"
# }
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `action` (String) Used to filter. Action when this policy is triggered.
- `applications` (List of String) Used to filter. Enter applications for NeuVector to allow or deny. NeuVector understands deep application behavior and will analyze the payload to determine application protocols. Protocols include HTTP, HTTPS, SSL, SSH, DNS, DNCP, NTP, TFTP, ECHO, RTSP, SIP, MySQL, Redis, Zookeeper, Cassandra, MongoDB, PostgresSQL, Kafka, Couchbase, ActiveMQ, ElasticSearch, RabbitMQ, Radius, VoltDB, Consul, Syslog, Etcd, Spark, Apache, Nginx, Jetty, NodeJS, Oracle, MSSQL, Memcached and gRPC. To select everything enter "any"
- `cfg_type` (String) Used to filter. The type of configuration, its scope, for example whether the rule applies to the whole federation or just to the cluster.
- `disable` (Boolean) Used to filter. Disable the policy.
- `from` (String) Used to filter. Specify the group from where the connection will originate.
- `learned` (Boolean) Used to filter. Indicates if the rules has been learned.
- `ports` (String) Used to filter. If there are specific ports to limit this rule to, enter them here. For ICMP traffic, enter icmp. Sample: 80,tcp/8080,udp/6142-6150,tcp/any,udp/any,icmp,any
- `priority` (Number) Used to filter. The rule priority level.
- `to` (String) Used to filter. Specify the destination GROUP where these connections are allowed or denied.

### Read-Only

- `id` (String) The ID of this resource.
- `ids` (Set of Number) List of every policy ID.

