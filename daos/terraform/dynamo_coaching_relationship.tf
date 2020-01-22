resource "aws_dynamodb_table" "coaching_relationship_dynamodb_table"  {
	name           = "${var.client_id}_coaching_relationship"
	billing_mode   = "PAY_PER_REQUEST"
	
	tags           = local.default_tags
	hash_key       = "coach_quarter_year"
	range_key      = "coachee"
	point_in_time_recovery {
		enabled = true
	}
	attribute {
	    name = "coach_quarter_year"
	    type = "S"
	}
	attribute {
	    name = "coachee_quarter_year"
	    type = "S"
	}
	attribute {
	    name = "coachee"
	    type = "S"
	}
	attribute {
	    name = "quarter"
	    type = "N"
	}
	attribute {
	    name = "year"
	    type = "N"
	}
	global_secondary_index {
		name            = "CoachQuarterYearIndex"
		hash_key        = "coach_quarter_year"
		
		projection_type = "ALL"
		write_capacity  = var.dynamo_ondemand_write_capacity
		read_capacity   = var.dynamo_ondemand_read_capacity
	}
	global_secondary_index {
		name            = "QuarterYearIndex"
		hash_key        = "quarter"
		range_key       = "year"
		projection_type = "ALL"
		write_capacity  = var.dynamo_ondemand_write_capacity
		read_capacity   = var.dynamo_ondemand_read_capacity
	}
	global_secondary_index {
		name            = "CoacheeQuarterYearIndex"
		hash_key        = "coachee_quarter_year"
		
		projection_type = "ALL"
		write_capacity  = var.dynamo_ondemand_write_capacity
		read_capacity   = var.dynamo_ondemand_read_capacity
	}
}
output "coaching_relationship_table_arn" {
	description = "ARN of the coaching_relationship table"
	value = aws_dynamodb_table.coaching_relationship_dynamodb_table.arn
}
output "coaching_relationship_table_name" {
	description = "Name of the coaching_relationship table"
	value = aws_dynamodb_table.coaching_relationship_dynamodb_table.name
}
