resource "neuvector_admission_rule" "test" {
  rule_type = "deny"
  category  = "Kubernetes"
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
