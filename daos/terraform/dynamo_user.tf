resource "aws_dynamodb_table" "user_dynamodb_table"  {
	name           = "${var.client_id}_user"
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
	    name = "timezone_offset"
	    type = "N"
	}
	attribute {
	    name = "adaptive_scheduled_time_in_utc"
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
		read_capacity   = var.dynamo_ondemand_read_capacity
	}
	global_secondary_index {
		name            = "PlatformIDTimezoneOffsetIndex"
		hash_key        = "platform_id"
		range_key       = "timezone_offset"
		projection_type = "ALL"
		write_capacity  = var.dynamo_ondemand_write_capacity
		read_capacity   = var.dynamo_ondemand_read_capacity
	}
	global_secondary_index {
		name            = "PlatformIDAdaptiveScheduledTimeInUTCIndex"
		hash_key        = "platform_id"
		range_key       = "adaptive_scheduled_time_in_utc"
		projection_type = "ALL"
		write_capacity  = var.dynamo_ondemand_write_capacity
		read_capacity   = var.dynamo_ondemand_read_capacity
	}
}
output "user_table_arn" {
	description = "ARN of the user table"
	value = aws_dynamodb_table.user_dynamodb_table.arn
}
output "user_table_name" {
	description = "Name of the user table"
	value = aws_dynamodb_table.user_dynamodb_table.name
}
