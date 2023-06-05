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

variable "cloudflare_api_token" {
  type        = string
  description = "cloudflare api token to setup the domain"
  sensitive   = true
}

variable "cloudflare_zoneid" {
  type        = string
  description = "the id of domain that will be used in production"
  sensitive   = true
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

variable "apigateway_cloudflare_domain" {
  type        = string
  description = "the domain (with the root and the subdomains) that will be used in api gateway"
}

variable "cloudflare_root_domain" {
  type        = string
  description = "the root domain that will be used in api gateway"
}

variable "website_cloudflare_domain" {
  type        = string
  description = "the domain (with the root and the subdomains) that will be used in the s3 website"
}

locals {
  mime_types = jsondecode(file("./0-mime.json"))
}
