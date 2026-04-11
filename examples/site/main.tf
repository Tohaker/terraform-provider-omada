terraform {
  required_providers {
    omada = {
      source = "registry.terraform.io/tohaker/omada"
    }
  }
}

provider "omada" {}

resource "omada_site" "example" {
  name      = "Terraform Example"
  region    = "United Kingdom"
  time_zone = "Europe/London"
  scenario  = "Home"

  device_account_setting = {
    username = "admin"
    password = "VeryStr@ngPas5w0rd"
  }
}

output "example_sites" {
  value     = omada_site.example
  sensitive = true
}
