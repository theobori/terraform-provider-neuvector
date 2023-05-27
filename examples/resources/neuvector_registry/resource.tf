resource "neuvector_registry" "foo" {
  name                   = "docker.io"
  registry_type          = "Docker Registry"
  username               = "test"
  password               = "test"
  filters                = ["neuvector/*", "devops/*"]
  registry               = "https://registry.hub.docker.com/"
  rescan_after_db_update = true
  auth_with_token        = false
  scan_layers            = true
}