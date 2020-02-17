// policy for writing all tables related to issues workflow
resource "aws_iam_policy" "write_issues" {
  name   = "${var.client_id}_write_issues"
  policy = data.aws_iam_policy_document.write_issues.json
}

data "aws_iam_policy_document" "write_issues" {
  statement {
    actions = ["dynamodb:*"]
    resources = [
      aws_dynamodb_table.adaptive_user_engagements_dynamo_table.arn,
      aws_dynamodb_table.adaptive_users_dynamodb_table.arn,

      aws_dynamodb_table.postponed_event_dynamodb_table.arn,

      aws_dynamodb_table.user_objective_dynamodb_table.arn,
      aws_dynamodb_table.user_objectives_progress.arn,

      aws_dynamodb_table.strategy_objectives.arn,
      aws_dynamodb_table.strategy_initiatives.arn,

      aws_dynamodb_table.accountability_partnership_rejections_table.arn,
      
    ]
  }
}
