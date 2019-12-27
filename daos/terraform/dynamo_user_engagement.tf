resource "aws_dynamodb_table" "user_engagement_dynamodb_table"  {
	name           = "${var.client_id}_user_engagement"
	billing_mode   = "PAY_PER_REQUEST"
	
	tags           = local.default_tags
	hash_key       = "user_id"
	range_key      = "id"
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
	    name = "answered"
	    type = "N"
	}
	local_secondary_index {
		name            = "UserIDAnsweredIndex"
		projection_type = "INCLUDE"
		range_key       = "answered"
		non_key_attributes =  [
			"script",
			"priority",
			"target_id",
			"ignored",
		]
	}
}
