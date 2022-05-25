output "cluster_endpoint" {
  description = "eks cluster endpoint"
  value       = data.aws_eks_cluster.cluster.endpoint
}

output "cluster_name" {
  description = "eks cluster name"
  value       = module.eks.cluster_name
}

output "region" {
  description = "aws region"
  value       = var.region
}

output "prometheus_release_namespace" {
  description = "Prometheus release namespace"
  value       = module.helm.prometheus_release_namespace
}

output "grafana_release_namespace" {
  description = "Grafana release namespace"
  value       = module.helm.grafana_release_namespace
}
