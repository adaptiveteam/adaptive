resource "aws_dynamodb_table" "adaptive_user_feedback_dynamodb_table" {
  name           = "${var.client_id}_adaptive_user_feedback"
  billing_mode = "PAY_PER_REQUEST"

  hash_key       = "id"

  attribute {
    name = "id"
    type = "S"
  }

  attribute {
    name = "target"
    type = "S"
  }

  attribute {
    name = "source"
    type = "S"
  }

  attribute {
    name = "quarter_year"
    type = "S"
  }

  stream_enabled   = true
  stream_view_type = var.dynamo_stream_view_type

  global_secondary_index {
    name            = var.feedback_target_quarter_year_index
    hash_key        = "quarter_year"
    range_key       = "target"
    projection_type = "ALL"
    write_capacity  = var.dynamo_ondemand_write_capacity
    read_capacity  = var.dynamo_ondemand_read_capacity
  }

  global_secondary_index {
    name            = var.feedback_source_quarter_year_index
    hash_key        = "quarter_year"
    range_key       = "source"
    projection_type = "ALL"
    write_capacity  = var.dynamo_ondemand_write_capacity
    read_capacity  = var.dynamo_ondemand_read_capacity
  }

  tags = local.default_tags
}
