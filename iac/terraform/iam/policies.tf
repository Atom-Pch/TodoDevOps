resource "aws_iam_policy" "S3_todo_env_read" {
  name = "S3TodoEnvReadOnly"
  description = "Allow services to get env file from S3 for todo app"

  policy = jsonencode({
    Version : "2012-10-17",
    Statement : [
      {
        Action : [
          "s3:GetObject"
        ],
        Effect : "Allow",
        Resource : "arn:aws:s3:::todo-env-131912109503-us-east-2-an"
      }
    ]
  })
}
