data "archive_file" "user_engagement_scheduler_lambda_zip" {
  type = "zip"
  source_file = "../../../bin/user-engagement-scheduler-lambda-go"
  output_path = "lambdas/user-engagement-scheduler-lambda-go.zip"
}

module "user_engagement_scheduler_lambda" {
  source = "../../../terraform-modules/adaptive-lambda"

  client_id = var.client_id
  filename = data.archive_file.user_engagement_scheduler_lambda_zip.output_path
  source_hash = data.archive_file.user_engagement_scheduler_lambda_zip.output_base64sha256
  function_name = "user-engagement-scheduler-lambda-go"
  handler = "user-engagement-scheduler-lambda-go"
  runtime = var.lambda_runtime
  timeout = var.lambda_timeout
  memory_size = var.multi_core_memory_size

  reserved_concurrent_executions = -1

  // Schedule the lambda
  schedule = true
  schedule_name = "user_engagement_scheduler_lambda_scheduled_run"
  schedule_description = "User Engagement Scheuduler Lambda Scheduled Run for ${var.client_id}"
  schedule_expression = "cron(0 10 ? * MON-FRI *)"
  schedule_invoke_json = "{}"

  environment_variables = {
    USER_ENGAGEMENTS_TABLE_NAME = aws_dynamodb_table.adaptive_user_engagements_dynamo_table.name
    USERS_TABLE_NAME = aws_dynamodb_table.adaptive_users_dynamodb_table.name
    USERS_PLATFORM_INDEX = var.dynamo_users_platform_index
    COMMUNITY_USERS_TABLE_NAME = aws_dynamodb_table.community_users.name
    COMMUNITY_USERS_COMMUNITY_INDEX = var.dynamo_community_users_community_index
    COACHING_RELATIONSHIPS_TABLE_NAME = aws_dynamodb_table.coaching_relationships.name
    COACHING_RELATIONSHIPS_COACHEE_QUARTER_YEAR_INDEX = var.dynamo_coaching_relationship_coachee_index
    COACHING_RELATIONSHIPS_COACH_QUARTER_YEAR_INDEX = var.dynamo_coaching_relationship_coach_index

    COACHING_RELATIONSHIPS_QUARTER_YEAR_INDEX = var.dynamo_coaching_relationship_quarter_year_index
    USER_OBJECTIVES_TABLE = aws_dynamodb_table.user_objective_dynamodb_table.name
    USER_OBJECTIVES_USER_ID_INDEX = var.dynamo_user_objectives_user_index
    USER_OBJECTIVES_PARTNER_INDEX = var.dynamo_user_objectives_partner_index
    USER_FEEDBACK_TABLE_NAME = aws_dynamodb_table.adaptive_user_feedback_dynamodb_table.name
    USER_FEEDBACK_SOURCE_QUARTER_YEAR_INDEX = var.feedback_source_quarter_year_index
    USER_PROFILE_LAMBDA_NAME = module.user_profile_lambda.function_name
    FEEDBACK_REPORTING_LAMBDA_NAME = local.reporting_lambda_name
    PLATFORM_NOTIFICATION_TOPIC = local.platform_notification_topic_arn
    ADAPTIVE_COMMUNITIES_TABLE = aws_dynamodb_table.user_communities.name
    USER_OBJECTIVES_TABLE_NAME = aws_dynamodb_table.user_objective_dynamodb_table.name
    USER_OBJECTIVES_TYPE_INDEX = var.dynamo_user_objectives_type_index
    USER_OBJECTIVES_PROGRESS_TABLE = aws_dynamodb_table.user_objectives_progress.name
    USER_OBJECTIVES_PROGRESS_ID_INDEX = var.dynamo_user_objectives_progress_index
    COMMUNITY_USERS_USER_COMMUNITY_INDEX = var.dynamo_community_users_user_community_index
    COMMUNITY_USERS_USER_INDEX = var.dynamo_community_users_user_index
    USER_ANSWERED_INDEX = var.user_engagement_answered_dynamo_index
    STRATEGY_INITIATIVES_TABLE_NAME = aws_dynamodb_table.strategy_initiatives.name
    STRATEGY_OBJECTIVES_TABLE_NAME = aws_dynamodb_table.strategy_objectives.name
    CAPABILITY_COMMUNITIES_TABLE_NAME = aws_dynamodb_table.capability_communities.name
    CAPABILITY_COMMUNITIES_PLATFORM_INDEX = var.dynamo_capability_communities_platform_index
    INITIATIVE_COMMUNITIES_TABLE_NAME = aws_dynamodb_table.initiative_communities.name
    INITIATIVE_COMMUNITIES_PLATFORM_INDEX = var.dynamo_strategy_initiative_communities_platform_index
    STRATEGY_OBJECTIVES_CAPABILITY_COMMUNITY_INDEX = var.dynamo_strategy_objectives_capability_community_index
    STRATEGY_INITIATIVES_INITIATIVE_COMMUNITY_ID_INDEX = var.dynamo_strategy_initiatives_initiative_community_index

    STRATEGY_INITIATIVES_TABLE = aws_dynamodb_table.strategy_initiatives.name
    STRATEGY_INITIATIVES_PLATFORM_INDEX = var.dynamo_strategy_initiatives_platform_index
    STRATEGY_OBJECTIVES_TABLE = aws_dynamodb_table.strategy_objectives.name
    STRATEGY_OBJECTIVES_PLATFORM_INDEX = var.dynamo_strategy_objectives_platform_index
    VISION_TABLE_NAME = aws_dynamodb_table.vision.name

    HOLIDAYS_AD_HOC_TABLE = aws_dynamodb_table.ad_hoc_holidays.name
    HOLIDAYS_PLATFORM_DATE_INDEX = var.dynamo_holidays_date_index

    // for schedules
    CLIENT_ID = var.client_id
    USER_OBJECTIVES_ID_INDEX = var.dynamo_user_objectives_id_index
    STRATEGY_COMMUNITIES_TABLE_NAME = aws_dynamodb_table.strategy_communities.name

    LOG_NAMESPACE = "user-engagement-scheduler"

    # Reporting
    REPORTS_BUCKET_NAME = aws_s3_bucket.adaptive-feedback-reports-bucket.bucket
  }

  // Attach extra policy
  attach_policy = true
  policy = data.aws_iam_policy_document.user_engagement_scheduler_policy.json

  tags = local.default_tags
}

data "aws_iam_policy_document" "user_engagement_scheduler_policy" {
  statement {
    effect = "Allow"
    actions = [
      "dynamodb:PutItem",]
    resources = [
      aws_dynamodb_table.adaptive_user_engagements_dynamo_table.arn,
    ]
  }

  statement {
    effect = "Allow"
    actions = [
      "dynamodb:GetItem",]
    resources = [
      aws_dynamodb_table.vision.arn,
      aws_dynamodb_table.user_communities.arn,
      aws_dynamodb_table.adaptive_users_dynamodb_table.arn,
      aws_dynamodb_table.strategy_communities.arn,
    ]
  }

  statement {
    effect = "Allow"
    actions = [
      "dynamodb:Query",]
    resources = [
      "${aws_dynamodb_table.adaptive_users_dynamodb_table.arn}/index/${var.dynamo_users_platform_index}",
      "${aws_dynamodb_table.adaptive_user_engagements_dynamo_table.arn}/index/${var.user_engagement_answered_dynamo_index}",
      "${aws_dynamodb_table.community_users.arn}/index/${var.dynamo_community_users_community_index}",
      "${aws_dynamodb_table.community_users.arn}/index/${var.dynamo_community_users_user_community_index}",
      "${aws_dynamodb_table.coaching_relationships.arn}/index/${var.dynamo_coaching_relationship_coachee_index}",
      "${aws_dynamodb_table.coaching_relationships.arn}/index/${var.dynamo_coaching_relationship_coach_index}",
      "${aws_dynamodb_table.coaching_relationships.arn}/index/${var.dynamo_coaching_relationship_quarter_year_index}",
      "${aws_dynamodb_table.user_objective_dynamodb_table.arn}/index/${var.dynamo_user_objectives_user_index}",
      "${aws_dynamodb_table.user_objective_dynamodb_table.arn}/index/${var.dynamo_user_objectives_partner_index}",
      "${aws_dynamodb_table.adaptive_user_feedback_dynamodb_table.arn}/index/${var.feedback_source_quarter_year_index}",
      "${aws_dynamodb_table.user_objectives_progress.arn}/index/${var.dynamo_user_objectives_progress_index}",
      "${aws_dynamodb_table.ad_hoc_holidays.arn}/index/${var.dynamo_holidays_date_index}",
      "${aws_dynamodb_table.strategy_objectives.arn}/index/${var.dynamo_strategy_objectives_platform_index}",
      "${aws_dynamodb_table.strategy_initiatives.arn}/index/${var.dynamo_strategy_initiatives_platform_index}",
      "${aws_dynamodb_table.strategy_initiatives.arn}/index/${var.dynamo_strategy_initiatives_initiative_community_index}",

      // New ADM perms
      "${aws_dynamodb_table.user_objective_dynamodb_table.arn}/index/${var.dynamo_user_objectives_type_index}",
      "${aws_dynamodb_table.community_users.arn}/index/${var.dynamo_community_users_user_index}",
      "${aws_dynamodb_table.adaptive_value_dynamodb_table.arn}/index/${var.dynamo_adaptive_values_platform_id_index}",
      "${aws_dynamodb_table.capability_communities.arn}/index/${var.dynamo_capability_communities_platform_index}",
      "${aws_dynamodb_table.initiative_communities.arn}/index/${var.dynamo_strategy_initiative_communities_platform_index}",

      // For schedules
      "${aws_dynamodb_table.user_objective_dynamodb_table.arn}/index/${var.dynamo_user_objectives_id_index}",
      aws_dynamodb_table.user_objectives_progress.arn,
      aws_dynamodb_table.postponed_event_dynamodb_table.arn,
      "${aws_dynamodb_table.postponed_event_dynamodb_table.arn}/index/*"
    ]
  }

  statement {
    effect = "Allow"
    actions = [
      "lambda:InvokeFunction",]
    resources = [
      module.user_profile_lambda.function_arn,
      "arn:aws:lambda:${var.aws_region}:${data.aws_caller_identity.current.account_id}:function:${local.reporting_lambda_name}",
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

  statement {
    effect = "Allow"
    actions = [
      "SNS:Publish",]
    resources = [
      local.platform_notification_topic_arn,]
  }
}
