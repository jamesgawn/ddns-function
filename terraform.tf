variable "profile" {
  type = string
  default = "default"
}

variable "region" {
  type = string
  default = "eu-west-2"
}

variable "domain" {
  type = string
}

variable "username" {
  type = string
}

variable "password" {
  type = string
}

variable "cert_domain" {
  type = string
}

variable "zone" {
  type = string
}

variable "name" {
  type = string
}

variable "notification_sns_queue_name" {
  type = string
  description = "The name of the SNS queue to send ok/error alarms if the lambda stops working."
}

provider "aws" {
  region = var.region
  profile = var.profile
}

terraform {
  backend "s3" {
    bucket = "ana-terraform-state-prod"
    key = "ddns-lambda/terraform.tfstate"
    region = "eu-west-2"
  }
}

data "archive_file" "lambda_code" {
  type        = "zip"
  output_path = "${path.root}/dist/ddns.zip"
  source_dir = "./bin"
}

resource "aws_lambda_function" "lambda" {
  function_name = var.name
  description = "A lambda to handle dynamic DNS requests to update Route53"

  handler = "main"
  runtime = "go1.x"
  filename = data.archive_file.lambda_code.output_path
  source_code_hash = data.archive_file.lambda_code.output_sha
  memory_size = 128
  timeout = 20
  environment {
    variables = {
      username = var.username
      password = var.password
      version = "0.0.0"
    }
  }

  role = aws_iam_role.lambda_execution_role.arn
}

resource "aws_iam_role" "lambda_execution_role" {
  name = "${var.name}-lambda-execution-role"

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

resource "aws_iam_role_policy_attachment" "lambda_execution_basis_role_attachment" {
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
  role = aws_iam_role.lambda_execution_role.name
}

// API Gateway Configuration for Endpoints

resource "aws_apigatewayv2_api" "api" {
  name          = "${var.name}-api"
  protocol_type = "HTTP"
  disable_execute_api_endpoint = true
}

resource "aws_cloudwatch_log_group" "api" {
  name = "${var.name}-api"
  retention_in_days = 14
}

resource "aws_apigatewayv2_stage" "api" {
  api_id = aws_apigatewayv2_api.api.id
  name   = "prod"
  auto_deploy = true
  access_log_settings {
    destination_arn = aws_cloudwatch_log_group.api.arn
    format = <<EOF
    { "requestId":"$context.requestId", "ip": "$context.identity.sourceIp", "requestTime":"$context.requestTime", "httpMethod":"$context.httpMethod","routeKey":"$context.routeKey", "status":"$context.status","protocol":"$context.protocol", "responseLength":$context.responseLength, "integrationStatus": $context.integrationStatus, "integrationErrorMessage": "$context.integrationErrorMessage", "integration": { "error": "$context.integration.error", "integrationstatus": $context.integration.integrationStatus, "latency": $context.integration.latency, "requestId": "$context.integration.requestId", "status": $context.integration.status } }
EOF
  }
}

data "aws_acm_certificate" "api_cert" {
  domain   = var.cert_domain
  statuses = ["ISSUED"]
}

resource "aws_apigatewayv2_domain_name" "api" {
  domain_name = var.domain

  domain_name_configuration {
    certificate_arn = data.aws_acm_certificate.api_cert.arn
    endpoint_type   = "REGIONAL"
    security_policy = "TLS_1_2"
  }
}

resource "aws_apigatewayv2_api_mapping" "example" {
  api_id      = aws_apigatewayv2_api.api.id
  domain_name = aws_apigatewayv2_domain_name.api.id
  stage       = aws_apigatewayv2_stage.api.id
}

data "aws_route53_zone" "api" {
  name = var.zone
  private_zone = false
}

resource "aws_route53_record" "api" {
  name    = aws_apigatewayv2_domain_name.api.domain_name
  type    = "A"
  zone_id = data.aws_route53_zone.api.zone_id

  alias {
    name                   = aws_apigatewayv2_domain_name.api.domain_name_configuration[0].target_domain_name
    zone_id                = aws_apigatewayv2_domain_name.api.domain_name_configuration[0].hosted_zone_id
    evaluate_target_health = false
  }
}

resource "aws_apigatewayv2_integration" "integration" {
  api_id                  = aws_apigatewayv2_api.api.id
  integration_type        = "AWS_PROXY"
  integration_uri         = aws_lambda_function.lambda.arn
  payload_format_version  = "2.0"
}

resource "aws_apigatewayv2_route" "route" {
  api_id = aws_apigatewayv2_api.api.id
  route_key = "GET /"
  target = "integrations/${aws_apigatewayv2_integration.integration.id}"
}

resource "aws_apigatewayv2_route" "route-update" {
  api_id = aws_apigatewayv2_api.api.id
  route_key = "GET /nic/update"
  target = "integrations/${aws_apigatewayv2_integration.integration.id}"
}

resource "aws_lambda_permission" "permission" {
  action = "lambda:InvokeFunction"
  function_name = aws_lambda_function.lambda.function_name
  principal = "apigateway.amazonaws.com"
  source_arn = "${aws_apigatewayv2_api.api.execution_arn}/*/*/*"
}

resource "aws_iam_group" "ddns-developer" {
  name = "ddns-developers"
}

resource "aws_iam_group_policy" "deployer-policy"  {
  name = "ddns-deployer-policy"
  group = aws_iam_group.ddns-developer.name
  policy = jsonencode({
    "Version": "2012-10-17",
    "Statement": [
      {
        Sid: "VisualEditor0",
        Effect: "Allow",
        Action: "lambda:UpdateFunctionCode",
        Resource: "arn:aws:lambda:eu-west-2:979779020614:function:ddns"
      }
    ]
  })
}