output "cluster_name" {
  description = "eks cluster name"
  value       = aws_eks_cluster.this.id
}


output "region" {
  description = "aws region"
  value       = var.region
}
