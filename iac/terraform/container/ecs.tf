module "ecs" {
  source  = "terraform-aws-modules/ecs/aws"
  version = ">= 7.5.0"

  cluster_name               = "todo-app-cluster"
  cluster_capacity_providers = ["FARGATE"]

  services = {
    todo-frontend-task = {
      cpu           = 512
      memory        = 1024
      desired_count = 1

      container_definitions = {
        backend-container = {
          image     = "131912109503.dkr.ecr.us-east-2.amazonaws.com/todo-frontend-repo:latest"
          essential = true
        }
      }

      portMappings = [
        {
          name          = "todo-frontend-task"
          containerPort = 3000
          hostPort      = 3000
          protocol      = "tcp"
        }
      ]

      security_group_ids = [module.frontend_sg.security_group_id]
      subnet_ids         = var.private_subnets
      assign_public_ip   = false
    }

    todo-backend-task = {
      cpu           = 512
      memory        = 1024
      desired_count = 1

      container_definitions = {
        backend-container = {
          image     = "131912109503.dkr.ecr.us-east-2.amazonaws.com/todo-backend-repo:latest"
          essential = true
        }
      }

      portMappings = [
        {
          name          = "todo-backend-task"
          containerPort = 8080
          hostPort      = 8080
          protocol      = "tcp"
        }
      ]

      security_group_ids = [module.backend_sg.security_group_id]
      subnet_ids         = var.private_subnets
      assign_public_ip   = false
    }
  }

  create_security_group     = false
  create_task_exec_iam_role = true
  create_task_exec_policy   = true
}
