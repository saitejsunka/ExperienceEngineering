# Random Password Generation for the DB user
resource "random_password" "db_password" {
  length  = 16
  special = true
}

# Application User (Replicates to replicas automatically)
resource "google_sql_user" "app_user" {
  name     = "dbexp_user"
  instance = google_sql_database_instance.primary.name
  password = random_password.db_password.result
}
