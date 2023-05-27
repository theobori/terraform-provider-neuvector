# resource "neuvector_admission_rule" "container_prevention" {
#   rule_type = "deny"
#   category  = "Docker"
#   comment   = "Containers prevention"

#   criteria {
#     name  = "runAsRoot"
#     op    = "="
#     path  = "runAsRoot"
#     value = "true"
#   }

#   criteria {
#     name  = "runAsPrivileged"
#     op    = "="
#     path  = "runAsPrivileged"
#     value = "true"
#   }

#   disable = false
# }

# resource "neuvector_registry" "docker" {
#   name                   = "docker.io"
#   registry_type          = "Docker Registry"
#   username               = "tes12"
#   password               = "test12"
#   filters                = ["neuvector/*"]
#   registry               = "https://registry.hub.docker.com/"
#   rescan_after_db_update = true
#   auth_with_token        = false
#   scan_layers            = true
# }

# data "neuvector_registry" "docker_metadata" {
#   name = "docker.io"
# }

# data "neuvector_registry_names" "registries" {
#   registry_type = "Dockker Registry"
# }

# data "neuvector_registry_names" "regiskkkktries" {
#   registry_type = "Docker Registry"
# }

# data "neuvector_registry_names" "regisaaatries" {
# }

# data "neuvector_registry" "docker_metadata" {
#   name = resource.neuvector_registry.docker.name

#   depends_on = [
#     resource.neuvector_registry.docker
#   ]
# }

# data "neuvector_policy_ids" "containers" {
#   from = "containers"
#   to = "containers"
#   ports = "any"
#   applications = [ "HTTP", "MySQL" ]
# }

resource "neuvector_policy" "basic_preventions" {
  rule {
    action       = "deny"
    applications = ["any"]
    comment      = "Containers constraints"
    disable      = false
    from         = "containers"
    to           = "containers"
    learned      = false
    ports        = "any"
    priority     = 0
  }

  rule {
    action       = "deny"
    applications = ["HTTP"]
    comment      = "Nodes web constraints"
    from         = "nodes"
    to           = "containers"
    ports        = "tcp/80"
    priority     = 0
  }

  rule {
    policy_id    = 45
    action       = "deny"
    applications = ["Redis"]
    comment      = "Excluding external Redis connection to the containers"
    from         = "external"
    to           = "containers"
    ports        = "any"
    priority     = 0
  }
}

resource "neuvector_policy" "federation_network" {
  rules_scope = "federal"
  rule {
    action       = "deny"
    applications = ["any"]
    comment      = "Containers constraints"
    disable      = false
    from         = "fed.containers"
    to           = "fed.containers"
    learned      = false
    ports        = "any"
    priority     = 0
    cfg_type = "federal"
  }

}

# resource "neuvector_promote" "promote_server" {
#   port   = 11443
#   server = "localhost"
#   user   = "admin"
#   name   = "cluster.local"

#   depends_on = [
#     neuvector_admission_rule.container_prevention,
#     neuvector_registry.docker,
#     neuneuvector_policy.basic_preventions
#   ]
# }
