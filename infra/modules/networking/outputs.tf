output "vpc_id" {
  value       = google_compute_network.vpc_network.id
  description = "The ID of the VPC network"
}

output "network_name" {
  value       = google_compute_network.vpc_network.name
  description = "The name of the VPC network"
}
