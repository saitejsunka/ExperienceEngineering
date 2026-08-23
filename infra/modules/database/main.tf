# Service Networking Connection for Private IP
resource "google_compute_global_address" "private_ip_address" {
  name          = "google-managed-services-dbexp"
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

# Primary Database (us-west1)
resource "google_sql_database_instance" "primary" {
  name             = "dbexp-write"
  region           = "us-west1"
  database_version = "POSTGRES_15"
  
  # For ease of learning/teardown. In prod, set to true.
  deletion_protection = false

  settings {
    tier = "db-f1-micro"
    
    ip_configuration {
      ipv4_enabled    = false
      private_network = var.vpc_id
    }
    
    disk_size = 10
  }
  
  depends_on = [google_service_networking_connection.private_vpc_connection]
}

# Replica 1 (us-west1)
resource "google_sql_database_instance" "replica_west" {
  name                 = "dbexp-read-west"
  master_instance_name = google_sql_database_instance.primary.name
  region               = "us-west1"
  database_version     = "POSTGRES_15"
  
  deletion_protection = false

  settings {
    tier = "db-f1-micro"
    ip_configuration {
      ipv4_enabled    = false
      private_network = var.vpc_id
    }
    disk_size = 10
  }
}

# Replica 2 (us-east1)
resource "google_sql_database_instance" "replica_east" {
  name                 = "dbexp-read-east"
  master_instance_name = google_sql_database_instance.primary.name
  region               = "us-east1"
  database_version     = "POSTGRES_15"
  
  deletion_protection = false

  settings {
    tier = "db-f1-micro"
    ip_configuration {
      ipv4_enabled    = false
      private_network = var.vpc_id
    }
    disk_size = 10
  }
}
