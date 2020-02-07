resource "aws_dynamodb_table" "adaptive_user_engagements_dynamo_table" {
  name = "${var.client_id}_user_engagement"
  billing_mode = "PAY_PER_REQUEST"

  hash_key = "user_id"
  range_key = "id"

  attribute {
    name = "user_id"
    type = "S"
  }

  # id of the engagement
  attribute {
    name = "id"
    type = "S"
  }

  stream_enabled = true
  stream_view_type = var.dynamo_stream_view_type

  # number indicating if an engagement is answered, 0 or 1
  # terraform doesn't support boolean attributes yet for dynamo
  # this is used for Global Secondary Index
  attribute {
    name = "answered"
    type = "N"
  }

  local_secondary_index {
    name = var.user_engagement_answered_dynamo_index
    range_key = "answered"
    projection_type = "ALL"
  }

  tags = local.default_tags
}
