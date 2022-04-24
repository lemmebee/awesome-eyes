variable "region" {
  type    = string
  default = "eu-west-3"
}

variable "cluster_tags" {
  type        = map(string)
  description = "Kubernetes Cluster Tags"
  default = {
    "Name" = "awesome-eyes-eks-cluster"
  }
}
