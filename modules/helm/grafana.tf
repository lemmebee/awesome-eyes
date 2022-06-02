resource "kubernetes_secret" "grafana" {
  metadata {
    name      = var.grafana_kubernetes_secret_name
    namespace = var.namespace
  }

  data = {
    admin-user     = var.grafana_admin_user
    admin-password = var.grafana_admin_password
  }

  depends_on = [
    helm_release.prometheus
  ]
}

resource "helm_release" "grafana" {
  chart      = var.grafana_chart
  name       = var.grafana_name
  repository = var.grafana_repository
  namespace  = var.namespace
  version    = var.grafana_version

  values = [
    templatefile("${path.module}/templates/grafana-values.yaml", {
      admin_existing_secret = kubernetes_secret.grafana.metadata[0].name
      admin_user_key        = var.grafana_admin_user_key
      admin_password_key    = var.grafana_admin_password_key
      prometheus_svc        = "${helm_release.prometheus.name}-server"
      replicas              = var.grafana_replicas
    })
  ]

  depends_on = [
    kubernetes_secret.grafana
  ]
}
