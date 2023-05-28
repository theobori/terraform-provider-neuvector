resource "neuvector_promote" "test" {
  port   = 11443
  server = "localhost"
  user   = "admin"
  name   = "cluster.local"
}
