module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = ">= 6.6.1"

  name = "todo-vpc"
  cidr = "10.0.0.0/20"

  azs = ["us-east-2a", "us-east-2b"]

  public_subnets  = ["10.0.0.0/24", "10.0.1.0/24"]
  private_subnets = ["10.0.8.0/24", "10.0.9.0/24"]

  public_subnet_names  = ["todo-subnet-public1-us-east-2a", "todo-subnet-public2-us-east-2b"]
  private_subnet_names = ["todo-subnet-private1-us-east-2a", "todo-subnet-private2-us-east-2b"]

  enable_nat_gateway = false

  create_igw = true
  igw_tags = {
    name = "todo-igw"
  }
  public_route_table_tags = {
    name = "todo-rtb-public"
  }
  private_route_table_tags = {
    name = "todo-rtb-private"
  }

  enable_flow_log                      = true
  create_flow_log_cloudwatch_log_group = true
  create_flow_log_cloudwatch_iam_role  = true
}

resource "aws_vpc_endpoint" "s3" {
  vpc_id            = module.vpc.vpc_id
  service_name      = "com.amazonaws.${var.aws_region}.s3"
  vpc_endpoint_type = "Gateway"

  route_table_ids = module.vpc.private_route_table_ids

  tags = {
    Name = "todo-vpce-s3"
  }
}
