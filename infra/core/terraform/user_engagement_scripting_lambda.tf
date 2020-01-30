locals {
  user_engagement_scripting_lambda_function_name_suffix = "user-engagement-scripting-lambda-go"
  user_engagement_scripting_lambda_function_name = "${var.client_id}_${local.user_engagement_scripting_lambda_function_name_suffix}"
}

module "user_engagement_scripting_lambda" {
  source = "../../../terraform-modules/adaptive-lambda"

  client_id = var.client_id
  filename      = data.archive_file.adaptive-lambda-zip.output_path
  source_hash   = data.archive_file.adaptive-lambda-zip.output_base64sha256
  handler       = "adaptive"
  function_name_suffix = local.user_engagement_scripting_lambda_function_name_suffix
  runtime = var.lambda_runtime
  timeout = var.lambda_timeout
  memory_size = var.multi_core_memory_size

  reserved_concurrent_executions = -1

  // Add environment variables.
  environment_variables = merge(local.environment_variables, {
    LAMBDA_ROLE   = "user-engagement-scripting"
    LOG_NAMESPACE = "user-engagement-scripting"
  })

  // Attach extra policy
  attach_policy = true
  policy = data.aws_iam_policy_document.user_engagement_scripting_policy.json

  tags = local.default_tags
}

data "aws_iam_policy_document" "user_engagement_scripting_policy" {
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
    actions = [
      "dynamodb:DescribeTable",
      "dynamodb:GetItem",
      "dynamodb:Query",
    ]
    resources = [
      "${aws_dynamodb_table.adaptive_user_engagements_dynamo_table.arn}/index/${var.user_engagement_answered_dynamo_index}",
    ]
  }

  statement {
    resources = [aws_dynamodb_table.adaptive_user_engagements_dynamo_table.arn]
    actions = [
      "dynamodb:Query",
      "dynamodb:UpdateItem",
    ]
  }
  statement {
    resources = [module.user_engagement_scheduler_lambda.function_arn]
    actions   = ["lambda:InvokeFunction"]
  }
  statement {
    resources = [aws_sns_topic.platform_notification.arn]
    actions = ["SNS:Publish"]
  }
}

