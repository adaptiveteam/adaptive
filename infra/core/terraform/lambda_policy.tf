// Policy for the main lambda role
resource "aws_iam_policy" "additional" {
  name   = "${var.client_id}_additional_policy_attachment_for_lamba_role"
  policy = data.aws_iam_policy_document.slack_message_processor_dynamo_write_policy.json
}

resource "aws_iam_policy_attachment" "additional" {
  name       = "${var.client_id}_additional_policy_attachment_for_lamba_role"
  roles      = [aws_iam_role.lambda.name]
  policy_arn = aws_iam_policy.additional.arn
}

data "aws_iam_policy_document" "slack_message_processor_dynamo_write_policy" {
  statement {
    actions = [
      "lambda:InvokeFunction",
    ]
    resources = [
      module.user_profile_lambda.function_arn,
    ]
  }

  statement {
    actions = [
      "dynamodb:GetItem",
    ]
    resources = [
      aws_dynamodb_table.client_config_dynamodb_table.arn,
      aws_dynamodb_table.adaptive_users_dynamodb_table.arn,
      aws_dynamodb_table.coaching_relationships.arn,
      aws_dynamodb_table.community_users.arn,
      aws_dynamodb_table.vision.arn,
      aws_dynamodb_table.strategy_initiatives.arn,
      aws_dynamodb_table.adaptive_users_dynamodb_table.arn,
      aws_dynamodb_table.user_communities.arn,
      aws_dynamodb_table.strategy_communities.arn,
      aws_dynamodb_table.adaptive_dialog_content.arn,
      aws_dynamodb_table.adaptive_dialog_aliases.arn,
      aws_dynamodb_table.user_objective_dynamodb_table.arn,
    ]
  }

  statement {
    actions = [
      "dynamodb:Query",
    ]
    resources = [
      aws_dynamodb_table.user_objectives_progress.arn,
      aws_dynamodb_table.client_config_dynamodb_table.arn,
      "${aws_dynamodb_table.adaptive_user_engagements_dynamo_table.arn}/index/${var.user_engagement_answered_dynamo_index}",
      "${aws_dynamodb_table.adaptive_value_dynamodb_table.arn}/index/${var.dynamo_adaptive_values_platform_id_index}",
      "${aws_dynamodb_table.community_users.arn}/index/${var.dynamo_community_users_user_community_index}",
      "${aws_dynamodb_table.community_users.arn}/index/${var.dynamo_community_users_user_index}",
      "${aws_dynamodb_table.adaptive_dialog_content.arn}/index/${var.dynamo_dialog_content_contect_subject_index}",
      "${aws_dynamodb_table.user_objective_dynamodb_table.arn}/index/${var.dynamo_user_objectives_partner_index}",
      "${aws_dynamodb_table.user_objective_dynamodb_table.arn}/index/${var.dynamo_user_objectives_user_index}",
      "${aws_dynamodb_table.user_objective_dynamodb_table.arn}/index/${var.dynamo_user_objectives_type_index}",
      "${aws_dynamodb_table.user_objective_dynamodb_table.arn}/index/${var.dynamo_user_objectives_id_index}",
      "${aws_dynamodb_table.strategy_communities.arn}/index/${var.dynamo_strategy_communities_platform_channel_created_index}",
      "${aws_dynamodb_table.capability_communities.arn}/index/${var.dynamo_capability_communities_platform_index}",
      "${aws_dynamodb_table.initiative_communities.arn}/index/${var.dynamo_strategy_initiative_communities_platform_index}",
      "${aws_dynamodb_table.strategy_initiatives.arn}/index/${var.dynamo_strategy_initiatives_platform_index}",
      "${aws_dynamodb_table.strategy_objectives.arn}/index/${var.dynamo_strategy_objectives_platform_index}",
      "${aws_dynamodb_table.strategy_objectives.arn}/index/${var.dynamo_strategy_objectives_capability_community_index}",
      "${aws_dynamodb_table.strategy_initiatives.arn}/index/${var.dynamo_strategy_initiatives_initiative_community_index}",
      "${aws_dynamodb_table.adaptive_user_feedback_dynamodb_table.arn}/index/${var.feedback_source_quarter_year_index}",
      "${aws_dynamodb_table.coaching_relationships.arn}/index/*",
      "${aws_dynamodb_table.ad_hoc_holidays.arn}/index/*",
      "${aws_dynamodb_table.adaptive_dialog_content.arn}/index/*",
      "${aws_dynamodb_table.adaptive_dialog_aliases.arn}/index/*",
    ]
  }

  statement {
    actions = [
      "lambda:InvokeFunction",
    ]
    resources = [
      "arn:aws:lambda:${var.aws_region}:${data.aws_caller_identity.current.account_id}:function:${var.client_id}_strategy-${local.slack_message_processor_suffix}",
      "arn:aws:lambda:${var.aws_region}:${data.aws_caller_identity.current.account_id}:function:${var.client_id}_feedback-${local.slack_message_processor_suffix}",
      module.user_profile_lambda.function_arn,
      module.user_engagement_scripting_lambda.function_arn,
    ]
  }

  statement {
    actions = [
      "SNS:Publish",
    ]
    resources = [
      aws_sns_topic.namespace_payload.arn,
      aws_sns_topic.platform_notification.arn,
    ]
  }

  statement {
    actions = [
      "s3:GetObject",
      "s3:GetObjectAcl",]
    resources = [
      "${aws_s3_bucket.adaptive-feedback-reports-bucket.arn}/*",
    ]
  }
}
