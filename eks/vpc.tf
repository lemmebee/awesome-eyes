data "aws_availability_zones" "available" {}

locals {
  cluster_name = "awesome-eyes-eks-${random_string.suffix.result}"
}

resource "random_string" "suffix" {
  length  = 8
  special = false
}

module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "3.2.0"

  name           = "awesome-eyes-vpc"
  cidr           = "10.0.0.0/16"
  azs            = data.aws_availability_zones.available.names
  public_subnets = ["10.0.4.0/24", "10.0.5.0/24"]

  tags = {
    "kubernetes.io/cluster/${local.cluster_name}" = "shared"
  }

  public_subnet_tags = {
    "kubernetes.io/cluster/${local.cluster_name}" = "shared"
    "kubernetes.io/role/elb"                      = "1"
  }
}
