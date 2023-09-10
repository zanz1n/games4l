variable "aws_access_key" {
  type        = string
  description = "aws access key to deploy the project"
  sensitive   = true
}

variable "aws_secret_key" {
  type        = string
  description = "aws secret key to deploy the project"
  sensitive   = true
}

variable "aws_route53_zoneid" {
  type        = string
  description = "aws route53 domain zoneid"
}

variable "aws_route53_root_domain" {
  type        = string
  description = "the root domain that the app will use"
}

variable "telemetry_mongo_database_uri" {
  type        = string
  description = "mongodb database uri"
  sensitive   = true
}

variable "telemetry_mongo_database_name" {
  type        = string
  description = "mongodb database name"
  sensitive   = true
}

variable "users_mongo_database_uri" {
  type        = string
  description = "mongodb database uri"
  sensitive   = true
}

variable "users_mongo_database_name" {
  type        = string
  description = "mongodb database name"
  sensitive   = true
}

variable "postgres_database_uri" {
  type        = string
  description = "postgresql uri"
  sensitive   = true
}

variable "webhook_signature" {
  type        = string
  description = "webhook signature"
  sensitive   = true
}

variable "environment_type" {
  type        = string
  description = "dev or prod"
  default     = "dev"
}

variable "jwt_signature" {
  type        = string
  description = "jwt signature"
  sensitive   = true
}

locals {
  mime_types = jsondecode(file("./0-mime.json"))
}
