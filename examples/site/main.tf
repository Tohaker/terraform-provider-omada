terraform {
  required_providers {
    omada = {
      source = "registry.terraform.io/tohaker/omada"
    }
  }
}

provider "omada" {}

resource "omada_site" "home" {
  name      = "Home"
  region    = "United Kingdom"
  time_zone = "Europe/London"
  scenario  = "Home"

  device_account_setting = {
    username = "admin"
    password = var.site_password
  }
}
