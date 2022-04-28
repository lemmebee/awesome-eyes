output "kubeconfig_filename" {
  description = "Kubernetes Cluster Kubeconfig Filename"
  value       = module.eks.kubeconfig_filename
}
