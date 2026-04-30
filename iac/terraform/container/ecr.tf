resource "aws_ecr_registry_scanning_configuration" "scan" {
  scan_type = "BASIC"
  rule {
    scan_frequency = "SCAN_ON_PUSH"
    repository_filter {
      filter      = "*"
      filter_type = "WILDCARD"
    }
  }
}

module "frontend-repo" {
  source  = "terraform-aws-modules/ecr/aws"
  version = ">= 3.2.0"

  repository_name = "todo-frontend-repo"
  repository_type = "private"

  repository_image_tag_mutability = var.tag_policy
  repository_image_tag_mutability_exclusion_filter = [
    {
      filter      = "latest"
      filter_type = "WILDCARD"
    }
  ]

  repository_lifecycle_policy = jsonencode({
    rules = [
      {
        rulePriority = 1,
        description  = "Keep last 10 images",
        selection = {
          tagStatus     = "tagged",
          tagPrefixList = ["v"],
          countType     = "imageCountMoreThan",
          countNumber   = 10
        },
        action = {
          type = "expire"
        }
      }
    ]
  })
}

module "backend-repo" {
  source  = "terraform-aws-modules/ecr/aws"
  version = ">= 3.2.0"

  repository_name = "todo-backend-repo"
  repository_type = "private"

  repository_image_tag_mutability = var.tag_policy
  repository_image_tag_mutability_exclusion_filter = [
    {
      filter      = "latest"
      filter_type = "WILDCARD"
    }
  ]

  repository_lifecycle_policy = jsonencode({
    rules = [
      {
        rulePriority = 1,
        description  = "Keep last 10 images",
        selection = {
          tagStatus     = "tagged",
          tagPrefixList = ["v"],
          countType     = "imageCountMoreThan",
          countNumber   = 10
        },
        action = {
          type = "expire"
        }
      }
    ]
  })
}
