resource "aws_dynamodb_table" "client_platform_token_dynamodb_table"  {
	name           = "${var.client_id}_client_platform_token"
	billing_mode   = "PAY_PER_REQUEST"
	
	tags           = local.default_tags
	hash_key       = "platform_id"
	point_in_time_recovery {
		enabled = true
	}
	attribute {
	    name = "platform_id"
	    type = "S"
	}
}
