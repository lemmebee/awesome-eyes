output "cluster_endpoint" {
  description = "eks cluster endpoint"
  value       = data.aws_eks_cluster.cluster.endpoint
}

output "cluster_name" {
  description = "eks cluster name"
  value       = module.eks.cluster_name
}

output "node_group_name" {
  description = "Eks node group name"
  value       = module.eks.node_group_name
}

output "prometheus_release_namespace" {
  description = "Prometheus release namespace"
  value       = module.helm.prometheus_release_namespace
}

output "grafana_release_namespace" {
  description = "Grafana release namespace"
  value       = module.helm.grafana_release_namespace
}
