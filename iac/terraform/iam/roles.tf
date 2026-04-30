data "aws_iam_policy_document" "ecs_assume_role_policy" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["ecs.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "ecs_role" {
  name        = "ReadOnly-TodoApp-Env"
  description = "Allow services to get env file from S3"
  assume_role_policy = data.aws_iam_policy_document.ecs_assume_role_policy.json
}

resource "aws_iam_role_policy_attachment" "attach" {
  role = aws_iam_role.ecs_role.name
  policy_arn = aws_iam_policy.S3_todo_env_read.arn
}