resource "google_firestore_database" "database" {
  project     = var.project_id
  name        = "dbexp-id"
  location_id = var.location_id
  type        = "FIRESTORE_NATIVE"
  
  # Disabling delete protection for easy teardown in a learning project.
  # In production, this should be set to DELETE_PROTECTION_ENABLED.
  delete_protection_state = "DELETE_PROTECTION_DISABLED"
}
