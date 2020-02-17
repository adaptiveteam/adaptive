locals {
  slack_message_processor_lambda_function_name_suffix = "slack-message-processor-lambda-go"
  slack_message_processor_lambda_function_name = "${var.client_id}_${local.slack_message_processor_lambda_function_name_suffix}"
}
module "slack_message_processor_lambda" {
  source = "../../../terraform-modules/adaptive-lambda"

  client_id = var.client_id
  filename = data.archive_file.adaptive-lambda-zip.output_path
  source_hash = data.archive_file.adaptive-lambda-zip.output_base64sha256
  handler = "adaptive"
  function_name_suffix = local.slack_message_processor_lambda_function_name_suffix
  runtime = var.lambda_runtime
  timeout = var.lambda_timeout
  memory_size = var.multi_core_memory_size

  reserved_concurrent_executions = -1

  // Schedule the lambda
  schedule = true
  schedule_name = "slack_message_processor_lambda-warmer"
  schedule_description = "Slack Message Processor lambda warmer for ${var.client_id}"
  schedule_expression = "rate(5 minutes)"
  schedule_invoke_json = data.local_file.slack_message_processor_lambda_warmer_json.content

  // Add environment variables.
  environment_variables = merge(local.environment_variables, {
    LAMBDA_ROLE   = "slack-message-processor"
    LOG_NAMESPACE = "slack-message-processor"
  })

  // Attach extra policy
  attach_policy = true
  policy = data.aws_iam_policy_document.slack_message_processor_dynamo_write_policy.json

  tags = local.default_tags
}

data "aws_iam_policy_document" "slack_message_processor_dynamo_write_policy" {
  statement {
    actions = ["lambda:InvokeFunction"]
    resources = [
      module.user_profile_lambda.function_arn,
      "arn:aws:lambda:${var.aws_region}:${data.aws_caller_identity.current.account_id}:function:${var.client_id}_strategy-${local.slack_message_processor_suffix}",
      "arn:aws:lambda:${var.aws_region}:${data.aws_caller_identity.current.account_id}:function:${var.client_id}_feedback-${local.slack_message_processor_suffix}",
      module.user_profile_lambda.function_arn,
      module.user_engagement_scripting_lambda.function_arn,
    ]
  }

  # statement {
  #   actions = [
  #     "dynamodb:GetItem",
  #   ]
  #   resources = [
  #     aws_dynamodb_table.client_config_dynamodb_table.arn,
  #     aws_dynamodb_table.adaptive_users_dynamodb_table.arn,
  #     aws_dynamodb_table.coaching_relationships.arn,
  #     aws_dynamodb_table.community_users.arn,
  #     aws_dynamodb_table.vision.arn,
  #     aws_dynamodb_table.strategy_initiatives.arn,
  #     aws_dynamodb_table.adaptive_users_dynamodb_table.arn,
  #     aws_dynamodb_table.user_communities.arn,
  #     aws_dynamodb_table.strategy_communities.arn,
  #     aws_dynamodb_table.adaptive_dialog_content.arn,
  #     aws_dynamodb_table.adaptive_dialog_aliases.arn,
  #     aws_dynamodb_table.user_objective_dynamodb_table.arn,
  #   ]
  # }

  # statement {
  #   actions = ["dynamodb:Query"]
  #   resources = [
  #     aws_dynamodb_table.user_objectives_progress.arn,
  #     aws_dynamodb_table.client_config_dynamodb_table.arn,
  #     "${aws_dynamodb_table.adaptive_user_engagements_dynamo_table.arn}/index/*",
  #     "${aws_dynamodb_table.adaptive_value_dynamodb_table.arn}/index/*",
  #     "${aws_dynamodb_table.community_users.arn}/index/*",
  #     "${aws_dynamodb_table.adaptive_dialog_content.arn}/index/*",
  #     "${aws_dynamodb_table.user_objective_dynamodb_table.arn}/index/*",
  #     "${aws_dynamodb_table.strategy_communities.arn}/index/*",
  #     "${aws_dynamodb_table.capability_communities.arn}/index/*",
  #     "${aws_dynamodb_table.initiative_communities.arn}/index/*",
  #     "${aws_dynamodb_table.strategy_initiatives.arn}/index/*",
  #     "${aws_dynamodb_table.strategy_objectives.arn}/index/*",
  #     "${aws_dynamodb_table.strategy_initiatives.arn}/index/${var.dynamo_strategy_initiatives_initiative_community_index}",
  #     "${aws_dynamodb_table.adaptive_user_feedback_dynamodb_table.arn}/index/${var.feedback_source_quarter_year_index}",
  #     "${aws_dynamodb_table.coaching_relationships.arn}/index/*",
  #     "${aws_dynamodb_table.ad_hoc_holidays.arn}/index/*",
  #     "${aws_dynamodb_table.adaptive_dialog_content.arn}/index/*",
  #     "${aws_dynamodb_table.adaptive_dialog_aliases.arn}/index/*",
  #     "${aws_dynamodb_table.postponed_event_dynamodb_table.arn}/index/*",
  #   ]
  # }

  statement {
    actions = ["SNS:Publish"]
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

data "local_file" "slack_message_processor_lambda_warmer_json" {
  filename = "${path.module}/templates/api_slack_warmup.json"
}

module "slack_message_processor_error_alarm" {
  // TODO: Pin to released version once this repo is released
  source = "github.com/dwp/terraform-aws-metric-filter-alarm?ref=master"
  log_group_name = module.slack_message_processor_lambda.log_group_name
  metric_namespace = "${var.client_id}-AWS/Lambda"
  pattern = "ERROR"
  alarm_name = "${var.client_id}-slack-message-processor-errors"
  alarm_action_arns = [
    aws_sns_topic.errors.arn,
  ]
  period = "60"
  threshold = "1"
  statistic = "SampleCount"
}

resource "aws_iam_role_policy_attachment" "slack_message_processor_lambda_holidays_additional_policy_attachment" {
  role       = module.slack_message_processor_lambda.role_name
  policy_arn = aws_iam_policy.holidays_additional_policy.arn
}

resource "aws_iam_role_policy_attachment" "slack_message_processor_lambda_read_all_tables" {
  role       = module.slack_message_processor_lambda.role_name
  policy_arn = aws_iam_policy.read_all_tables.arn
}

resource "aws_iam_role_policy_attachment" "slack_message_processor_lambda_competencies_additional_policy_attachment" {
  role       = module.slack_message_processor_lambda.role_name
  policy_arn = aws_iam_policy.competencies_additional_policy.arn
}
