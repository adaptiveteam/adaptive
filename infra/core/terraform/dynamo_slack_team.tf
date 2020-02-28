resource "aws_dynamodb_table" "slack_team_dynamodb_table"  {
	name           = "${var.client_id}_slack_team"
	billing_mode   = "PAY_PER_REQUEST"
	
	tags           = local.default_tags
	hash_key       = "team_id"
	point_in_time_recovery {
		enabled = true
	}
	attribute {
	    name = "team_id"
	    type = "S"
	}
}
output "slack_team_table_arn" {
	description = "ARN of the slack_team table"
	value = aws_dynamodb_table.slack_team_dynamodb_table.arn
}
output "slack_team_table_name" {
	description = "Name of the slack_team table"
	value = aws_dynamodb_table.slack_team_dynamodb_table.name
}
