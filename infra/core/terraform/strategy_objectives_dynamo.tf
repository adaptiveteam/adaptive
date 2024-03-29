resource "aws_dynamodb_table" "strategy_objectives" {
  name           = "${var.client_id}_strategy_objectives"
  billing_mode = "PAY_PER_REQUEST"

  hash_key       = "id"
  range_key      = "platform_id"

  attribute {
    name = "id"
    type = "S"
  }

  // GSI
  attribute {
    name = "platform_id"
    type = "S"
  }

  stream_enabled   = true
  stream_view_type = var.dynamo_stream_view_type

  global_secondary_index {
    name            = var.dynamo_strategy_objectives_platform_index
    hash_key        = "platform_id"
    projection_type = "ALL"
    write_capacity  = var.dynamo_ondemand_write_capacity
    read_capacity  = var.dynamo_ondemand_read_capacity
  }

  tags = local.default_tags
}
