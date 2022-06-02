output "cluster_name" {
  description = "eks cluster name"
  value       = aws_eks_cluster.this.id
}


output "region" {
  description = "aws region"
  value       = var.region
}

output "node_group_name" {
  description = "Eks node group name"
  value       = aws_eks_node_group.this.node_group_name
}
