# ==============================================================================
# Database Module - High Availability PostgreSQL (Primary & Replica)
# ==============================================================================

# 1. Private Services Access (VPC Peering for Private IP)
resource "google_compute_global_address" "private_ip_address" {
  name          = "dbexp-private-ip"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = var.vpc_id
}

resource "google_service_networking_connection" "private_vpc_connection" {
  network                 = var.vpc_id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_ip_address.name]
}

# 2. Primary Database Instance (DBExp_Write)
resource "google_sql_database_instance" "primary" {
  name             = "dbexp-write"
  database_version = "POSTGRES_15"
  region           = var.region
  
  # Ensure the private VPC connection is established before creating the DB
  depends_on = [google_service_networking_connection.private_vpc_connection]

  settings {
    tier = "db-f1-micro" # Use smallest tier for dev, easily scaled up later
    ip_configuration {
      ipv4_enabled    = false # STRICTLY Private
      private_network = var.vpc_id
    }
    backup_configuration {
      enabled = true # Backups must be enabled for read replicas
    }
  }
}

# 3. Read Replica Instance (DBExp_Read)
resource "google_sql_database_instance" "replica" {
  name                 = "dbexp-read"
  master_instance_name = google_sql_database_instance.primary.name
  database_version     = "POSTGRES_15"
  region               = var.region

  settings {
    tier = "db-f1-micro"
    ip_configuration {
      ipv4_enabled    = false
      private_network = var.vpc_id
    }
  }
}

# 4. Database Creation
resource "google_sql_database" "database" {
  name     = "dbplay_db"
  instance = google_sql_database_instance.primary.name
}

# 5. Default User
resource "google_sql_user" "users" {
  name     = "dbplay_admin"
  instance = google_sql_database_instance.primary.name
  password = "SuperSecretPassword123!" # For production, inject this via Secret Manager
}
