module "alb" {
  source  = "terraform-aws-modules/alb/aws"
  version = ">= 10.5.0"
  
  name = "todo-app-alb"

  vpc_id  = var.vpc
  subnets = var.public_subnets

  security_groups = [module.alb_sg.security_group_id]

  enable_deletion_protection = false

#   target_groups = {
#     tg-frontend = {
#       protocol = "HTTP"
#       port     = 3000
#       health_check = {
#         enabled = true
#         path    = "/"
#       }
#     },
#     tg-backend = {
#       protocol = "HTTP"
#       port     = 8080
#       health_check = {
#         enabled = true
#         path    = "/"
#       }
#     }
#   }
}
