resource "aws_iam_role" "auth_lambda_exec" {
  name = "auth-lambda"

  assume_role_policy = data.aws_iam_policy_document.lambda_exec_policy.json
}

resource "aws_iam_role_policy_attachment" "auth_lambda_policy" {
  role = aws_iam_role.auth_lambda_exec.name

  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_lambda_permission" "auth_lambda_gtw" {
  statement_id = "AllowExecutionFromApiGateway"

  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.auth.function_name
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_apigatewayv2_api.main.execution_arn}/*/*"
}

resource "aws_lambda_function" "auth" {
  function_name = "auth"

  s3_bucket = aws_s3_bucket.lambda_bucket.id
  s3_key    = aws_s3_object.lambda_auth.key

  runtime = "go1.x"
  handler = "main"

  environment {
    variables = {
      BCRYPT_SALT_LEN     = 12
      MONGO_URI           = var.users_mongo_database_uri
      MONGO_DATABASE_NAME = var.users_mongo_database_name
      WEBHOOK_SIG         = var.webhook_signature
      API_GATEWAY_PREFIX  = var.environment_type
      JWT_SIG             = var.jwt_signature
    }
  }

  memory_size = 512

  timeout = 8

  source_code_hash = data.archive_file.lambda_auth.output_base64sha256

  role = aws_iam_role.auth_lambda_exec.arn
}

resource "aws_cloudwatch_log_group" "auth" {
  name = "/aws/lambda/${aws_lambda_function.auth.function_name}"

  retention_in_days = 7
}

data "archive_file" "lambda_auth" {
  type = "zip"

  source_dir  = "../${path.module}/services/auth_lambda/dist"
  output_path = "../${path.module}/services/dist/auth_lambda.zip"
}

resource "aws_s3_object" "lambda_auth" {
  bucket = aws_s3_bucket.lambda_bucket.id

  key    = "auth.zip"
  source = data.archive_file.lambda_auth.output_path

  etag = filemd5(data.archive_file.lambda_auth.output_path)
}
