terraform {
  required_providers {
    omada = {
      source = "registry.terraform.io/tohaker/omada"
    }
  }
}

provider "omada" {}

data "omada_site" "example" {}
