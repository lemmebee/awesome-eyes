locals {
  cluster_name = "awesome-eyes-eks-${random_string.suffix.result}"
}

resource "random_string" "suffix" {
  length  = 8
  special = false
}

resource "aws_eks_cluster" "this" {
  name     = local.cluster_name
  version  = var.kubernetes_version
  role_arn = aws_iam_role.eks.arn

  vpc_config {
    subnet_ids = aws_subnet.this.*.id
  }

  depends_on = [
    aws_iam_role_policy_attachment.eks_AmazonEKSClusterPolicy,
  ]
}

resource "aws_eks_node_group" "this" {
  cluster_name    = local.cluster_name
  node_group_name = "ng-${local.cluster_name}"
  node_role_arn   = aws_iam_role.eks_node.arn
  subnet_ids      = aws_subnet.this.*.id
  instance_types  = ["t2.large"]

  # https://stackoverflow.com/questions/72161772/k8s-deployment-is-not-scaling-on-eks-cluster-too-many-pods
  scaling_config {
    desired_size = 2
    max_size     = 3
    min_size     = 2
  }

  depends_on = [
    aws_iam_role_policy_attachment.eks_AmazonEKSWorkerNodePolicy,
    aws_iam_role_policy_attachment.eks_AmazonEKS_CNI_Policy,
    aws_iam_role_policy_attachment.eks_AmazonEC2ContainerRegistryReadOnly,
    # aws_eks_cluster.this
  ]
}

