resource "neuvector_admission_rule" "foo" {
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
}
