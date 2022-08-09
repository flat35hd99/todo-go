output "api_endpoint_url" {
  value = aws_lambda_function_url.endpoint.function_url
}

output "web_url" {
  value = aws_cloudfront_distribution.web.domain_name
}
