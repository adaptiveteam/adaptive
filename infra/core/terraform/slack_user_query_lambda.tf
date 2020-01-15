data "archive_file" "slack-user-query-lambda-zip" {
  type = "zip"
  source_file = "../../../bin/slack-user-query-lambda-go"
  output_path = "lambdas/slack-user-query-lambda-go.zip"
}

module "slack_user_query_lambda" {
  source = "../../../terraform-modules/adaptive-lambda"

  client_id = var.client_id
  filename = data.archive_file.slack-user-query-lambda-zip.output_path
  source_hash = data.archive_file.slack-user-query-lambda-zip.output_base64sha256
  function_name = "slack-user-query-lambda-go"
  handler = "slack-user-query-lambda-go"
  runtime = var.lambda_runtime
  timeout = 900
  memory_size = var.multi_core_memory_size

  reserved_concurrent_executions = -1

  environment_variables = {
    USER_SETUP_LAMBDA_NAME = module.user_setup_lambda.function_name
    USER_TABLE_NAME = aws_dynamodb_table.adaptive_users_dynamodb_table.name
    CLIENT_CONFIG_TABLE_NAME = aws_dynamodb_table.client_config_dynamodb_table.name
    CLIENT_ID = var.client_id
    LOG_NAMESPACE = "slack-user-query"
    USER_COMMUNITY_TABLE_NAME = aws_dynamodb_table.user_communities.name
    USER_COMMUNITY_PLATFORM_INDEX = var.user_community_platform_dynamo_index
    COMMUNITY_USERS_TABLE_NAME = aws_dynamodb_table.community_users.name
    COMMUNITY_USERS_COMMUNITY_INDEX = var.dynamo_community_users_community_index
  }

  // Attach extra policy
  attach_policy = true
  policy = data.aws_iam_policy_document.slack_user_query_policy.json

  tags = local.default_tags
}

data "aws_iam_policy_document" "slack_user_query_policy" {
  statement {
    effect = "Allow"
    actions = [
      "dynamodb:PutItem",
      "dynamodb:GetItem",
      "dynamodb:DeleteItem",
    ]
    resources = [
      aws_dynamodb_table.adaptive_users_dynamodb_table.arn,
    ]
  }

  statement {
    effect = "Allow"
    actions = [
      "lambda:InvokeFunction",
    ]
    resources = [
      module.user_setup_lambda.function_arn,
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
      "dynamodb:Query",
    ]
    resources = [
      "${aws_dynamodb_table.user_communities.arn}/index/${var.user_community_platform_dynamo_index}",
      "${aws_dynamodb_table.community_users.arn}/index/${var.dynamo_community_users_community_index}",
      "${aws_dynamodb_table.adaptive_users_dynamodb_table.arn}/index/${var.dynamo_users_platform_index}",
    ]
  }
}

