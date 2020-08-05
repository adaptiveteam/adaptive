resource "aws_dynamodb_table" "capability_communities" {
  name           = "${var.client_id}_capability_communities"
  billing_mode = "PAY_PER_REQUEST"

  hash_key       = "id"
  range_key      = "platform_id"
  stream_enabled   = true
  stream_view_type = var.dynamo_stream_view_type

  attribute {
    name = "id"
    type = "S"
  }
  attribute {
    name = "platform_id"
    type = "S"
  }

  global_secondary_index {
    name            = var.dynamo_capability_communities_platform_index
    hash_key        = "platform_id"
    projection_type = "ALL"
    write_capacity  = var.dynamo_ondemand_write_capacity
    read_capacity  = var.dynamo_ondemand_read_capacity
  }

  tags = local.default_tags
}
