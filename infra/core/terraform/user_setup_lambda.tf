data "archive_file" "user-setup-lambda-zip" {
  type        = "zip"
  source_file = "../../../bin/user-setup-lambda-go"
  output_path = "lambdas/user-setup-lambda-go.zip"
}

module "user_setup_lambda" {
  source = "../../../terraform-modules/adaptive-lambda"

  client_id     = var.client_id
  filename      = data.archive_file.user-setup-lambda-zip.output_path
  source_hash   = data.archive_file.user-setup-lambda-zip.output_base64sha256
  function_name = "user-setup-lambda-go"
  handler       = "user-setup-lambda-go"
  runtime       = var.lambda_runtime
  timeout       = var.lambda_timeout

  reserved_concurrent_executions = -1

  environment_variables = {
    CLIENT_ID                   = var.client_id
    PLATFORM_NOTIFICATION_TOPIC = aws_sns_topic.platform_notification.arn
    USER_ENGAGEMENTS_TABLE_NAME = aws_dynamodb_table.adaptive_user_engagements_dynamo_table.name
    USERS_TABLE_NAME            = aws_dynamodb_table.adaptive_users_dynamodb_table.name
    LOG_NAMESPACE               = "user-setup"
  }

  // Attach extra policy
  attach_policy = true
  policy        = data.aws_iam_policy_document.user_setup_policy.json

  tags = local.default_tags
}

data "aws_iam_policy_document" "user_setup_policy" {
  statement {
    effect    = "Allow"
    actions   = [
      "dynamodb:PutItem",
    ]
    resources = [
      aws_dynamodb_table.adaptive_user_engagements_dynamo_table.arn,
    ]
  }

  statement {
    effect    = "Allow"
    actions   = [
      "dynamodb:GetItem",
      "dynamodb:Query",
    ]
    resources = [
      aws_dynamodb_table.adaptive_users_dynamodb_table.arn,
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

