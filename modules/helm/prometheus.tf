resource "helm_release" "prometheus" {
  create_namespace = true
  chart            = var.prometheus_chart
  name             = var.prometheus_name
  namespace        = var.namespace
  repository       = var.prometheus_repository
  version          = var.prometheus_version

  set {
    name  = "podSecurityPolicy.enabled"
    value = true
  }

  set {
    name  = "server.persistentVolume.enabled"
    value = false
  }

  set {
    name = "server\\.resources"
    value = yamlencode({
      limits = {
        cpu    = "200m"
        memory = "50Mi"
      }
      requests = {
        cpu    = "100m"
        memory = "30Mi"
      }
    })
  }
}
