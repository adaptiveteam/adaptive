resource "aws_dynamodb_table" "adaptive_users_dynamodb_table" {
  name           = "${var.client_id}_adaptive_users"
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
    name = "timezone_offset"
    type = "N"
  }
  attribute {
    name = "adaptive_scheduled_time_in_utc"
    type = "S"
  }
  global_secondary_index {
    name            = var.dynamo_users_platform_index
    hash_key        = "platform_id"
    write_capacity  = var.dynamo_ondemand_write_capacity
    read_capacity  = var.dynamo_ondemand_read_capacity
    projection_type = "ALL"
  }
  global_secondary_index {
    name            = var.dynamo_users_timezone_offset_index
    hash_key        = "platform_id"
    range_key       = "timezone_offset"
    projection_type = "ALL"
    write_capacity  = var.dynamo_ondemand_write_capacity
    read_capacity  = var.dynamo_ondemand_read_capacity
  }
  global_secondary_index {
    name            = var.dynamo_users_scheduled_time_index
    hash_key        = "platform_id"
    range_key       = "adaptive_scheduled_time_in_utc"
    projection_type = "ALL"
    write_capacity  = var.dynamo_ondemand_write_capacity
    read_capacity  = var.dynamo_ondemand_read_capacity
  }

  stream_enabled   = true
  stream_view_type = var.dynamo_stream_view_type

  tags = local.default_tags
}

//module "backup_user" {
//  source = "github.com/adaptiveteam/adaptive-terraform-modules//adaptive-backup?ref=backup-module"
//  bucket = aws_s3_bucket.backup_bucket.bucket
//  name_prefix = "${local.client_id}User"
//  role_name = aws_iam_role.role.name
//  resource_role_name = aws_iam_role.resource_role.name
//  table_name = aws_dynamodb_table.adaptive_users_dynamodb_table.name
//}
