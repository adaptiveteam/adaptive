resource "aws_dynamodb_table" "community_users" {
  name           = "${var.client_id}_community_users"
  billing_mode = "PAY_PER_REQUEST"

  hash_key       = "channel_id"
  range_key      = "user_id"

  attribute {
    name = "channel_id"
    type = "S"
  }

  attribute {
    name = "user_id"
    type = "S"
  }

  attribute {
    name = "community_id"
    type = "S"
  }
  attribute {
    name = "platform_id"
    type = "S"
  }

  stream_enabled   = true
  stream_view_type = var.dynamo_stream_view_type

  global_secondary_index {
    name            = var.dynamo_community_users_channel_index
    hash_key        = "channel_id"
    projection_type = "ALL"
    write_capacity  = var.dynamo_ondemand_write_capacity
    read_capacity  = var.dynamo_ondemand_read_capacity
  }

  global_secondary_index {
    name            = var.dynamo_community_users_user_community_index
    hash_key        = "user_id"
    range_key       = "community_id"
    projection_type = "ALL"
    write_capacity  = var.dynamo_ondemand_write_capacity
    read_capacity  = var.dynamo_ondemand_read_capacity
  }

  global_secondary_index {
    name            = var.dynamo_community_users_user_index
    hash_key        = "user_id"
    projection_type = "ALL"
    write_capacity  = var.dynamo_ondemand_write_capacity
    read_capacity  = var.dynamo_ondemand_read_capacity
  }

  global_secondary_index {
    name            = var.dynamo_community_users_community_index
    hash_key        = "platform_id"
    range_key       = "community_id"
    projection_type = "ALL"
    write_capacity  = var.dynamo_ondemand_write_capacity
    read_capacity  = var.dynamo_ondemand_read_capacity
  }

  tags = local.default_tags
}
