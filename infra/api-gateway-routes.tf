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