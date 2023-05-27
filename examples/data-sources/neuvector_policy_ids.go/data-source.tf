data "neuvector_policy_ids" "from_containers" {
  from = "containers"
  to = "containers"
  ports = "any"
  applications = [ "HTTP", "MySQL" ]
}

data "neuvector_policy_ids" "http" {
  applications = [ "HTTP" ]
}

data "neuvector_policy_ids" "http" {
  applications = [ "HTTP" ]
}

data "neuvector_policy_ids" "federation_rules" {
  cfg_type = "federal"
}
