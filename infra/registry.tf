# ==============================================================================
# Artifact Registry (registry.tf)
# ==============================================================================
# Provisions a Docker container registry to store backend images for Cloud Run.

resource "google_artifact_registry_repository" "backend_repo" {
  location      = "us-west1" # Should match the Cloud Run deployment region
  repository_id = "dbexp-repo"
  description   = "Docker repository for DBPlay backend images"
  format        = "DOCKER"

  cleanup_policies {
    id     = "keep-top-20"
    action = "KEEP"
    most_recent_versions {
      keep_count = 20
    }
  }

  depends_on = [
    google_project_service.enabled_apis
  ]
}
