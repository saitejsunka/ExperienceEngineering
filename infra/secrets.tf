# ==============================================================================
# Secret Manager (secrets.tf)
# ==============================================================================
# Provisions GCP Secret Manager and populates it with the database credentials.

resource "google_secret_manager_secret" "dbexp_config" {
  secret_id = "dbexp-backend-config"
  project   = var.gcp_project_id

  replication {
    auto {}
  }

  depends_on = [
    google_project_service.enabled_apis
  ]
}

resource "google_secret_manager_secret_version" "dbexp_config_version" {
  secret = google_secret_manager_secret.dbexp_config.id
  
  # Inject the database credentials as a JSON payload matching the AppConfig struct.
  # The jsonencode function safely escapes strings.
  secret_data = jsonencode({
    DatabaseHost     = module.database.primary_private_ip
    DatabasePort     = "5432"
    DatabaseUser     = module.database.db_user
    DatabasePassword = module.database.db_password
    DatabaseName     = "postgres" # Default DB created by Cloud SQL
  })
}
