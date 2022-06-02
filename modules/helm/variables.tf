variable "namespace" {
  type    = string
  default = "awesomeeyes"
}

variable "grafana_chart" {
  type    = string
  default = "grafana"
}

variable "grafana_name" {
  type    = string
  default = "grafana"
}

variable "grafana_repository" {
  type    = string
  default = "https://grafana.github.io/helm-charts"
}

variable "grafana_version" {
  type    = string
  default = "6.24.1"
}

variable "grafana_admin_user" {
  type    = string
  default = "admin"
}

variable "grafana_admin_password" {
  type = string
}

variable "grafana_admin_user_key" {
  type    = string
  default = "admin-user"
}

variable "grafana_admin_password_key" {
  type    = string
  default = "admin-password"
}

variable "grafana_replicas" {
  type    = number
  default = 1
}

variable "grafana_kubernetes_secret_name" {
  type    = string
  default = "grafana"
}

variable "prometheus_chart" {
  type    = string
  default = "prometheus"
}

variable "prometheus_name" {
  type    = string
  default = "prometheus"
}

variable "prometheus_repository" {
  type    = string
  default = "https://prometheus-community.github.io/helm-charts"
}

variable "prometheus_version" {
  type    = string
  default = "15.5.3"
}
