terraform {
    required_providers {
            neuvector = {
            version = "~> 1.0.0"
            source  = "github.com/theobori/neuvector"
        }
    }
}

provider "neuvector" {
    base_url = "https://localhost:10443/v1"
    username = "admin"
    password = "admin"
}
