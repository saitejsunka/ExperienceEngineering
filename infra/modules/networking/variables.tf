variable "vpc_name" {
  type        = string
  description = "The name of the VPC network"
  default     = "dbexp-vpc"
}

variable "project_id" {
  type        = string
  description = "The GCP project ID"
}

variable "regions" {
  type        = list(string)
  description = "The list of GCP regions to deploy subnets into"
}

variable "vpc_cidr_block" {
  type        = string
  description = "The base CIDR block for the VPC"
  default     = "10.0.0.0/16"
}
