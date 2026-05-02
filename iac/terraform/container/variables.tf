variable "tag_policy" {
  description = "ECR tags policy"
}
variable "alb_sg" {
  description = "ALB security group for containers"
}
variable "vpc" {
  description = "todo vpc"
}
variable "private_subnets" {
  description = "private subnets for RDS"
}
variable "alb_tg" {
  description = "ALB target group arn"
}
variable "todo_env_policy" {
  description = "IAM policy for S3 access env"
}
variable "todo_files_policy" {
  description = "IAM policy for S3 access files"
}
variable "alb_dns" {
  description = "DNS name from ALB"
}
variable "s3_files_name" {
  description = "name of todo files S3 bucket"
}
variable "s3_env_arn" {
  description = "ARN of todo env S3 bucket"
}
variable "db_address" {
  description = "address of todo RDS"
}