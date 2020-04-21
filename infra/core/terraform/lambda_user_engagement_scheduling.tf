locals {
  user_engagement_scheduling_lambda_function_name_suffix = "user-engagement-scheduling-lambda-go"
  user_engagement_scheduling_lambda_function_name = "${var.client_id}_${local.user_engagement_scheduling_lambda_function_name_suffix}"
}

module "user_engagement_scheduling_lambda" {
  source = "../../../terraform-modules/adaptive-lambda-s3"
  s3_bucket = aws_s3_bucket.binary_bucket.bucket
  s3_key = aws_s3_bucket_object.adaptive_zip.key
  source_hash = data.archive_file.adaptive_lambda_zip.output_md5

  client_id     = var.client_id
  handler = "adaptive"
  function_name_suffix = local.user_engagement_scheduling_lambda_function_name_suffix
  runtime       = var.lambda_runtime
  timeout       = var.lambda_timeout
  memory_size   = var.multi_core_memory_size

  // Attach extra policy
  attach_policy = true
  policy        = data.aws_iam_policy_document.user_engagement_scheduling_policy.json

  reserved_concurrent_executions = -1

  // Schedule the lambda
  schedule             = true
  schedule_name        = "user_engagement_scheduling_lambda_runner"
  schedule_description = "User Engagement Scheduling Lambda Runner"
  schedule_expression  = "cron(0/15 0-23 ? * MON-FRI *)"
  schedule_invoke_json = "{}"

  tags = local.default_tags

//     USER_ENGAGEMENT_SCRIPTING_LAMBDA_ARN    = local.user_engagement_scripting_lambda_function_arn

  // Add environment variables.
  environment_variables = merge(local.environment_variables, {
    LAMBDA_ROLE   = "user-engagement-scheduling"
    LOG_NAMESPACE = "user-engagement-scheduling"
    USER_ENGAGEMENT_SCRIPTING_LAMBDA_ARN = module.user_engagement_scripting_lambda.function_arn
    USER_ENGAGEMENT_SCHEDULER_LAMBDA_NAME= module.user_engagement_scheduler_lambda.function_name
  })

}

data "aws_iam_policy_document" "user_engagement_scheduling_policy" {
  statement {
    resources = [aws_dynamodb_table.adaptive_users_dynamodb_table.arn, "${aws_dynamodb_table.adaptive_users_dynamodb_table.arn}/index/*"]
    actions   = ["dynamodb:*"]
  }
  statement {
    resources = [aws_dynamodb_table.user_objective_dynamodb_table.arn]
    actions   = ["dynamodb:*"]
  }
  statement {
    resources = [aws_dynamodb_table.strategy_objectives.arn]
    actions   = ["dynamodb:*"]
  }
  statement {
    resources = [aws_dynamodb_table.strategy_initiatives.arn]
    actions   = ["dynamodb:*"]
  }
  statement {
    resources = [aws_dynamodb_table.vision.arn]
    actions   = ["dynamodb:*"]
  }
  statement {
    resources = [aws_dynamodb_table.client_config_dynamodb_table.arn]
    actions   = ["dynamodb:Scan", "dynamodb:GetItem"]
  }
  statement {
    resources = ["${aws_dynamodb_table.community_users.arn}/index/*"]
    actions   = ["dynamodb:Query"]
  }
  statement {
    resources = [aws_dynamodb_table.postponed_event_dynamodb_table.arn,"${aws_dynamodb_table.postponed_event_dynamodb_table.arn}/index/*"]
    actions   = ["dynamodb:*"]
  }
  statement {
    resources = [module.user_engagement_scripting_lambda.function_arn,module.user_engagement_scheduler_lambda.function_arn]
    actions   = ["lambda:InvokeFunction"]
  }
}

resource "aws_iam_role_policy_attachment" "user_engagement_scheduling_lambda_read_all_tables" {
  role       = module.user_engagement_scheduling_lambda.role_name
  policy_arn = aws_iam_policy.read_all_tables.arn
}
