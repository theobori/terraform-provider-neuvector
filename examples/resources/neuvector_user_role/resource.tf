resource "neuvector_user_role" "test" {
  name = "roletest"

  permission {
    id    = "ci_scan"
    read  = false
    write = true
  }
}
