resource "neuvector_admission_rule" "container_prevention" {
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

resource "neuvector_registry" "docker_registry" {
  name                   = "docker-registry"
  registry_type          = "Docker Registry"
  username               = "test"
  password               = "test"
  filters                = ["*"]
  registry               = "https://registry.hub.docker.com/"
  rescan_after_db_update = true
  auth_with_token        = false
  scan_layers            = false
}

resource "neuvector_promote" "promote_server" {
  port   = 11443
  server = "localhost"
  user   = "admin"
  name   = "cluster.local"

  depends_on = [
    neuvector_admission_rule.container_prevention,
    neuvector_registry.docker_registry
  ]
}
