data "archive_file" "lambda_code" {
  type        = "zip"
  output_path = "../dist/ddns.zip"
  source_dir = "./bin"
}

resource "aws_lambda_function" "lambda" {
  function_name = var.name
  description = "A lambda to handle dynamic DNS requests to update Route53"

  handler = "main"
  runtime = "go1.x"
  filename = data.archive_file.lambda_code.output_path
  source_code_hash = filebase64sha256(data.archive_file.lambda_code.output_path)
  memory_size = 128
  timeout = 20
  environment {
    variables = {
      username = var.username
      password = var.password
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

resource "aws_iam_policy" "ddns-service-update" {
  name = "${var.name}-service-dns-update"
  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": "route53:ChangeResourceRecordSets",
            "Resource": "arn:aws:route53:::hostedzone/Z1Z4NARII70GPM"
        },
        {
            "Effect": "Allow",
            "Action": "route53:ChangeResourceRecordSets",
            "Resource": "arn:aws:route53:::hostedzone/Z02568913AU6JLN00TL3A"
        },
        {
            "Effect": "Allow",
            "Action": "route53:ListHostedZones",
            "Resource": "*"
        }
    ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "ddns_update_attachment" {
  policy_arn = aws_iam_policy.ddns-service-update.arn
  role = aws_iam_role.lambda_execution_role.name
}