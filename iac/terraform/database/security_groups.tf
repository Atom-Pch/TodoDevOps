module "rds_sg" {
  source  = "terraform-aws-modules/security-group/aws"
  version = ">= 5.3.1"

  name        = "todo-rds-sg"
  description = "Allow todo RDS service to receive connections from backend and local"
  vpc_id      = var.vpc

  ingress_rules       = ["postgresql-tcp"]
  ingress_cidr_blocks = [var.my_ip]
  ingress_with_source_security_group_id = [
    {
      rule = "postgresql-tcp"
      description = "ALB port",
      source_security_group_id = var.backend_sg
    }
  ]

  egress_rules       = ["all-tcp"]
  egress_cidr_blocks = ["0.0.0.0/0"]
}
