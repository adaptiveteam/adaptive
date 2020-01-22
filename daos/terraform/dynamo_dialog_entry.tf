resource "aws_dynamodb_table" "dialog_entry_dynamodb_table"  {
	name           = "${var.client_id}_dialog_entry"
	billing_mode   = "PAY_PER_REQUEST"
	
	tags           = local.default_tags
	hash_key       = "dialog_id"
	point_in_time_recovery {
		enabled = true
	}
	server_side_encryption {
		enabled = true
	}
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
	global_secondary_index {
		name            = "ContextSubjectIndex"
		hash_key        = "context"
		range_key       = "subject"
		projection_type = "ALL"
		write_capacity  = var.dynamo_ondemand_write_capacity
		read_capacity   = var.dynamo_ondemand_read_capacity
	}
}
output "dialog_entry_table_arn" {
	description = "ARN of the dialog_entry table"
	value = aws_dynamodb_table.dialog_entry_dynamodb_table.arn
}
output "dialog_entry_table_name" {
	description = "Name of the dialog_entry table"
	value = aws_dynamodb_table.dialog_entry_dynamodb_table.name
}
