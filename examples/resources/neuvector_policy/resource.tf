resource "neuvector_policy" "test" {
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
    policy_id    = 123
    action       = "deny"
    applications = ["HTTP", "Redis"]
    comment      = "Nodes web constraints"
    from         = "nodes"
    to           = "containers"
    ports        = "tcp/80"
  }

  rule {
    action       = "deny"
    applications = ["any"]
    comment      = "Nodes web constraints"
    from         = "containers"
    to           = "nodes"
    ports        = "tcp/80"
  }
}

# resource "neuvector_policy" "fed_containers" {
#   rules_scope = "federal"
#   rule {
#     action       = "deny"
#     applications = ["any"]
#     comment      = "Containers constraints"
#     disable      = false
#     from         = "fed.containers"
#     to           = "fed.containers"
#     learned      = false
#     ports        = "any"
#     priority     = 0
#     cfg_type     = "federal"
#   }
# }