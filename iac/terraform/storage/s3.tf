data "aws_caller_identity" "current" {}

module "todo-bucket" {
  source  = "terraform-aws-modules/s3-bucket/aws"
  version = ">= 5.12.0"

  bucket           = format("todo-files-%s-%s-an", data.aws_caller_identity.current.account_id, var.aws_region)
  bucket_namespace = "account-regional"

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true

  cors_rule = [
    {
      allowed_headers = ["*"],
      allowed_methods = ["GET", "PUT", "DELETE"],
      allowed_origins = ["*"],
      max_age_seconds = 3600
    }
  ]
}

module "env-bucket" {
  source  = "terraform-aws-modules/s3-bucket/aws"
  version = ">= 5.12.0"

  bucket           = format("todo-env-%s-%s-an", data.aws_caller_identity.current.account_id, var.aws_region)
  bucket_namespace = "account-regional"

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}