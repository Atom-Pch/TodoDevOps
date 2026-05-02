resource "aws_iam_policy" "S3_todo_env" {
  name        = "S3TodoEnv"
  description = "Allow services to get env file from S3 for todo app"

  policy = jsonencode({
    Version : "2012-10-17",
    Statement : [
      {
        Action : [
          "s3:GetBucketLocation",
          "s3:GetObject",
          "secretsmanager:GetSecretValue"
        ],
        Effect : "Allow",
        Resource : [
          "arn:aws:s3:::todo-env-131912109503-us-east-2-an",
          "arn:aws:s3:::todo-env-131912109503-us-east-2-an/*",
          "arn:aws:secretsmanager:us-east-2:131912109503:secret:rds!db-f7c461f3-645c-40ac-891a-590b5c39d73a-QUSUCP"
        ]
      }
    ]
  })
}

resource "aws_iam_policy" "S3_todo_files_getPutDel" {
  name        = "S3TodoFilesGETPUTDEL"
  description = "Allow services to get/put/del files from S3 for todo app"

  policy = jsonencode({
    Version : "2012-10-17",
    Statement : [
      {
        Action : [
          "s3:GetObject",
          "s3:PutObject",
          "s3:DeleteObject"
        ],
        Effect : "Allow",
        Resource : [
          "arn:aws:s3:::todo-files-131912109503-us-east-2-an",
          "arn:aws:s3:::todo-files-131912109503-us-east-2-an/*"
        ]
      }
    ]
  })
}
