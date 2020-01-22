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
output "client_platform_token_table_arn" {
	description = "ARN of the client_platform_token table"
	value = aws_dynamodb_table.client_platform_token_dynamodb_table.arn
}
output "client_platform_token_table_name" {
	description = "Name of the client_platform_token table"
	value = aws_dynamodb_table.client_platform_token_dynamodb_table.name
}
