# ==============================================================================
# Database Configuration (db.tf)
# ==============================================================================
# Instantiates the Primary and Replica PostgreSQL Database setup.

module "database" {
  source     = "./modules/database"
  vpc_id     = module.networking.vpc_id
  project_id = var.gcp_project_id
}
