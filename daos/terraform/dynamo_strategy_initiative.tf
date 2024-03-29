resource "aws_dynamodb_table" "strategy_initiative_dynamodb_table"  {
	name           = "${var.client_id}_strategy_initiative"
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
	    name = "initiative_community_id"
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
		name            = "InitiativeCommunityIDIndex"
		hash_key        = "initiative_community_id"
		
		projection_type = "ALL"
		write_capacity  = var.dynamo_ondemand_write_capacity
		read_capacity   = var.dynamo_ondemand_read_capacity
	}
}
output "strategy_initiative_table_arn" {
	description = "ARN of the strategy_initiative table"
	value = aws_dynamodb_table.strategy_initiative_dynamodb_table.arn
}
output "strategy_initiative_table_name" {
	description = "Name of the strategy_initiative table"
	value = aws_dynamodb_table.strategy_initiative_dynamodb_table.name
}
