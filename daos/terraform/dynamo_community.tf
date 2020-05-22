resource "aws_dynamodb_table" "community_dynamodb_table"  {
	name           = "${var.client_id}_community"
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
	    name = "channel_id"
	    type = "S"
	}
	attribute {
	    name = "community_kind"
	    type = "S"
	}
	global_secondary_index {
		name            = "ChannelIDPlatformIDIndex"
		hash_key        = "channel_id"
		range_key       = "platform_id"
		projection_type = "ALL"
		write_capacity  = var.dynamo_ondemand_write_capacity
		read_capacity   = var.dynamo_ondemand_read_capacity
	}
	global_secondary_index {
		name            = "PlatformIDCommunityKindIndex"
		hash_key        = "platform_id"
		range_key       = "community_kind"
		projection_type = "ALL"
		write_capacity  = var.dynamo_ondemand_write_capacity
		read_capacity   = var.dynamo_ondemand_read_capacity
	}
}
output "community_table_arn" {
	description = "ARN of the community table"
	value = aws_dynamodb_table.community_dynamodb_table.arn
}
output "community_table_name" {
	description = "Name of the community table"
	value = aws_dynamodb_table.community_dynamodb_table.name
}
