# ==============================================================================
# Database Configuration (db.tf)
# ==============================================================================
# Instantiates the Primary/Replica PostgreSQL Database setup.

module "database" {
  source     = "./modules/database"
  vpc_id     = module.networking.vpc_id
  region     = var.gcp_region
  project_id = var.gcp_project_id
}
