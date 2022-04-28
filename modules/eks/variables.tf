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
