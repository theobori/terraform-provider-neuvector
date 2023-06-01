resource "neuvector_service_config" "test" {
  services = [
    "test"
  ]

  not_scored = true
}
