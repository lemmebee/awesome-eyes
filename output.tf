output "prometheus_release_namespace" {
  description = "Prometheus release namespace"
  value       = helm_release.prometheus.namespace
}

output "grafana_release_namespace" {
  description = "Grafana release namespace"
  value       = helm_release.grafana.namespace
}
