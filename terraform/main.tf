terraform {
  required_providers {
    jellyfin = {
      source = "registry.terraform.io/TheWhale01/jellyfin"
    }
  }
}

provider "jellyfin" {
  endpoint = "http://localhost:8096"
  api_key = "this is an api key!"
}
