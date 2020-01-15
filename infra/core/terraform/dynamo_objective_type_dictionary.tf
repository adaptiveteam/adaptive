resource "aws_dynamodb_table" "objective_type_dictionary_dynamodb_table"  {
	name           = "${var.client_id}_objective_type_dictionary"
	billing_mode = "PAY_PER_REQUEST"

	tags           = local.default_tags
	hash_key       = "id"
	range_key      = "platform_id"
	point_in_time_recovery {
		enabled = true
	}
	attribute {
	    name = "id"
	    type = "S"
	}
	attribute {
	    name = "platform_id"
	    type = "S"
	}
	global_secondary_index {
		name            = "PlatformIDIndex"
		hash_key        = "platform_id"
		
		projection_type = "ALL"
		write_capacity  = var.dynamo_ondemand_write_capacity
		read_capacity  = var.dynamo_ondemand_read_capacity
	}
}

# module "backup_objective_type_dictionary" {
#   source = "github.com/adaptiveteam/adaptive-terraform-modules//adaptive-backup?ref=backup-module"
#   bucket = aws_s3_bucket.backup_bucket.bucket
#   name_prefix = "${local.client_id}ObjectiveTypeDictionary"
#   role_name = aws_iam_role.role.name
#   resource_role_name = aws_iam_role.resource_role.name
#   table_name = aws_dynamodb_table.objective_type_dictionary_dynamodb_table.name
# }
