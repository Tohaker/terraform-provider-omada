data "omada_sites" "all" {}

output "first_site" {
  value = data.omada_sites.all.sites[0].name
}