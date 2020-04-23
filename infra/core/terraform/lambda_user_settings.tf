module "user_settings_lambda" {
  source = "../../../terraform-modules/adaptive-lambda-s3"
  s3_bucket = aws_s3_bucket.binary_bucket.bucket
  s3_key = aws_s3_bucket_object.adaptive_zip.key
  source_hash = data.archive_file.adaptive_lambda_zip.output_md5

  client_id     = var.client_id
  handler       = "adaptive"
  function_name_suffix = "user-settings-lambda-go"
  runtime       = var.lambda_runtime
  timeout       = var.lambda_timeout

  // Attach extra policy
  attach_policy = true
  policy        = data.aws_iam_policy_document.adaptive_user_settings_policy.json

  reserved_concurrent_executions = -1

  // Schedule the lambda
  schedule             = true
  schedule_name        = "user_settings_lambda_warmer"
  schedule_description = "User Settings lambda warmer for ${var.client_id}"
  schedule_expression  = "rate(5 minutes)"
  schedule_invoke_json = data.local_file.adaptive_user_settings_lambda_warmer_json.content

  // Add environment variables.
  environment_variables = merge(local.environment_variables, {
    LAMBDA_ROLE   = "user-settings"
    LOG_NAMESPACE = "user-settings"
  })

  tags = local.default_tags
}

data "aws_iam_policy_document" "adaptive_user_settings_policy" {
  statement {
    effect    = "Allow"
    actions   = [
      "dynamodb:DescribeTable",
      "dynamodb:GetItem",
      "dynamodb:UpdateItem",
      "dynamodb:PutItem",
    ]
    resources = [
      aws_dynamodb_table.adaptive_user_engagements_dynamo_table.arn,
    ]
  }
  statement {
    effect    = "Allow"
    actions   = [
      "dynamodb:UpdateItem",
      "dynamodb:GetItem",
    ]
    resources = [
      aws_dynamodb_table.adaptive_users_dynamodb_table.arn,
    ]
  }
  statement {
    effect    = "Allow"
    actions   = [
      "lambda:InvokeFunction",
    ]
    resources = [
      module.user_setup_lambda.function_arn,
      module.user_profile_lambda.function_arn,
      module.user_engagement_scripting_lambda.function_arn,
      module.user_engagement_scheduling_lambda.function_arn,
    ]
  }
  statement {
    effect    = "Allow"
    actions   = [
      "SNS:Publish",
    ]
    resources = [
      aws_sns_topic.platform_notification.arn,
    ]
  }
}

resource "aws_sns_topic_subscription" "adaptive_user_settings_lambda_sns" {
  topic_arn = aws_sns_topic.namespace_payload.arn
  protocol  = "lambda"
  endpoint  = module.user_settings_lambda.function_arn
}

resource "aws_lambda_permission" "adaptive_user_settings_lambda_sns" {
  statement_id  = "AllowExecutionFromSNS"
  action        = "lambda:InvokeFunction"
  function_name = module.user_settings_lambda.function_arn
  principal     = "sns.amazonaws.com"
  source_arn    = aws_sns_topic.namespace_payload.arn
}

data "local_file" "adaptive_user_settings_lambda_warmer_json" {
  filename = "${path.module}/templates/sns_warmup.json"
}

