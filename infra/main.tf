# ==============================================================================
# DBPlay Infrastructure Entrypoint (main.tf)
# Triggering GitHub Action Update
# ==============================================================================
# Core terraform block defining providers and backend.
# In Terraform, you don't need to explicitly "call" other .tf files here.
# Terraform automatically reads all .tf files (like vpc.tf and commons.tf) 
# in this directory and combines them during execution.
# ==============================================================================

terraform {
  # Remote backend for state management.
  # The bucket is automatically created and injected by GitHub Actions.
  backend "gcs" {
    bucket  = "__BACKEND_BUCKET_NAME__"
    prefix  = "terraform/state"
  }

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }
}
