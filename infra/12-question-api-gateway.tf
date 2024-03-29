resource "aws_apigatewayv2_integration" "lambda_question" {
  api_id = aws_apigatewayv2_api.main.id

  integration_uri    = aws_lambda_function.question.invoke_arn
  integration_type   = "AWS_PROXY"
  integration_method = "POST"
}

resource "aws_apigatewayv2_route" "get_question" {
  api_id = aws_apigatewayv2_api.main.id

  route_key = "GET /question/{id}"

  target = "integrations/${aws_apigatewayv2_integration.lambda_question.id}"
}

resource "aws_apigatewayv2_route" "get_questions" {
  api_id = aws_apigatewayv2_api.main.id

  route_key = "GET /question"

  target = "integrations/${aws_apigatewayv2_integration.lambda_question.id}"
}

resource "aws_apigatewayv2_route" "put_question" {
  api_id = aws_apigatewayv2_api.main.id

  route_key = "PUT /question/{id}"

  target = "integrations/${aws_apigatewayv2_integration.lambda_question.id}"
}

resource "aws_apigatewayv2_route" "patch_question" {
  api_id = aws_apigatewayv2_api.main.id

  route_key = "PATCH /question/{id}"

  target = "integrations/${aws_apigatewayv2_integration.lambda_question.id}"
}

resource "aws_apigatewayv2_route" "post_question" {
  api_id = aws_apigatewayv2_api.main.id

  route_key = "POST /question"

  target = "integrations/${aws_apigatewayv2_integration.lambda_question.id}"
}
