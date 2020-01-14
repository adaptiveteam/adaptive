resource "aws_dynamodb_table" "strategy_communities" {
  name           = "${var.client_id}_strategy_communities"
  billing_mode = "PAY_PER_REQUEST"

  hash_key       = "id"

  attribute {
    name = "id"
    type = "S"
  }

  attribute {
    name = "channel_created"
    type = "N"
  }
  attribute {
    name = "platform_id"
    type = "S"
  }
  attribute {
    name = "channel_id"
    type = "S"
  }

  stream_enabled   = true
  stream_view_type = var.dynamo_stream_view_type

  global_secondary_index {
    name            = var.dynamo_strategy_communities_platform_channel_created_index
    hash_key        = "platform_id"
    range_key       = "channel_created"
    projection_type = "ALL"
    write_capacity  = var.dynamo_ondemand_write_capacity
    read_capacity  = var.dynamo_ondemand_read_capacity
  }

  global_secondary_index {
    name            = var.dynamo_strategy_communities_platform_index
    hash_key        = "platform_id"
    projection_type = "ALL"
    write_capacity  = var.dynamo_ondemand_write_capacity
    read_capacity  = var.dynamo_ondemand_read_capacity
  }

  global_secondary_index {
    name            = var.dynamo_strategy_communities_channel_index
    hash_key        = "channel_id"
    projection_type = "ALL"
    write_capacity  = var.dynamo_ondemand_write_capacity
    read_capacity  = var.dynamo_ondemand_read_capacity
  }

  tags = local.default_tags
}
