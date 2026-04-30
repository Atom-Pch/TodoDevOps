module "frontend_sg" {
  source  = "terraform-aws-modules/security-group/aws"
  version = ">= 5.3.1"

  name        = "todo-frontend-service-sg"
  description = "Allow todo frontend service to receive connections from ALB"
  vpc_id      = var.vpc

  ingress_with_source_security_group_id = [
    {
      from_port   = 3000
      to_port     = 3000
      protocol    = "tcp"
      description = "Frontend port"
      source_security_group_id = var.alb_sg
    }
  ]

  egress_rules       = ["all-tcp"]
  egress_cidr_blocks = ["0.0.0.0/0"]
}

module "backend_sg" {
  source  = "terraform-aws-modules/security-group/aws"
  version = ">= 5.3.1"

  name        = "todo-backend-service-sg"
  description = "Allow todo backend service to receive connections from ALB"
  vpc_id      = var.vpc

  ingress_with_source_security_group_id = [
    {
      from_port   = 8080
      to_port     = 8080
      protocol    = "tcp"
      description = "Backend port"
      source_security_group_id = var.alb_sg
    }
  ]

  egress_rules       = ["all-tcp"]
  egress_cidr_blocks = ["0.0.0.0/0"]
}
