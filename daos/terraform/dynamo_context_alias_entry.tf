resource "aws_dynamodb_table" "context_alias_entry_dynamodb_table"  {
	name           = "${var.client_id}_context_alias_entry"
	billing_mode   = "PAY_PER_REQUEST"
	
	tags           = local.default_tags
	hash_key       = "application_alias"
	point_in_time_recovery {
		enabled = true
	}
	attribute {
	    name = "application_alias"
	    type = "S"
	}
}
output "context_alias_entry_table_arn" {
	description = "ARN of the context_alias_entry table"
	value = aws_dynamodb_table.context_alias_entry_dynamodb_table.arn
}
output "context_alias_entry_table_name" {
	description = "Name of the context_alias_entry table"
	value = aws_dynamodb_table.context_alias_entry_dynamodb_table.name
}
