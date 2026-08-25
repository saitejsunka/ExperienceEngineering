# ==============================================================================
# Observability (observability.tf)
# ==============================================================================
# Manages Cloud Logging and Monitoring configurations.

# Update the default log bucket to retain logs for only 1 day to save costs.
resource "google_logging_project_bucket_config" "default" {
  project        = var.gcp_project_id
  location       = "global"
  retention_days = 1
  bucket_id      = "_Default"

  depends_on = [
    google_project_service.enabled_apis
  ]
}
