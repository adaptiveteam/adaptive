resource "aws_dynamodb_table" "client_config_dynamodb_table" {
  name           = "${var.client_id}_adaptive_client_config"
  billing_mode = "PAY_PER_REQUEST"

  hash_key       = "platform_id"

  attribute {
    name = "platform_id"
    type = "S"
  }

  stream_enabled   = true
  stream_view_type = var.dynamo_stream_view_type

  tags = local.default_tags
}
