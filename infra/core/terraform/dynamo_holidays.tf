resource "aws_dynamodb_table" "ad_hoc_holidays" {
  name           = "${var.client_id}_ad_hoc_holidays"
  billing_mode = "PAY_PER_REQUEST"

  hash_key       = "id"

  attribute {
    name = "id"
    type = "S"
  }

  // GSI
  attribute {
    name = "platform_id"
    type = "S"
  }

  attribute {
    name = "date"
    type = "S"
  }

  stream_enabled   = true
  stream_view_type = var.dynamo_stream_view_type

  global_secondary_index {
    name            = var.dynamo_holidays_date_index
    hash_key        = "platform_id"
    range_key       = "date"
    projection_type = "ALL"
    write_capacity  = var.dynamo_ondemand_write_capacity
    read_capacity  = var.dynamo_ondemand_read_capacity
  }

  tags = local.default_tags
}
