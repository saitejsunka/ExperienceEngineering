output "primary_private_ip" {
  value       = google_sql_database_instance.primary.private_ip_address
  description = "The private IP address of the primary database instance"
}

output "db_user" {
  value       = google_sql_user.app_user.name
  description = "The database application username"
}

output "db_password" {
  value       = google_sql_user.app_user.password
  description = "The database application password"
  sensitive   = true
}
