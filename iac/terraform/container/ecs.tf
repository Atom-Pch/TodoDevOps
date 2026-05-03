module "ecs" {
  source  = "terraform-aws-modules/ecs/aws"
  version = ">= 7.5.0"

  cluster_name               = "todo-app-cluster"
  cluster_capacity_providers = ["FARGATE"]

  services = {
    todo-frontend-task = {
      cpu           = 1024
      memory        = 2048
      desired_count = 2

      container_definitions = {
        frontend-container = {
          image     = "131912109503.dkr.ecr.us-east-2.amazonaws.com/todo-frontend-repo:latest"
          essential = true

          portMappings = [
            {
              name          = "todo-frontend-task"
              containerPort = 3000
              hostPort      = 3000
              protocol      = "tcp"
            }
          ]

          environment = [
            {
              name  = "VITE_API_URL",
              value = var.alb_dns
            }
          ]
        }
      }

      load_balancer = {
        service = {
          target_group_arn = var.alb_tg["tg-frontend"].arn
          container_name   = "frontend-container"
          container_port   = 3000
        }
      }

      security_group_ids = [module.frontend_sg.security_group_id]
      subnet_ids         = var.private_subnets
      assign_public_ip   = false

      task_exec_iam_role_policies = {
        env_policy   = var.todo_env_policy
        files_policy = var.todo_files_policy
      }
    }

    todo-backend-task = {
      cpu           = 1024
      memory        = 2048
      desired_count = 2

      container_definitions = {
        backend-container = {
          image     = "131912109503.dkr.ecr.us-east-2.amazonaws.com/todo-backend-repo:latest"
          essential = true

          portMappings = [
            {
              name          = "todo-backend-task"
              containerPort = 8080
              hostPort      = 8080
              protocol      = "tcp"
            }
          ]

          environment = [
            {
              name = "DB_HOST"
              value  = var.db_address
            },
            {
              name = "DB_USER"
              value  = "atom"
            },
            {
              name = "DB_NAME"
              value  = "todo_db"
            },
            {
              name ="S3_BUCKET_NAME"
              value = var.s3_files_name
            },
            {
              name = "AWS_REGION"
              value = "us-east-2"
            }
          ]

          secrets = [
            {
              name = "DB_PASS"
              valueFrom = "arn:aws:secretsmanager:us-east-2:131912109503:secret:rds!db-f7c461f3-645c-40ac-891a-590b5c39d73a-QUSUCP:password::"
            }
          ]

          environmentFiles = [
            {
              value = "${var.s3_env_arn}/.env"
              type  = "s3"
            }
          ]
        }
      }

      load_balancer = {
        service = {
          target_group_arn = var.alb_tg["tg-backend"].arn
          container_name   = "backend-container"
          container_port   = 8080
        }
      }

      security_group_ids = [module.backend_sg.security_group_id]
      subnet_ids         = var.private_subnets
      assign_public_ip   = false

      task_exec_iam_role_policies = {
        env_policy   = var.todo_env_policy
        files_policy = var.todo_files_policy
      }
    }
  }

  task_exec_secret_arns = []

  create_security_group     = false
  create_task_exec_iam_role = true
  create_task_exec_policy   = true
}
