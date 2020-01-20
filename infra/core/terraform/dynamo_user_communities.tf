resource "aws_dynamodb_table" "user_communities" {
  name           = "${var.client_id}_communities"
  billing_mode = "PAY_PER_REQUEST"

  hash_key       = "id"
  range_key      = "platform_id"

  attribute {
    name = "id"
    type = "S"
  }

  attribute {
    name = "platform_id"
    type = "S"
  }

  // For GSI
  attribute {
    name = "channel"
    type = "S"
  }

  stream_enabled   = true
  stream_view_type = var.dynamo_stream_view_type

  global_secondary_index {
    name            = var.user_community_channel_dynamo_index
    hash_key        = "channel"
    projection_type = "ALL"
    write_capacity  = var.dynamo_ondemand_write_capacity
    read_capacity  = var.dynamo_ondemand_read_capacity
  }

  global_secondary_index {
    name            = var.user_community_platform_dynamo_index
    hash_key        = "platform_id"
    projection_type = "ALL"
    write_capacity  = var.dynamo_ondemand_write_capacity
    read_capacity  = var.dynamo_ondemand_read_capacity
  }

  tags = local.default_tags
}
