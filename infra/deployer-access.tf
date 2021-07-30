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