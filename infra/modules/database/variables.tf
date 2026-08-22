variable "vpc_id" {
  type        = string
  description = "The ID of the VPC network where the database will reside"
}

variable "region" {
  type        = string
  description = "The primary region for the database"
}

variable "project_id" {
  type        = string
  description = "The GCP Project ID"
}
