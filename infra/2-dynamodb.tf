resource "aws_dynamodb_table" "telemetry_table" {
  name = "games4l_telemetry"

  hash_key  = "pacient_name"
  range_key = "created_at"

  billing_mode = "PROVISIONED"

  read_capacity  = 5
  write_capacity = 5

  attribute {
    name = "pacient_name"
    type = "S"
  }

  attribute {
    name = "created_at"
    type = "S"
  }
}
