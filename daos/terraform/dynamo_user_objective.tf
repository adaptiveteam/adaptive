resource "aws_dynamodb_table" "user_objective_dynamodb_table"  {
	name           = "${var.client_id}_user_objective"
	billing_mode   = "PAY_PER_REQUEST"
	
	tags           = local.default_tags
	hash_key       = "id"
	point_in_time_recovery {
		enabled = true
	}
	attribute {
	    name = "id"
	    type = "S"
	}
	attribute {
	    name = "user_id"
	    type = "S"
	}
	attribute {
	    name = "accountability_partner"
	    type = "S"
	}
	attribute {
	    name = "accepted"
	    type = "N"
	}
	attribute {
	    name = "type"
	    type = "S"
	}
	attribute {
	    name = "completed"
	    type = "N"
	}
	global_secondary_index {
		name            = "UserIDCompletedIndex"
		hash_key        = "user_id"
		range_key       = "completed"
		projection_type = "ALL"
		write_capacity  = var.dynamo_ondemand_write_capacity
		read_capacity   = var.dynamo_ondemand_read_capacity
	}
	global_secondary_index {
		name            = "AcceptedIndex"
		hash_key        = "accepted"
		
		projection_type = "ALL"
		write_capacity  = var.dynamo_ondemand_write_capacity
		read_capacity   = var.dynamo_ondemand_read_capacity
	}
	global_secondary_index {
		name            = "AccountabilityPartnerIndex"
		hash_key        = "accountability_partner"
		
		projection_type = "ALL"
		write_capacity  = var.dynamo_ondemand_write_capacity
		read_capacity   = var.dynamo_ondemand_read_capacity
	}
	global_secondary_index {
		name            = "UserIDTypeIndex"
		hash_key        = "user_id"
		range_key       = "type"
		projection_type = "ALL"
		write_capacity  = var.dynamo_ondemand_write_capacity
		read_capacity   = var.dynamo_ondemand_read_capacity
	}
}
output "user_objective_table_arn" {
	description = "ARN of the user_objective table"
	value = aws_dynamodb_table.user_objective_dynamodb_table.arn
}
output "user_objective_table_name" {
	description = "Name of the user_objective table"
	value = aws_dynamodb_table.user_objective_dynamodb_table.name
}
