resource "neuvector_admission_rule" "containers_prevention" {
  rule_type = "deny"
  category  = "Docker"
  comment   = "Containers prevention"

  criteria {
    name  = "runAsRoot"
    op    = "="
    path  = "runAsRoot"
    value = "true"
  }

  disable = false
}

resource "neuvector_admission_rule" "containers_prevjjention" {
  rule_type = "deny"
  category  = "Docker"
  comment   = "Containers prevention"

  criteria {
    name  = "runAsRoot"
    op    = "="
    path  = "runAsRoot"
    value = "true"
  }

  criteria {
    name  = "runAsPrivileged"
    op    = "="
    path  = "runAsPrivileged"
    value = "true"
  }

  disable = false
}

resource "neuvector_promote" "promote_server" {
  port   = 11443
  server = "localhost"
  user   = "admin"
  name   = "cluster.local"
}
