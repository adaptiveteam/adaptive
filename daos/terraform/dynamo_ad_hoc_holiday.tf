resource "aws_dynamodb_table" "ad_hoc_holiday_dynamodb_table"  {
	name           = "${var.client_id}_ad_hoc_holiday"
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
	    name = "date"
	    type = "S"
	}
	global_secondary_index {
		name            = "PlatformIDDateIndex"
		hash_key        = "platform_id"
		range_key       = "date"
		projection_type = "ALL"
		write_capacity  = var.dynamo_ondemand_write_capacity
		read_capacity   = var.dynamo_ondemand_read_capacity
	}
}
output "ad_hoc_holiday_table_arn" {
	description = "ARN of the ad_hoc_holiday table"
	value = aws_dynamodb_table.ad_hoc_holiday_dynamodb_table.arn
}
output "ad_hoc_holiday_table_name" {
	description = "Name of the ad_hoc_holiday table"
	value = aws_dynamodb_table.ad_hoc_holiday_dynamodb_table.name
}
