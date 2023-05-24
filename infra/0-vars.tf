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

variable "environment_type" {
  type        = string
  description = "dev or prod"
  default     = "dev"
}
