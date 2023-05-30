resource "neuvector_group" "test" {
  name = "mytestgroup"

  criteria {
    key   = "pattern"
    value = "[a-z]"
    op    = "regex"
  }
}
