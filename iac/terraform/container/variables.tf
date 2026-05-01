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