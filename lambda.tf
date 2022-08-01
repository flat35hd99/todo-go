locals {
  lambda_function_name = "server"
}

resource "aws_lambda_function" "server" {
  function_name = local.lambda_function_name
  role          = aws_iam_role.for_lambda

  runtime = "go1.x"

  depends_on = [
    aws_iam_role_policy_attachment.lambda_log,
    aws_cloudwatch_log_group.lambda_log
  ]
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
  name              = "/aws/lambda/${locals.lambda_function_name}"
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
  roles      = aws_iam_role.for_lambda.name
  policy_arn = aws_iam_policy.lambda_logging.arn
}
