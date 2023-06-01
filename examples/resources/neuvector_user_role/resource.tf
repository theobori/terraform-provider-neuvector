resource "neuvector_user_role" "roletest" {
  name = "roletest"

  permission {
    id    = "ci_scan"
    read  = false
    write = true
  }
}
