resource "aws_dynamodb_table" "vision_mission_dynamodb_table"  {
	name           = "${var.client_id}_vision_mission"
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
output "vision_mission_table_arn" {
	description = "ARN of the vision_mission table"
	value = aws_dynamodb_table.vision_mission_dynamodb_table.arn
}
output "vision_mission_table_name" {
	description = "Name of the vision_mission table"
	value = aws_dynamodb_table.vision_mission_dynamodb_table.name
}
