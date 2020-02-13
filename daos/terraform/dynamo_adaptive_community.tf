resource "aws_dynamodb_table" "adaptive_community_dynamodb_table"  {
	name           = "${var.client_id}_adaptive_community"
	billing_mode   = "PAY_PER_REQUEST"
	
	tags           = local.default_tags
	hash_key       = "id"
	range_key      = "platform_id"
	point_in_time_recovery {
		enabled = true
	}
	attribute {
	    name = "platform_id"
	    type = "S"
	}
	attribute {
	    name = "id"
	    type = "S"
	}
	attribute {
	    name = "channel"
	    type = "S"
	}
	global_secondary_index {
		name            = "ChannelIndex"
		hash_key        = "channel"
		
		projection_type = "ALL"
		write_capacity  = var.dynamo_ondemand_write_capacity
		read_capacity   = var.dynamo_ondemand_read_capacity
	}
	global_secondary_index {
		name            = "PlatformIDIndex"
		hash_key        = "platform_id"
		
		projection_type = "ALL"
		write_capacity  = var.dynamo_ondemand_write_capacity
		read_capacity   = var.dynamo_ondemand_read_capacity
	}
}
output "adaptive_community_table_arn" {
	description = "ARN of the adaptive_community table"
	value = aws_dynamodb_table.adaptive_community_dynamodb_table.arn
}
output "adaptive_community_table_name" {
	description = "Name of the adaptive_community table"
	value = aws_dynamodb_table.adaptive_community_dynamodb_table.name
}
