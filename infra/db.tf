# ==============================================================================
# Database Configuration (db.tf)
# ==============================================================================
# Instantiates the Multi-Region Firestore Database setup.

module "firestore" {
  source      = "./modules/firestore"
  project_id  = var.gcp_project_id
  location_id = "nam5" # North America Multi-Region
}
