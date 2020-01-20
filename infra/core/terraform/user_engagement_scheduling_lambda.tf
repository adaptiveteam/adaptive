data "archive_file" "user-engagement-scheduling-lambda-zip" {
  type        = "zip"
  source_file = "../../../bin/user-engagement-scheduling-lambda-go"
  output_path = "lambdas/user-engagement-scheduling-lambda-go.zip"
}

module "user_engagement_scheduling_lambda" {
  source = "../../../terraform-modules/adaptive-lambda"

  client_id     = var.client_id
  filename      = data.archive_file.user-engagement-scheduling-lambda-zip.output_path
  source_hash   = data.archive_file.user-engagement-scheduling-lambda-zip.output_base64sha256
  function_name = "user-engagement-scheduling-lambda-go"
  handler       = "user-engagement-scheduling-lambda-go"
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

  environment_variables = {
    CLIENT_ID                            = var.client_id
    LOG_NAMESPACE                        = "user-engagement-scheduling"
    CLIENT_CONFIG_TABLE_NAME             = aws_dynamodb_table.client_config_dynamodb_table.name
    USERS_TABLE_NAME                     = aws_dynamodb_table.adaptive_users_dynamodb_table.name
    USERS_SCHEDULED_TIME_INDEX           = var.dynamo_users_scheduled_time_index
    USER_ENGAGEMENT_SCRIPTING_LAMBDA_ARN = module.user_engagement_scripting_lambda.function_arn
    USERS_TIMEZONE_OFFSET_INDEX          = var.dynamo_users_timezone_offset_index
    COMMUNITY_USERS_TABLE_NAME           = aws_dynamodb_table.community_users.name
    COMMUNITY_USERS_USER_INDEX           = var.dynamo_community_users_user_index
  }
}

data "aws_iam_policy_document" "user_engagement_scheduling_policy" {
  statement {
    effect    = "Allow"
    actions   = [
      "dynamodb:GetItem",
    ]
    resources = [
      aws_dynamodb_table.adaptive_users_dynamodb_table.arn,
    ]
  }
  statement {
    effect    = "Allow"
    actions   = [
      "dynamodb:Scan",
    ]
    resources = [
      aws_dynamodb_table.client_config_dynamodb_table.arn,
    ]
  }
  statement {
    effect    = "Allow"
    actions   = [
      "dynamodb:Query",
    ]
    resources = [
      "${aws_dynamodb_table.adaptive_users_dynamodb_table.arn}/index/${var.dynamo_users_scheduled_time_index}",
      "${aws_dynamodb_table.adaptive_users_dynamodb_table.arn}/index/${var.dynamo_users_timezone_offset_index}",
      "${aws_dynamodb_table.community_users.arn}/index/${var.dynamo_community_users_user_index}",
    ]
  }
  statement {
    effect    = "Allow"
    actions   = [
      "lambda:InvokeFunction",
    ]
    # TF-UPGRADE-TODO: In Terraform v0.10 and earlier, it was sometimes necessary to
    # force an interpolation expression to be interpreted as a list by wrapping it
    # in an extra set of list brackets. That form was supported for compatibilty in
    # v0.11, but is no longer supported in Terraform v0.12.
    #
    # If the expression in the following list itself returns a list, remove the
    # brackets to avoid interpretation as a list of lists. If the expression
    # returns a single list item then leave it as-is and remove this TODO comment.
    resources = [
      module.user_engagement_scripting_lambda.function_arn,
    ]
  }
}

