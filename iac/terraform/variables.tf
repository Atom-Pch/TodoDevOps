variable "aws_region" {
  description = "AWS region"
  default     = "us-east-2"
}

variable "aws_profile" {
  description = "CLI profile"
  sensitive   = true
}

variable "db_password" {
  description = "Password for RDS"
  sensitive   = true
}

variable "my_ip" {
  description = "local IP for development"
  sensitive = true
}