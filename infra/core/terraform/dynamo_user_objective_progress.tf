resource "aws_dynamodb_table" "user_objectives_progress" {
    name           = "${var.client_id}_user_objectives_progress"
	billing_mode = "PAY_PER_REQUEST"

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

	stream_enabled   = true
	stream_view_type = var.dynamo_stream_view_type

	global_secondary_index {
		name            = var.dynamo_user_objectives_progress_index
		hash_key        = "id"
		projection_type = "ALL"
		write_capacity  = var.dynamo_ondemand_write_capacity
		read_capacity  = var.dynamo_ondemand_read_capacity
	}
	global_secondary_index {
		name            = "CreatedOnIndex"
		hash_key        = "created_on"

		projection_type = "ALL"
		write_capacity  = var.dynamo_ondemand_write_capacity
		read_capacity  = var.dynamo_ondemand_read_capacity
	}
}
