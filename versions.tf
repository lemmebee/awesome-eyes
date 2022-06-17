terraform {
    backend "s3" {
    bucket         = "awesome-eyes-terrafrom-state"
    key            = "terraform.tfstate"
    region         = "eu-west-3"
    dynamodb_table = "awesome-eyes-locks"
  }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "4.13.0"
    }

    random = {
      source  = "hashicorp/random"
      version = "3.1.2"
    }

    local = {
      source  = "hashicorp/local"
      version = "2.1.0"
    }

    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "2.11.0"
    }

    helm = {
      source  = "hashicorp/helm"
      version = "2.5.1"
    }

  }

  required_version = ">= 1.1.2"
}
