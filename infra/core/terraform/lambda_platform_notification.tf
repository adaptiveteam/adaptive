module "adaptive-platform-notification-lambda" {
  source = "../../../terraform-modules/adaptive-lambda"

  client_id     = var.client_id
  filename = data.archive_file.adaptive-lambda-zip.output_path
  source_hash = data.archive_file.adaptive-lambda-zip.output_base64sha256
  handler = "adaptive"
  function_name_suffix = "adaptive-platform-notification-lambda-go"
  runtime       = var.lambda_runtime
  timeout       = var.lambda_timeout
  memory_size   = var.multi_core_memory_size

  reserved_concurrent_executions = -1

  // Attach extra policy
  attach_policy = true
  policy        = data.aws_iam_policy_document.adaptive-platform-notification-policy.json

  tags = local.default_tags

  // Add environment variables.
  environment_variables = merge(local.environment_variables, {
    LAMBDA_ROLE   = "platform-notification"
    LOG_NAMESPACE = "platform-notification"
  })

}

data "aws_iam_policy_document" "adaptive-platform-notification-policy" {
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
      module.user_profile_lambda.function_arn,
    ]
  }
  statement {
    effect    = "Allow"
    actions   = [
      "dynamodb:DescribeTable",
      "dynamodb:Query",
      "dynamodb:GetItem",
    ]
    resources = [
      aws_dynamodb_table.client_config_dynamodb_table.arn,
    ]
  }
  statement {
    effect    = "Allow"
    actions   = [
      "dynamodb:GetItem",
    ]
    resources = [
      aws_dynamodb_table.adaptive_users_dynamodb_table.arn,
    ]
  }
}

resource "aws_sns_topic_subscription" "adaptive_platform_notification_lambda_sns" {
  topic_arn = aws_sns_topic.platform_notification.arn
  protocol  = "lambda"
  endpoint  = module.adaptive-platform-notification-lambda.function_arn
}

resource "aws_lambda_permission" "adaptive_platform_notification_lambda_sns" {
  statement_id  = "AllowExecutionFromSNS"
  action        = "lambda:InvokeFunction"
  function_name = module.adaptive-platform-notification-lambda.function_arn
  principal     = "sns.amazonaws.com"
  source_arn    = aws_sns_topic.platform_notification.arn
}

resource "aws_iam_role_policy_attachment" "platform_notification_lambda_read_all_tables" {
  role       = module.adaptive-platform-notification-lambda.role_name
  policy_arn = aws_iam_policy.read_all_tables.arn
}
