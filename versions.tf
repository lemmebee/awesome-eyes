terraform {

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 3.74.0"
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
      version = "2.9.0"
    }

    helm = {
      source  = "hashicorp/helm"
      version = "2.4.1"
    }

  }

  required_version = ">= 1.1.2"
}
