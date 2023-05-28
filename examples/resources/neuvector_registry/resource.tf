resource "neuvector_registry" "test" {
  name                   = "docker.io"
  registry_type          = "Docker Registry"
  filters                = ["neuvector/*"]
  registry               = "https://registry.hub.docker.com/"
  rescan_after_db_update = true
  auth_with_token        = false
  scan_layers            = true
}
