module "slack_message_processor_lambda" {
  source = "../../../terraform-modules/adaptive-lambda"

  client_id = var.client_id
  filename = data.archive_file.adaptive-lambda-zip.output_path
  source_hash = data.archive_file.adaptive-lambda-zip.output_base64sha256
  function_name = "slack-message-processor-lambda-go"
  handler = "adaptive"
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
  environment_variables = {
    LAMBDA_ROLE="slack-message-processor"

    ADAPTIVE_HELP_PAGE = "https://adaptiveteam.github.io/docs/general/commands"
    NAMESPACE_PAYLOAD_TOPIC_ARN = aws_sns_topic.namespace_payload.arn
    PLATFORM_NOTIFICATION_TOPIC = aws_sns_topic.platform_notification.arn
    USER_COMMUNITIES_TABLE = aws_dynamodb_table.user_communities.name
    COMMUNITY_USERS_TABLE_NAME = aws_dynamodb_table.community_users.name
    COMMUNITY_USERS_USER_COMMUNITY_INDEX = var.dynamo_community_users_user_community_index
    COMMUNITY_USERS_USER_INDEX = var.dynamo_community_users_user_index
    USER_OBJECTIVES_TABLE_NAME = aws_dynamodb_table.user_objective_dynamodb_table.name
    USER_OBJECTIVES_PARTNER_INDEX = var.dynamo_user_objectives_partner_index
    USER_OBJECTIVES_ID_INDEX = var.dynamo_user_objectives_id_index
    USER_OBJECTIVES_TYPE_INDEX = var.dynamo_user_objectives_type_index
    DIALOG_TABLE = aws_dynamodb_table.adaptive_dialog_content.name
    VISION_TABLE_NAME = aws_dynamodb_table.vision.name
    CAPABILITY_COMMUNITIES_TABLE_NAME = aws_dynamodb_table.capability_communities.name
    CAPABILITY_COMMUNITIES_PLATFORM_INDEX = var.dynamo_capability_communities_platform_index
    INITIATIVE_COMMUNITIES_TABLE_NAME = aws_dynamodb_table.initiative_communities.name
    INITIATIVE_COMMUNITIES_PLATFORM_INDEX = var.dynamo_strategy_initiative_communities_platform_index
    STRATEGY_COMMUNITIES_TABLE_NAME = aws_dynamodb_table.strategy_communities.name
    STRATEGY_COMMUNITIES_PLATFORM_CHANNEL_CREATED_INDEX = var.dynamo_strategy_communities_platform_channel_created_index
    LOG_NAMESPACE = "slack-message-processor"
    // ADM
    USER_OBJECTIVES_TABLE = aws_dynamodb_table.user_objective_dynamodb_table.name
    USER_OBJECTIVES_USER_ID_INDEX = var.dynamo_user_objectives_user_index
    USER_OBJECTIVES_PROGRESS_TABLE = aws_dynamodb_table.user_objectives_progress.name
    USER_OBJECTIVES_PROGRESS_ID_INDEX = var.dynamo_user_objectives_progress_index
    CLIENT_ID = var.client_id
    ADAPTIVE_VALUES_TABLE = aws_dynamodb_table.adaptive_value_dynamodb_table.name
    USER_ENGAGEMENT_SCRIPTING_LAMBDA_NAME = module.user_engagement_scripting_lambda.function_name
    USER_ENGAGEMENTS_TABLE_NAME = aws_dynamodb_table.adaptive_user_engagements_dynamo_table.name
    USER_ANSWERED_INDEX = var.user_engagement_answered_dynamo_index
    STRATEGY_INITIATIVES_TABLE_NAME = aws_dynamodb_table.strategy_initiatives.name
    STRATEGY_INITIATIVES_PLATFORM_INDEX = var.dynamo_strategy_initiatives_platform_index
    STRATEGY_OBJECTIVES_TABLE_NAME = aws_dynamodb_table.strategy_objectives.name
    STRATEGY_OBJECTIVES_PLATFORM_INDEX = var.dynamo_strategy_objectives_platform_index
    STRATEGY_OBJECTIVES_CAPABILITY_COMMUNITY_INDEX = var.dynamo_strategy_objectives_capability_community_index
    STRATEGY_INITIATIVES_INITIATIVE_COMMUNITY_ID_INDEX = var.dynamo_strategy_initiatives_initiative_community_index

    USERS_TABLE_NAME = aws_dynamodb_table.adaptive_users_dynamodb_table.name
    USER_PROFILE_LAMBDA_NAME = module.user_profile_lambda.function_name

    // for adaptive-utils-go
    USERS_PLATFORM_INDEX = var.dynamo_users_platform_index
    COMMUNITY_USERS_COMMUNITY_INDEX = var.dynamo_community_users_community_index
    COMMUNITY_USERS_COMMUNITY_INDEX = var.dynamo_community_users_community_index
    COACHING_RELATIONSHIPS_TABLE_NAME = aws_dynamodb_table.coaching_relationships.name
    COACHING_RELATIONSHIPS_COACHEE_QUARTER_YEAR_INDEX = var.dynamo_coaching_relationship_coachee_index
    COACHING_RELATIONSHIPS_COACH_QUARTER_YEAR_INDEX = var.dynamo_coaching_relationship_coach_index
    COACHING_RELATIONSHIPS_QUARTER_YEAR_INDEX = var.dynamo_coaching_relationship_quarter_year_index
    ADAPTIVE_COMMUNITIES_TABLE = aws_dynamodb_table.user_communities.name
    USER_FEEDBACK_TABLE_NAME = aws_dynamodb_table.adaptive_user_feedback_dynamodb_table.name
    USER_FEEDBACK_SOURCE_QUARTER_YEAR_INDEX = var.feedback_source_quarter_year_index
    FEEDBACK_REPORTING_LAMBDA_NAME = local.reporting_lambda_name
    HOLIDAYS_AD_HOC_TABLE = aws_dynamodb_table.ad_hoc_holidays.name
    HOLIDAYS_PLATFORM_DATE_INDEX = var.dynamo_holidays_date_index
    STRATEGY_INITIATIVES_TABLE = aws_dynamodb_table.strategy_initiatives.name
    STRATEGY_INITIATIVES_PLATFORM_INDEX = var.dynamo_strategy_initiatives_platform_index
    STRATEGY_OBJECTIVES_TABLE = aws_dynamodb_table.strategy_objectives.name
    STRATEGY_OBJECTIVES_PLATFORM_INDEX = var.dynamo_strategy_objectives_platform_index
    VISION_TABLE_NAME = aws_dynamodb_table.vision.name
    SLACK_MESSAGE_PROCESSOR_SUFFIX = local.slack_message_processor_suffix

    # Reporting
    REPORTS_BUCKET_NAME = aws_s3_bucket.adaptive-feedback-reports-bucket.bucket

    RDS_HOST     = var.RDS_HOST
		RDS_USER     = var.RDS_USER
		RDS_PASSWORD = var.RDS_PASSWORD
		RDS_PORT     = var.RDS_PORT
		RDS_DB_NAME  = var.RDS_DB_NAME


    CLIENT_CONFIG_TABLE_NAME = aws_dynamodb_table.client_config_dynamodb_table.name
  }

  // Attach extra policy
  attach_policy = true
  policy = data.aws_iam_policy_document.slack_message_processor_dynamo_write_policy.json

  tags = local.default_tags
}

data "aws_iam_policy_document" "slack_message_processor_dynamo_write_policy" {
  statement {
    effect = "Allow"
    actions = [
      "lambda:InvokeFunction",
    ]
    resources = [
      module.user_profile_lambda.function_arn,
    ]
  }

  statement {
    effect = "Allow"
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
    effect = "Allow"
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
    effect = "Allow"
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
    effect = "Allow"
    actions = [
      "SNS:Publish",
    ]
    resources = [
      aws_sns_topic.namespace_payload.arn,
      aws_sns_topic.platform_notification.arn,
    ]
  }

  statement {
    effect = "Allow"
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
  alarm_name = "slack-message-processor-errors"
  alarm_action_arns = [
    aws_sns_topic.errors.arn,
  ]
  period = "60"
  threshold = "1"
  statistic = "SampleCount"
}

