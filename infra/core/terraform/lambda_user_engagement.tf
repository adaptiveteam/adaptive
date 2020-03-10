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
    actions = [
      "dynamodb:UpdateItem",
    ]
    resources = [
      aws_dynamodb_table.adaptive_user_engagements_dynamo_table.arn,
    ]
  }

  statement {
    actions = [
      "lambda:InvokeFunction",
    ]
    resources = [
      module.user_profile_lambda.function_arn,
    ]
  }

  statement {
    actions = [
      "dynamodb:Query",
    ]
    resources = [
      "${aws_dynamodb_table.community_users.arn}/index/*",
    ]
  }
}

resource "aws_iam_role_policy_attachment" "user_engagement_lambda_read_all_tables" {
  role       = module.user_engagement_lambda.role_name
  policy_arn = aws_iam_policy.read_all_tables.arn
}

resource "aws_lambda_event_source_mapping" "user_engagement_lambda_source_mapping" {
  event_source_arn = aws_dynamodb_table.adaptive_user_engagements_dynamo_table.stream_arn
  function_name = module.user_engagement_lambda.function_arn
  starting_position = "LATEST"
}

