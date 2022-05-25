variable "region" {
  type = string
}

variable "cluster_tags" {
  type        = map(string)
  description = "Kubernetes Cluster Tags"
  default = {
    "Name" = "awesome-eyes-eks-cluster"
  }
}

variable "kubernetes_version" {
  type        = string
  default     = "1.22"
  description = "kubernetes version for eks cluster"
}

variable "tag" {
  type        = string
  default     = "awesome-eyes-eks"
  description = "tag for aws resources"
}
