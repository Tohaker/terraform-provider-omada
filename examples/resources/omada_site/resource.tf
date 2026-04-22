# Set a variable to store the site password securely
variable "site_password" {
  type      = string
  sensitive = true
}

resource "omada_site" "example" {
  name      = "Example"
  type      = 0
  region    = "United Kingdom"
  time_zone = "Europe/London"
  scenario  = "Home"
  tag_ids = [
    "tag_id_1",
    "tag_id_2"
  ]
  longitude = -0.124681
  latitude  = 51.500786
  address   = "123 Fake Street"

  device_account_setting = {
    username = "admin"
    password = var.site_password
  }
}

output "home_site_id" {
  value = omada_site.example.site_id
}