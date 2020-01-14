resource "aws_dynamodb_table" "coaching_rejections" {
  name           = "${var.client_id}_coaching_rejections"
  billing_mode = "PAY_PER_REQUEST"

  hash_key       = "id"

  attribute {
    name = "id"
    type = "S"
  }

  tags = local.default_tags
}
