resource "aws_apigatewayv2_integration" "lambda_auth" {
  api_id = aws_apigatewayv2_api.main.id

  integration_uri    = aws_lambda_function.auth.invoke_arn
  integration_type   = "AWS_PROXY"
  integration_method = "POST"
}

resource "aws_apigatewayv2_route" "post_sign_in" {
  api_id = aws_apigatewayv2_api.main.id

  route_key = "POST /auth/signin"

  target = "integrations/${aws_apigatewayv2_integration.lambda_auth.id}"
}

resource "aws_apigatewayv2_route" "post_user" {
  api_id = aws_apigatewayv2_api.main.id

  route_key = "POST /user"

  target = "integrations/${aws_apigatewayv2_integration.lambda_auth.id}"
}
