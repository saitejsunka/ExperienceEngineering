# ==============================================================================
# Common Shared Configuration (commons.tf)
# ==============================================================================
# Contains shared variables, provider setup, and common locals used across resources.

variable "gcp_project_id" {
  type        = string
  description = "The GCP Project ID injected from GitHub Actions (TF_VAR_gcp_project_id)"
}

variable "gcp_region" {
  type        = string
  description = "The GCP Region"
  default     = "us-central1"
}

provider "google" {
  project = var.gcp_project_id
  region  = var.gcp_region
}

# Common labels/tags to be shared across all resources
locals {
  common_labels = {
    environment = "dev"
    project     = "dbplay"
    managed_by  = "terraform"
  }
}

variable "us_regions" {
  type        = list(string)
  description = "The target US regions for the global deployment"
  default     = ["us-east1", "us-central1", "us-west1"]
}
