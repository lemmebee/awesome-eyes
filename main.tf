data "aws_eks_cluster" "cluster" {
  name = module.eks.cluster_name
}

data "aws_eks_cluster_auth" "cluster" {
  name = module.eks.cluster_name
}

module "eks" {
  source = "./modules/eks"
  region = var.region
}

module "helm" {
  source                 = "./modules/helm"
  grafana_admin_password = var.grafana_admin_password
  depends_on = [
    module.eks
  ]
}
