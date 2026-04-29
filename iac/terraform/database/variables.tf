variable "backend_sg" {
  description = "Security group from backend ECS service"
}
variable "vpc" {
  description = "todo VPC"
}
variable "private_subnets" {
  description = "private subnets for RDS"
}
variable "my_ip" {}
