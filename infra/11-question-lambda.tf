resource "random_pet" "lambda_question_name" {
  prefix = "games4l-question"
  length = 2
}

resource "aws_iam_role" "question_lambda_exec" {
  name = random_pet.lambda_question_name.id

  assume_role_policy = data.aws_iam_policy_document.lambda_exec_policy.json
}

resource "aws_iam_role_policy_attachment" "question_lambda_policy" {
  role = aws_iam_role.question_lambda_exec.name

  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_lambda_permission" "question_lambda_gtw" {
  statement_id = "AllowExecutionFromApiGateway"

  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.question.function_name
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_apigatewayv2_api.main.execution_arn}/*/*"
}

resource "aws_lambda_function" "question" {
  function_name = random_pet.lambda_question_name.id

  s3_bucket = aws_s3_bucket.lambda_bucket.id
  s3_key    = aws_s3_object.lambda_question.key

  runtime = "go1.x"
  handler = "main"

  environment {
    variables = {
      MONGO_URI           = var.question_mongo_database_uri
      MONGO_DATABASE_NAME = var.question_mongo_database_name
      WEBHOOK_SIG         = var.webhook_signature
      API_GATEWAY_PREFIX  = var.environment_type
      JWT_SIG             = var.jwt_signature
    }
  }

  memory_size = 512

  timeout = 8

  source_code_hash = data.archive_file.lambda_question.output_base64sha256

  role = aws_iam_role.question_lambda_exec.arn
}

resource "aws_cloudwatch_log_group" "question" {
  name = "/aws/lambda/${aws_lambda_function.question.function_name}"

  retention_in_days = 7
}

data "archive_file" "lambda_question" {
  type = "zip"

  source_dir  = "../${path.module}/services/question_lambda/dist"
  output_path = "../${path.module}/services/dist/question_lambda.zip"
}

resource "aws_s3_object" "lambda_question" {
  bucket = aws_s3_bucket.lambda_bucket.id

  key    = "question.zip"
  source = data.archive_file.lambda_question.output_path

  etag = filemd5(data.archive_file.lambda_question.output_path)
}
