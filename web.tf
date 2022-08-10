locals {
  // Thx https://engineering.statefarm.com/blog/terraform-s3-upload-with-mime/
  mime_types = {
    ".html" : "text/html"
    ".js" : "application/javascript; charset=utf-8"
    ".css" : "text/css"
  }
}

resource "aws_s3_bucket" "web" {
  bucket_prefix = "flat35hd99oiwie"
}

resource "aws_s3_bucket_website_configuration" "web" {
  bucket = aws_s3_bucket.web.bucket
  index_document {
    suffix = "index.html"
  }
  error_document {
    key = "index.html"
  }
}

resource "aws_s3_object" "contents" {
  for_each     = fileset("${path.module}/front/dist/", "**")
  bucket       = aws_s3_bucket.web.bucket
  key          = each.value
  source       = "${path.module}/front/dist/${each.value}"
  content_type = lookup(local.mime_types, regex("\\.[^.]+$", each.value), null)
  etag         = filemd5("${path.module}/front/dist/${each.value}")
  depends_on = [
    null_resource.build_frontend
  ]
}

resource "null_resource" "build_frontend" {
  triggers = {
    # Fix me.
    # I should check all files without front/node_modules and front/dist/
    "source_hash" = join("", [for f in fileset(path.module, "/front/src/**/*.tsx") : filebase64sha256(f)])
    "dist_hash"   = fileexists("${path.module}/front/dist") ? join("", [for f in fileset(path.module, "/front/dist/**") : filebase64sha256(f)]) : ""
  }
  provisioner "local-exec" {
    command     = "yarn && yarn build"
    working_dir = "${path.module}/front"
    environment = {
      VITE_API_ENDPOINT_URL = aws_lambda_function_url.endpoint.function_url
    }
  }
}

resource "aws_s3_bucket_acl" "web" {
  bucket = aws_s3_bucket.web.bucket
  acl    = "private"
}

resource "aws_s3_bucket_policy" "web" {
  bucket = aws_s3_bucket.web.id
  policy = data.aws_iam_policy_document.web_bucket.json
}

data "aws_iam_policy_document" "web_bucket" {
  statement {
    sid    = "Allow CloudFront"
    effect = "Allow"
    principals {
      type        = "AWS"
      identifiers = [aws_cloudfront_origin_access_identity.web.iam_arn]
    }
    actions = [
      "s3:GetObject"
    ]

    resources = [
      "${aws_s3_bucket.web.arn}/*"
    ]
  }
}

resource "aws_cloudfront_distribution" "web" {
  origin {
    domain_name = aws_s3_bucket.web.bucket_regional_domain_name
    origin_id   = aws_s3_bucket.web.id
    s3_origin_config {
      origin_access_identity = aws_cloudfront_origin_access_identity.web.cloudfront_access_identity_path
    }
  }

  enabled = true

  default_root_object = "index.html"

  custom_error_response {
    error_code         = 403
    response_code      = 200
    response_page_path = "/"
  }

  default_cache_behavior {
    allowed_methods  = ["GET", "HEAD"]
    cached_methods   = ["GET", "HEAD"]
    target_origin_id = aws_s3_bucket.web.id

    forwarded_values {
      query_string = false
      cookies {
        forward = "none"
      }
    }

    viewer_protocol_policy = "redirect-to-https"
    min_ttl                = 0
    default_ttl            = 20
    max_ttl                = 30
  }

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  viewer_certificate {
    cloudfront_default_certificate = true
  }
}

resource "aws_cloudfront_origin_access_identity" "web" {}
