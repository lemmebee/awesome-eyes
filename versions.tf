terraform {

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 4.12.0"
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
      version = "2.0.1"
    }

    helm = {
      source  = "hashicorp/helm"
      version = "2.0.1"
    }

  }

  required_version = ">= 1.1.2"
}
