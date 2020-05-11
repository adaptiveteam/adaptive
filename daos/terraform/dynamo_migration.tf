resource "aws_dynamodb_table" "migration_dynamodb_table"  {
	name           = "${var.client_id}_migration"
	billing_mode   = "PAY_PER_REQUEST"
	
	tags           = local.default_tags
	hash_key       = "platform_id"
	range_key      = "migration_id"
	point_in_time_recovery {
		enabled = true
	}
	attribute {
	    name = "platform_id"
	    type = "S"
	}
	attribute {
	    name = "migration_id"
	    type = "S"
	}
}
output "migration_table_arn" {
	description = "ARN of the migration table"
	value = aws_dynamodb_table.migration_dynamodb_table.arn
}
output "migration_table_name" {
	description = "Name of the migration table"
	value = aws_dynamodb_table.migration_dynamodb_table.name
}
