resource "google_compute_network" "vpc_network" {
  name                    = var.vpc_name
  auto_create_subnetworks = false
  routing_mode            = "GLOBAL"
  description             = "VPC Network for DBPlay Architecture"
}

resource "google_compute_subnetwork" "regional_subnets" {
  for_each = toset(var.regions)

  name          = "${var.vpc_name}-${each.key}-subnet"
  ip_cidr_range = cidrsubnet(var.vpc_cidr_block, 4, index(var.regions, each.key))
  region        = each.key
  network       = google_compute_network.vpc_network.id

  # Best practice for future internal TCP/UDP Load Balancing or Private Service Connect
  private_ip_google_access = true
}
