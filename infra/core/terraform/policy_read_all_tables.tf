resource "aws_iam_policy" "read_all_tables" {
  name   = "${var.client_id}_read_all_tables"
  policy = data.aws_iam_policy_document.read_all_tables.json
}

data "aws_iam_policy_document" "read_all_tables" {
  statement {
    actions = ["dynamodb:GetItem","dynamodb:Query","dynamodb:DescribeTable", "dynamodb:Scan"]
    resources = [
      aws_dynamodb_table.adaptive_dialog_content.arn,
      aws_dynamodb_table.adaptive_dialog_aliases.arn,
      aws_dynamodb_table.adaptive_user_engagements_dynamo_table.arn,
      aws_dynamodb_table.adaptive_users_dynamodb_table.arn,
      aws_dynamodb_table.adaptive_value_dynamodb_table.arn,
      aws_dynamodb_table.client_config_dynamodb_table.arn,
      aws_dynamodb_table.coaching_relationships.arn,
      aws_dynamodb_table.community_users.arn,
      aws_dynamodb_table.user_communities.arn,
      aws_dynamodb_table.user_objective_dynamodb_table.arn,
      aws_dynamodb_table.user_objectives_progress.arn,
      aws_dynamodb_table.strategy_communities.arn,
      aws_dynamodb_table.strategy_initiatives.arn,
      aws_dynamodb_table.strategy_objectives.arn,
      aws_dynamodb_table.vision.arn,

      aws_dynamodb_table.accountability_partnership_rejections_table.arn,
      aws_dynamodb_table.initiative_communities.arn,
      aws_dynamodb_table.capability_communities.arn,
      aws_dynamodb_table.postponed_event_dynamodb_table.arn,
      aws_dynamodb_table.slack_team_dynamodb_table.arn,

      aws_dynamodb_table.coaching_rejections.arn,
    ]
  }
  statement {
    actions = ["dynamodb:Query"]
    resources = [
      "${aws_dynamodb_table.ad_hoc_holidays.arn}/index/*",
      "${aws_dynamodb_table.adaptive_dialog_aliases.arn}/index/*",
      "${aws_dynamodb_table.adaptive_dialog_content.arn}/index/*",
      "${aws_dynamodb_table.adaptive_user_engagements_dynamo_table.arn}/index/*",
      "${aws_dynamodb_table.adaptive_user_feedback_dynamodb_table.arn}/index/*",
      "${aws_dynamodb_table.adaptive_value_dynamodb_table.arn}/index/*",
      "${aws_dynamodb_table.capability_communities.arn}/index/*",
      "${aws_dynamodb_table.coaching_relationships.arn}/index/*",
      "${aws_dynamodb_table.community_users.arn}/index/*",
      "${aws_dynamodb_table.initiative_communities.arn}/index/*",
      "${aws_dynamodb_table.postponed_event_dynamodb_table.arn}/index/*",
      "${aws_dynamodb_table.user_objective_dynamodb_table.arn}/index/*",
      "${aws_dynamodb_table.strategy_communities.arn}/index/*",
      "${aws_dynamodb_table.strategy_initiatives.arn}/index/*",
      "${aws_dynamodb_table.strategy_objectives.arn}/index/*",
      "${aws_dynamodb_table.strategy_initiatives.arn}/index/*",
      "${aws_dynamodb_table.slack_team_dynamodb_table.arn}/index/*",
    ]
  }
}
