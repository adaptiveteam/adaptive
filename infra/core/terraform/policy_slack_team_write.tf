resource "aws_iam_policy" "slack_team_write" {
  name   = "${var.client_id}_slack_team_write"
  policy = data.aws_iam_policy_document.slack_team_write.json
}

data "aws_iam_policy_document" "slack_team_write" {
  statement {
    actions = ["dynamodb:*"]
    resources = [aws_dynamodb_table.slack_team_dynamodb_table.arn]
  }
}
