terraform {
  required_version = ">= 1.14.8"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.0"
    }
  }
}

provider "aws" {
  region  = var.aws_region
  profile = var.aws_profile

  default_tags {
    tags = {
      Environment = "Practice"
      ManagedBy   = "Terraform"
      Project     = "TodoApp"
      # awsApplication = aws_servicecatalogappregistry_application.todo_app.application_tag
    }
  }
}

resource "aws_servicecatalogappregistry_application" "todo_app" {
  name = "todo-app-for-devops"
  description = "Todo web application managed via Terraform"
}

module "networks" {
  source = "./networks"

  aws_region = var.aws_region
}

module "database" {
  source = "./database"

  backend_sg  = module.container.backend_sg
  vpc         = module.networks.vpc
  my_ip       = var.my_ip
  private_subnets = module.networks.private_subnets
}

module "load_balancer" {
  source = "./load_balancer"

  vpc = module.networks.vpc
  public_subnets = module.networks.pubic_subnets
}

module "container" {
  source = "./container"

  tag_policy = "IMMUTABLE_WITH_EXCLUSION"
  alb_sg     = module.load_balancer.alb_sg
  vpc        = module.networks.vpc
}

module "storage" {
  source = "./storage"

  aws_region = var.aws_region
}
