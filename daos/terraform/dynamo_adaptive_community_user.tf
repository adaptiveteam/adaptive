resource "aws_dynamodb_table" "adaptive_community_user_dynamodb_table"  {
	name           = "${var.client_id}_adaptive_community_user"
	billing_mode   = "PAY_PER_REQUEST"
	
	tags           = local.default_tags
	hash_key       = "channel_id"
	range_key      = "user_id"
	point_in_time_recovery {
		enabled = true
	}
	attribute {
	    name = "channel_id"
	    type = "S"
	}
	attribute {
	    name = "user_id"
	    type = "S"
	}
	attribute {
	    name = "platform_id"
	    type = "S"
	}
	attribute {
	    name = "community_id"
	    type = "S"
	}
	global_secondary_index {
		name            = "ChannelIDIndex"
		hash_key        = "channel_id"
		
		projection_type = "ALL"
		write_capacity  = var.dynamo_ondemand_write_capacity
		read_capacity   = var.dynamo_ondemand_read_capacity
	}
	global_secondary_index {
		name            = "UserIDCommunityIDIndex"
		hash_key        = "user_id"
		range_key       = "community_id"
		projection_type = "ALL"
		write_capacity  = var.dynamo_ondemand_write_capacity
		read_capacity   = var.dynamo_ondemand_read_capacity
	}
	global_secondary_index {
		name            = "UserIDIndex"
		hash_key        = "user_id"
		
		projection_type = "ALL"
		write_capacity  = var.dynamo_ondemand_write_capacity
		read_capacity   = var.dynamo_ondemand_read_capacity
	}
	global_secondary_index {
		name            = "PlatformIDCommunityIDIndex"
		hash_key        = "platform_id"
		range_key       = "community_id"
		projection_type = "ALL"
		write_capacity  = var.dynamo_ondemand_write_capacity
		read_capacity   = var.dynamo_ondemand_read_capacity
	}
}
