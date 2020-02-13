module "user_engagement_lambda" {
  source = "../../../terraform-modules/adaptive-lambda"

  client_id = var.client_id
  filename = data.archive_file.adaptive-lambda-zip.output_path
  source_hash = data.archive_file.adaptive-lambda-zip.output_base64sha256
  handler = "adaptive"
  function_name_suffix = "user-engagement-lambda-go"
  runtime = var.lambda_runtime
  timeout = var.lambda_timeout

  // Attach extra policy
  attach_policy = true
  policy = data.aws_iam_policy_document.user_engagement_policy.json

  reserved_concurrent_executions = -1

  tags = local.default_tags

  // Add environment variables.
  environment_variables = merge(local.environment_variables, {
    LAMBDA_ROLE   = "user-engagement"
    LOG_NAMESPACE = "user-engagement"
    DEBUG = "0"
    DEBUG_USER = "UE48A5TC0"
  })

}

data "aws_iam_policy_document" "user_engagement_policy" {
  statement {
    effect = "Allow"
    actions = [
      "dynamodb:GetRecords",
      "dynamodb:GetShardIterator",
      "dynamodb:DescribeStream",
      "dynamodb:ListStreams",
    ]
    resources = [
      aws_dynamodb_table.adaptive_user_engagements_dynamo_table.stream_arn,
    ]
  }

  statement {
    effect = "Allow"
    actions = [
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
      module.user_profile_lambda.function_arn,
    ]
  }

  statement {
    effect = "Allow"
    actions = [
      "dynamodb:Query",
    ]
    resources = [
      "${aws_dynamodb_table.community_users.arn}/index/${var.dynamo_community_users_user_index}",
      "${aws_dynamodb_table.community_users.arn}/index/${var.dynamo_community_users_channel_index}",
    ]
  }
}

resource "aws_lambda_event_source_mapping" "user_enagagement_lambda_source_mapping" {
  event_source_arn = aws_dynamodb_table.adaptive_user_engagements_dynamo_table.stream_arn
  function_name = module.user_engagement_lambda.function_arn
  starting_position = "LATEST"
}

