variable "site_password" {
  type = string
}

resource "omada_site" "example" {
  name      = "Example"
  region    = "United Kingdom"
  time_zone = "Europe/London"
  scenario  = "Home"

  device_account_setting = {
    username = "admin"
    password = var.site_password
  }
}