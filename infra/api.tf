# ==============================================================================
# Enabled APIs (api.tf)
# ==============================================================================
# Ensures required GCP Services are enabled for the infrastructure to function.

locals {
  services = [
    "cloudresourcemanager.googleapis.com",
    "compute.googleapis.com",
    "servicenetworking.googleapis.com",
    "sqladmin.googleapis.com",
    "secretmanager.googleapis.com",
    "logging.googleapis.com",
    "monitoring.googleapis.com",
    "run.googleapis.com",
    "artifactregistry.googleapis.com",
  ]
}

resource "google_project_service" "enabled_apis" {
  for_each           = toset(local.services)
  project            = var.gcp_project_id
  service            = each.key
  disable_on_destroy = false
}
