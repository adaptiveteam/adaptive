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
    LAMBDA_ROLE   = "user-engagement-scipting"
    LOG_NAMESPACE = "user-engagement-scipting"
  })

  // Attach extra policy
  attach_policy = true
  policy = data.aws_iam_policy_document.user_engagement_scripting_policy.json

  tags = local.default_tags
}

data "aws_iam_policy_document" "user_engagement_scripting_policy" {
  statement {
    effect = "Allow"
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
    effect = "Allow"
    actions = [
      "dynamodb:Query",
      "dynamodb:UpdateItem",
    ]
    resources = [
      aws_dynamodb_table.adaptive_user_engagements_dynamo_table.arn,
    ]
  }

  statement {
    effect = "Allow"
    actions = [
      "lambda:InvokeFunction",
    ]
    resources = [
      "arn:aws:lambda:${var.aws_region}:${data.aws_caller_identity.current.account_id}:function:${var.client_id}_${var.user_engagement_scheduler_lambda_prefix}",
    ]
  }

  statement {
    effect = "Allow"
    actions = [
      "dynamodb:GetItem",
    ]
    resources = [
      aws_dynamodb_table.client_config_dynamodb_table.arn,
    ]
  }

  statement {
    effect = "Allow"
    actions = [
      "SNS:Publish",
    ]
    resources = [
      aws_sns_topic.platform_notification.arn,
    ]
  }
}

