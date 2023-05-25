resource "aws_iam_role" "telemetry_lambda_exec" {
  name = "telemetry-lambda"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "",
      "Effect": "Allow",
      "Principal": {
        "Service": [
          "apigateway.amazonaws.com",
          "lambda.amazonaws.com"
        ]
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "telemetry_lambda_policy" {
  role       = aws_iam_role.telemetry_lambda_exec.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_lambda_permission" "telemetry_lambda_gtw" {
  statement_id = "AllowExecutionFromApiGateway"

  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.telemetry.function_name
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_apigatewayv2_api.main.execution_arn}/*/*"
}

resource "aws_lambda_function" "telemetry" {
  function_name = "telemetry"

  s3_bucket = aws_s3_bucket.lambda_bucket.id
  s3_key    = aws_s3_object.lambda_telemetry.key

  runtime = "go1.x"
  handler = "main"

  source_code_hash = data.archive_file.lambda_telemetry.output_base64sha256

  role = aws_iam_role.telemetry_lambda_exec.arn
}

resource "aws_cloudwatch_log_group" "telemetry" {
  name = "/aws/lambda/${aws_lambda_function.telemetry.function_name}"

  retention_in_days = 7
}

data "archive_file" "lambda_telemetry" {
  type = "zip"

  source_dir  = "../${path.module}/services/telemetry_lambda/dist"
  output_path = "../${path.module}/services/dist/telemetry_lambda.zip"
}

resource "aws_s3_object" "lambda_telemetry" {
  bucket = aws_s3_bucket.lambda_bucket.id

  key    = "telemetry.zip"
  source = data.archive_file.lambda_telemetry.output_path

  etag = filemd5(data.archive_file.lambda_telemetry.output_path)
}
