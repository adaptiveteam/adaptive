locals {
  user_objectives_lambda_function_name_suffix = "user-objectives-lambda-go"
  user_objectives_lambda_function_name = "${var.client_id}_${local.user_objectives_lambda_function_name_suffix}"
}

module "user_objectives_lambda" {
  source = "../../../terraform-modules/adaptive-lambda-s3"
  s3_bucket = aws_s3_bucket.binary_bucket.bucket
  s3_key = aws_s3_bucket_object.adaptive_zip.key
  source_hash = data.archive_file.adaptive_lambda_zip.output_md5

  client_id     = local.client_id
  handler       = "adaptive"
  function_name_suffix = local.user_objectives_lambda_function_name_suffix
  runtime       = var.lambda_runtime
  timeout       = var.lambda_timeout

  memory_size   = var.multi_core_memory_size

  // Add environment variables.
  environment_variables = merge(local.environment_variables, {
    LAMBDA_ROLE   = "user-objectives"
    LOG_NAMESPACE = "user-objectives"
    # REPORTS_BUCKET_NAME = aws_s3_bucket.adaptive-feedback-reports-bucket.bucket
    # USER_OBJECTIVES_LEARN_MORE_PATH          = "user-objectives"
    # USER_OBJECTIVES_CLOSEOUT_LEARN_MORE_PATH = "user-objectives"

  })

  reserved_concurrent_executions = -1

  // Schedule the lambda
  schedule             = true
  schedule_name        = "user_objectives_lambda_warmer"
  schedule_description = "User Objectives Lambda Warmer for ${local.client_id}"
  schedule_expression  = "rate(5 minutes)"
  # schedule_invoke_json = data.local_file.sns_lambda_warmer_json.content

  // Attach extra policy
  attach_policy = true
  policy        = data.aws_iam_policy_document.user_objectives_policy.json

  tags = local.default_tags

}

data "aws_iam_policy_document" "user_objectives_policy" {

  statement {
    actions   = ["lambda:InvokeFunction"]
    resources = [module.user_profile_lambda.function_arn]
  }

  statement {
    actions   = ["SNS:Publish"]
    resources = [local.platform_notification_topic_arn,]
  }

  statement {
    actions   = [
      "comprehend:DetectSyntax",
      "translate:TranslateText"]
    resources = ["*"]
  }
}

resource "aws_sns_topic_subscription" "user_objectives_lambda_sns" {
  topic_arn = aws_sns_topic.namespace_payload.arn
  protocol  = "lambda"
  endpoint  = module.user_objectives_lambda.function_arn
}

resource "aws_lambda_permission" "user_objectives_lambda_sns" {
  statement_id  = "AllowExecutionFromSNS"
  action        = "lambda:InvokeFunction"
  function_name = module.user_objectives_lambda.function_arn
  principal     = "sns.amazonaws.com"
  source_arn    = aws_sns_topic.namespace_payload.arn
}

resource "aws_iam_role_policy_attachment" "user_objectives_lambda_read_all_tables" {
  role       = module.user_objectives_lambda.role_name
  policy_arn = aws_iam_policy.read_all_tables.arn
}

resource "aws_iam_role_policy_attachment" "user_objectives_lambda_write_issues_policy_attachment" {
  role       = module.user_objectives_lambda.role_name
  policy_arn = aws_iam_policy.write_issues.arn
}
