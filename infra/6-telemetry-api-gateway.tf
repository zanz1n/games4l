resource "aws_apigatewayv2_integration" "lambda_telemetry" {
  api_id = aws_apigatewayv2_api.main.id

  integration_uri    = aws_lambda_function.telemetry.invoke_arn
  integration_type   = "AWS_PROXY"
  integration_method = "POST"
}

resource "aws_apigatewayv2_route" "post_telemetry" {
  api_id = aws_apigatewayv2_api.main.id

  route_key = "POST /telemetry"

  target = "integrations/${aws_apigatewayv2_integration.lambda_telemetry.id}"
}

resource "aws_apigatewayv2_route" "get_telemetry" {
  api_id = aws_apigatewayv2_api.main.id

  route_key = "GET /telemetry"

  target = "integrations/${aws_apigatewayv2_integration.lambda_telemetry.id}"
}
