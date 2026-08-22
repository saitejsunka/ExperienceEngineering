variable "project_id" {
  description = "The GCP project ID"
  type        = string
}

variable "location_id" {
  description = "The location for the Firestore database (e.g., nam5 for North America multi-region)"
  type        = string
  default     = "nam5"
}
