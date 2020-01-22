resource "aws_dynamodb_table" "user_feedback_dynamodb_table"  {
	name           = "${var.client_id}_user_feedback"
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
	    name = "source"
	    type = "S"
	}
	attribute {
	    name = "target"
	    type = "S"
	}
	attribute {
	    name = "quarter_year"
	    type = "S"
	}
	global_secondary_index {
		name            = "QuarterYearSourceIndex"
		hash_key        = "quarter_year"
		range_key       = "source"
		projection_type = "ALL"
		write_capacity  = var.dynamo_ondemand_write_capacity
		read_capacity   = var.dynamo_ondemand_read_capacity
	}
	global_secondary_index {
		name            = "QuarterYearTargetIndex"
		hash_key        = "quarter_year"
		range_key       = "target"
		projection_type = "ALL"
		write_capacity  = var.dynamo_ondemand_write_capacity
		read_capacity   = var.dynamo_ondemand_read_capacity
	}
}
output "user_feedback_table_arn" {
	description = "ARN of the user_feedback table"
	value = aws_dynamodb_table.user_feedback_dynamodb_table.arn
}
output "user_feedback_table_name" {
	description = "Name of the user_feedback table"
	value = aws_dynamodb_table.user_feedback_dynamodb_table.name
}
