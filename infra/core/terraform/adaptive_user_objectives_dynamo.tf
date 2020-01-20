# resource "aws_dynamodb_table" "user_objectives_progress" {
#   name           = "${var.client_id}_user_objectives_progress"
#   read_capacity  = var.default_dynamo_read_capacity
#   write_capacity = var.default_dynamo_write_capacity
#   hash_key       = "id"
#   range_key      = "created_on"

#   attribute {
#     name = "id"
#     type = "S"
#   }

#   attribute {
#     name = "created_on"
#     type = "S"
#   }

#   stream_enabled   = true
#   stream_view_type = var.dynamo_stream_view_type

#   global_secondary_index {
#     name            = var.dynamo_user_objectives_progress_index
#     hash_key        = "id"
#     projection_type = "ALL"
#     write_capacity  = var.default_dynamo_read_capacity
#     read_capacity   = var.default_dynamo_write_capacity
#   }

#   global_secondary_index {
#     name            = var.dynamo_user_objectives_progress_created_on_index
#     hash_key        = "created_on"
#     projection_type = "ALL"
#     write_capacity  = var.default_dynamo_read_capacity
#     read_capacity   = var.default_dynamo_write_capacity
#   }

#   tags = local.default_tags
# }

# module "user_objectives_progress_table_scaling" {
#   source = "github.com/adaptiveteam/adaptive-terraform-modules//adaptive-dynamo-autoscaler?ref=v1.0.3"

#   client_id  = var.client_id
#   table_name = aws_dynamodb_table.user_objectives_progress.name
#   table_arn  = aws_dynamodb_table.user_objectives_progress.arn
# }

resource "aws_dynamodb_table" "accountability_partnership_rejections_table" {
  name           = "${var.client_id}_partnership_rejections"
  billing_mode = "PAY_PER_REQUEST"

  hash_key       = "objective_id"
  range_key      = "created_on"

  attribute {
    name = "objective_id"
    type = "S"
  }

  attribute {
    name = "created_on"
    type = "S"
  }

  stream_enabled   = true
  stream_view_type = var.dynamo_stream_view_type

  tags = local.default_tags
}

//module "backup_user_objective" {
//  source = "github.com/adaptiveteam/adaptive-terraform-modules//adaptive-backup?ref=backup-module"
//  bucket = aws_s3_bucket.backup_bucket.bucket
//  name_prefix = "${local.client_id}UserObjective"
//  role_name = aws_iam_role.role.name
//  resource_role_name = aws_iam_role.resource_role.name
//  table_name = aws_dynamodb_table.user_objective_dynamodb_table.name
//}
