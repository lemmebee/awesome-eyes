resource "aws_security_group" "worker_mgmt" {
  name_prefix = "worker_management"
  vpc_id      = module.vpc.vpc_id

  ingress {
    from_port = 22
    to_port   = 22
    protocol  = "tcp"
  }
}
