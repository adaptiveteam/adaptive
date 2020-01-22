resource "aws_dynamodb_table" "strategy_community_dynamodb_table"  {
	name           = "${var.client_id}_strategy_community"
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
	    name = "platform_id"
	    type = "S"
	}
	attribute {
	    name = "channel_id"
	    type = "S"
	}
	attribute {
	    name = "channel_created"
	    type = "N"
	}
	global_secondary_index {
		name            = "PlatformIDChannelCreatedIndex"
		hash_key        = "platform_id"
		range_key       = "channel_created"
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
	global_secondary_index {
		name            = "ChannelIDIndex"
		hash_key        = "channel_id"
		
		projection_type = "ALL"
		write_capacity  = var.dynamo_ondemand_write_capacity
		read_capacity   = var.dynamo_ondemand_read_capacity
	}
}
output "strategy_community_table_arn" {
	description = "ARN of the strategy_community table"
	value = aws_dynamodb_table.strategy_community_dynamodb_table.arn
}
output "strategy_community_table_name" {
	description = "Name of the strategy_community table"
	value = aws_dynamodb_table.strategy_community_dynamodb_table.name
}
