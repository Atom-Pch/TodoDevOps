output "todo_env_policy" {
  value = aws_iam_policy.S3_todo_env_read.arn
}
output "todo_files_policy" {
  value = aws_iam_policy.S3_todo_files_getPutDel.arn
}