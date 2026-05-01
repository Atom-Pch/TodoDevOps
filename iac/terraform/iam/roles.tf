# data "aws_iam_policy_document" "ecs_assume_role_policy" {
#   statement {
#     effect  = "Allow"
#     actions = ["sts:AssumeRole"]
#     principals {
#       type        = "Service"
#       identifiers = ["ecs-tasks.amazonaws.com"]
#     }
#   }
# }

# # resource "aws_iam_role" "todo_env_role" {
# #   name        = "ReadOnly-TodoApp-Env"
# #   description = "Allow services to get env file from todo S3"
# #   assume_role_policy = data.aws_iam_policy_document.ecs_assume_role_policy.json
# # }

# # resource "aws_iam_role" "todo_S3files_role" {
# #   name = "GetPutDel-todoApp-files"
# #   description = "allow services to get/put/delete files from todo S3"
# #   assume_role_policy = data.aws_iam_policy_document.ecs_assume_role_policy.json
# # }

# # resource "aws_iam_role_policy_attachment" "ecs_env" {
# #   role = aws_iam_role.todo_env_role.name
# #   policy_arn = aws_iam_policy.S3_todo_env_read.arn
# # }

# # resource "aws_iam_role_policy_attachment" "ecs_files" {
# #   role = aws_iam_role.todo_S3files_role.name
# #   policy_arn = aws_iam_policy.S3_todo_files_getPutDel.arn
# # }