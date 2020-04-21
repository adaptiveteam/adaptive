locals {
  user_setup_lambda_function_name_suffix = "user-setup-lambda-go"
  user_setup_lambda_function_name = "${var.client_id}_${local.user_setup_lambda_function_name_suffix}"
}
module "user_setup_lambda" {
  source = "../../../terraform-modules/adaptive-lambda-s3"
  s3_bucket = aws_s3_bucket.binary_bucket.bucket
  s3_key = aws_s3_bucket_object.adaptive_zip.key
  source_hash = data.archive_file.adaptive_lambda_zip.output_md5

  client_id     = var.client_id
  handler       = "adaptive"
  function_name_suffix = local.user_setup_lambda_function_name_suffix
  runtime       = var.lambda_runtime
  timeout       = var.lambda_timeout

  reserved_concurrent_executions = -1

  // Add environment variables.
  environment_variables = merge(local.environment_variables, {
    LAMBDA_ROLE   = "user-setup"
    LOG_NAMESPACE = "user-setup"
  })

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

