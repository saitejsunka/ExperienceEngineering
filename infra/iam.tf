# ==============================================================================
# Service Account (iam.tf)
# ==============================================================================
# Provisions a dedicated Service Account for the Cloud Run backend to run as,
# and grants it the principle of least privilege (only what it needs).

resource "google_service_account" "cloudrun_sa" {
  account_id   = "dbexp-run-sa"
  display_name = "Cloud Run Service Account for DBExp"
  project      = var.gcp_project_id

  depends_on = [
    google_project_service.enabled_apis
  ]
}

# Grant the Cloud Run service account access to read the database config secret
resource "google_secret_manager_secret_iam_member" "secret_accessor" {
  project   = var.gcp_project_id
  secret_id = google_secret_manager_secret.dbexp_config.id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${google_service_account.cloudrun_sa.email}"
}
