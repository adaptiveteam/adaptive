
locals {
  slack_user_query_lambda_function_name_suffix = "slack-user-query-lambda-go"
  slack_user_query_lambda_function_name = "${var.client_id}_${local.slack_user_query_lambda_function_name_suffix}"
}

module "slack_user_query_lambda" {
  source = "../../../terraform-modules/adaptive-lambda-s3"
  s3_bucket = aws_s3_bucket.binary_bucket.bucket
  s3_key = aws_s3_bucket_object.adaptive_zip.key
  source_hash = data.archive_file.adaptive_lambda_zip.output_md5

  client_id = var.client_id
  handler = "adaptive"
  function_name_suffix = local.slack_user_query_lambda_function_name_suffix
  runtime = var.lambda_runtime
  timeout = 900
  memory_size = var.multi_core_memory_size

  reserved_concurrent_executions = -1

  // Add environment variables.
  environment_variables = merge(local.environment_variables, {
    LAMBDA_ROLE   = "slack-user-query"
    LOG_NAMESPACE = "slack-user-query"
  })

  // Attach extra policy
  attach_policy = true
  policy = data.aws_iam_policy_document.slack_user_query_policy.json

  tags = local.default_tags
}

data "aws_iam_policy_document" "slack_user_query_policy" {
  statement {
    actions = ["dynamodb:*"]
    resources = [
      aws_dynamodb_table.adaptive_users_dynamodb_table.arn,
      aws_dynamodb_table.user_communities.arn,
      aws_dynamodb_table.community_users.arn,
    ]
  }

  statement {
    actions = ["lambda:InvokeFunction"]
    resources = [
      module.user_setup_lambda.function_arn,
    ]
  }

  statement {
    actions = ["dynamodb:GetItem"]
    resources = [
      aws_dynamodb_table.client_config_dynamodb_table.arn,
    ]
  }

  statement {
    actions = ["dynamodb:Query"]
    resources = [
      "${aws_dynamodb_table.user_communities.arn}/index/*",
      "${aws_dynamodb_table.community_users.arn}/index/*",
      "${aws_dynamodb_table.adaptive_users_dynamodb_table.arn}/index/*",
    ]
  }
}

