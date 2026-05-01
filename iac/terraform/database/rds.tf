resource "aws_db_subnet_group" "this" {
  name        = "todo-db-subnet-group"
  description = "subnet group for todo DB"

  subnet_ids = var.private_subnets
}

module "rds" {
  source  = "terraform-aws-modules/rds/aws"
  version = ">= 7.2.0"

  vpc_security_group_ids = [module.rds_sg.security_group_id]

  identifier        = "todo-db"
  allocated_storage = 20
  storage_type      = "gp2"
  engine            = "postgres"
  engine_version    = "18.2"
  instance_class    = "db.t3.micro"
  username          = "atom"
  db_name           = "todo_db"

  family               = "postgres18"
  major_engine_version = "18.0"

  publicly_accessible = false

  db_subnet_group_name = aws_db_subnet_group.this.name

  skip_final_snapshot                                    = true
  manage_master_user_password                            = true
  manage_master_user_password_rotation                   = true
  master_user_password_rotate_immediately                = false
  master_user_password_rotation_automatically_after_days = 30
}
