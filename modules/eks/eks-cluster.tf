module "eks" {
  source                 = "terraform-aws-modules/eks/aws"
  version                = "17.24.0"
  cluster_name           = local.cluster_name
  cluster_tags           = var.cluster_tags
  cluster_version        = "1.20"
  subnets                = module.vpc.public_subnets
  write_kubeconfig       = true
  kubeconfig_output_path = "/home/runner/.kube/config"
  kubeconfig_name        = local.cluster_name

  vpc_id = module.vpc.vpc_id

  workers_group_defaults = {
    root_volume_type = "gp2"
  }

  worker_groups = [
    {
      name                          = "worker-group-1"
      instance_type                 = "t2.small"
      additional_security_group_ids = [aws_security_group.worker_mgmt.id]
      asg_desired_capacity          = 1
    },
  ]
}
