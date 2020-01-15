resource "aws_dynamodb_table" "adaptive_dialog_content" {
  name           = "${var.client_id}_dialog_content"
  billing_mode = "PAY_PER_REQUEST"

  hash_key       = "dialog_id"

  attribute {
    name = "dialog_id"
    type = "S"
  }
  attribute {
    name = "context"
    type = "S"
  }
  attribute {
    name = "subject"
    type = "S"
  }

  server_side_encryption {
    enabled = true
  }

  tags = local.default_tags

  global_secondary_index {
    name            = var.dynamo_dialog_content_contect_subject_index
    hash_key        = "context"
    range_key       = "subject"
    read_capacity  = var.dynamo_ondemand_read_capacity
    write_capacity  = var.dynamo_ondemand_write_capacity
    projection_type = "ALL"
  }
}

resource "aws_dynamodb_table" "adaptive_dialog_aliases" {
  name           = "${aws_dynamodb_table.adaptive_dialog_content.name}_alias"
  billing_mode = "PAY_PER_REQUEST"

  hash_key       = "application_alias"

  attribute {
    name = "application_alias"
    type = "S"
  }

  tags = local.default_tags
}
