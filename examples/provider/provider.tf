terraform {
  required_providers {
    omada = {
      source = "Tohaker/omada"
    }
  }
}

variable "client_secret" {
  type      = string
  sensitive = true
}


provider "omada" {
  # Your software controller is hosted in the US region
  host          = "https://use1-omada-cloud.tplinkcloud.com"
  controller_id = "example-controller-id"
  client_id     = "example-client-id"

  # Alternatively, omit this field and supply it securely 
  # with the OMADA_CLIENT_SECRET environment variable
  client_secret = var.client_secret

  # Explicitly set the TLS verification setting
  tls_skip_verify = false
}
