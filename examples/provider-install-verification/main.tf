terraform {
  required_providers {
    omada = {
      source = "registry.terraform.io/tohaker/omada"
    }
  }
}

provider "omada" {}

data "omada_sites" "example" {}

output "example_sites" {
  value = data.omada_sites.example
}
