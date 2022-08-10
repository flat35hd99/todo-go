locals {
  lambda_function_name   = "server"
  lambda_binary_filename = "app"
}

resource "aws_lambda_function" "server" {
  function_name = local.lambda_function_name
  role          = aws_iam_role.for_lambda.arn

  runtime          = "go1.x"
  filename         = data.archive_file.archive_binary.output_path
  handler          = local.lambda_binary_filename
  source_code_hash = data.archive_file.archive_binary.output_base64sha256

  # TODO: This configuration can make terraform to deploy when source of lambda functions
  # are changed but cause of unknown reason, someitmes apply will fail. 
  lifecycle {
    replace_triggered_by = [
      data.archive_file.archive_binary
    ]
  }
}

resource "aws_lambda_function_url" "endpoint" {
  function_name      = aws_lambda_function.server.function_name
  authorization_type = "NONE"
  cors {
    allow_headers = ["*"]
    allow_methods = ["*"]
    allow_origins = ["https://${aws_cloudfront_distribution.web.domain_name}"]
    max_age       = 60
  }
}

resource "aws_iam_role" "for_lambda" {
  name               = "iam_for_lambda"
  assume_role_policy = <<EOF
  {
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_cloudwatch_log_group" "lambda_log" {
  name              = "/aws/lambda/${local.lambda_function_name}"
  retention_in_days = 14
}

resource "aws_iam_policy" "lambda_logging" {
  name        = "lambda_logging"
  path        = "/"
  description = "IAM policy for logging from server lambda"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Resource": "arn:aws:logs:*:*:*",
      "Effect": "Allow"
    }
  ]
}
EOF
}

resource "aws_iam_policy_attachment" "lambda_log" {
  name       = "lambda_log"
  roles      = [aws_iam_role.for_lambda.name]
  policy_arn = aws_iam_policy.lambda_logging.arn
}

// Archive binary file to upload
data "archive_file" "archive_binary" {
  source_file = "${path.module}/backend/cmd/lambda/${local.lambda_binary_filename}"
  output_path = "${path.module}/function_payload.zip"
  type        = "zip"
}
