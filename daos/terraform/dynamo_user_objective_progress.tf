resource "aws_dynamodb_table" "user_objective_progress_dynamodb_table"  {
	name           = "${var.client_id}_user_objective_progress"
	billing_mode   = "PAY_PER_REQUEST"
	
	tags           = local.default_tags
	hash_key       = "id"
	range_key      = "created_on"
	point_in_time_recovery {
		enabled = true
	}
	attribute {
	    name = "id"
	    type = "S"
	}
	attribute {
	    name = "created_on"
	    type = "S"
	}
	global_secondary_index {
		name            = "IDIndex"
		hash_key        = "id"
		
		projection_type = "ALL"
		write_capacity  = var.dynamo_ondemand_write_capacity
		read_capacity   = var.dynamo_ondemand_read_capacity
	}
	global_secondary_index {
		name            = "CreatedOnIndex"
		hash_key        = "created_on"
		
		projection_type = "ALL"
		write_capacity  = var.dynamo_ondemand_write_capacity
		read_capacity   = var.dynamo_ondemand_read_capacity
	}
}
output "user_objective_progress_table_arn" {
	description = "ARN of the user_objective_progress table"
	value = aws_dynamodb_table.user_objective_progress_dynamodb_table.arn
}
output "user_objective_progress_table_name" {
	description = "Name of the user_objective_progress table"
	value = aws_dynamodb_table.user_objective_progress_dynamodb_table.name
}
