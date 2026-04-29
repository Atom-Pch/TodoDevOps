module "ecs" {
  source  = "terraform-aws-modules/ecs/aws"
  version = ">= 7.5.0"

  cluster_name = "todo-app-cluster"
}
