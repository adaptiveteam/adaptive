locals {
  user_engagement_scheduler_lambda_function_name_suffix = "user-engagement-scheduler-lambda-go"
  user_engagement_scheduler_lambda_function_name = "${var.client_id}_${local.user_engagement_scheduler_lambda_function_name_suffix}"
}
module "user_engagement_scheduler_lambda" {
  source = "../../../terraform-modules/adaptive-lambda"

  client_id = var.client_id
  filename = data.archive_file.adaptive-lambda-zip.output_path
  source_hash = data.archive_file.adaptive-lambda-zip.output_base64sha256
  handler = "adaptive"
  function_name_suffix = local.user_engagement_scheduler_lambda_function_name_suffix
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

  // Add environment variables.
  environment_variables = merge(local.environment_variables, {
    LAMBDA_ROLE   = "user-engagement-scheduler"
    LOG_NAMESPACE = "user-engagement-scheduler"
  })

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
    actions = [
      "dynamodb:GetItem",]
    resources = [
      aws_dynamodb_table.vision.arn,
      aws_dynamodb_table.user_communities.arn,
      aws_dynamodb_table.client_config_dynamodb_table.arn,
      aws_dynamodb_table.adaptive_users_dynamodb_table.arn,
      aws_dynamodb_table.strategy_communities.arn,
      aws_dynamodb_table.user_objective_dynamodb_table.arn,
    ]
  }

  statement {
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
      "${aws_dynamodb_table.strategy_initiatives.arn}/index/*",

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
