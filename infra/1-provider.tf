terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0.0"
    }

    random = {
      source  = "hashicorp/random"
      version = "~> 3.5.1"
    }

    archive = {
      source  = "hashicorp/archive"
      version = "~> 2.3.0"
    }

    time = {
      source  = "hashicorp/time"
      version = "~> 0.9.1"
    }
  }

  required_version = "~> 1.0"
}

provider "aws" {
  region = "sa-east-1"

  access_key = var.aws_access_key
  secret_key = var.aws_secret_key

  default_tags {
    tags = {
      Project     = "Games 4 life"
      ManagedBy   = "Terraform"
      CreatedAt   = "May 2023"
      Environment = var.environment_type
    }
  }
}
