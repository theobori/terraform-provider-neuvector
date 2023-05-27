resource "neuvector_promote" "foo" {
  port   = 11443
  server = "localhost"
  user   = "admin"
  name   = "cluster.local"
}