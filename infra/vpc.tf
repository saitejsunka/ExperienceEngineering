# ==============================================================================
# VPC Configuration (vpc.tf)
# ==============================================================================
# This file instantiates the networking module.
# Note: Terraform automatically discovers this file when run from the infra/ directory.

module "networking" {
  source     = "./modules/networking"
  vpc_name   = "dbexp-vpc"
  project_id = var.gcp_project_id
  regions    = var.us_regions
}
