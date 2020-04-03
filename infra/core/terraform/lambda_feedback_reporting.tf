locals {
  feedback_reporting_function_name_suffix = "feedback-reporting-lambda-go"
  feedback_reporting_function_name = "${var.client_id}_${local.feedback_reporting_function_name_suffix}"
}
module "feedback_reporting_lambda" {
  source = "../../../terraform-modules/adaptive-lambda"

  client_id     = var.client_id
  filename      = data.archive_file.adaptive-lambda-zip.output_path
  source_hash   = data.archive_file.adaptive-lambda-zip.output_base64sha256
  handler       = "adaptive"
  function_name_suffix = local.feedback_reporting_function_name_suffix
  runtime       = var.lambda_runtime
  timeout       = var.lambda_timeout

  reserved_concurrent_executions = -1

  // Add environment variables.
  environment_variables = merge(local.environment_variables, {
    LAMBDA_ROLE   = "feedback-reporting"
    LOG_NAMESPACE = "feedback-reporting"
  })

  // Attach extra policy
  attach_policy = true
  policy        = data.aws_iam_policy_document.feedback_reporting_policy.json

  tags = local.default_tags
}

resource "aws_iam_role_policy_attachment" "feedback_reporting_lambda_read_all_tables" {
  role       = module.feedback_reporting_lambda.role_name
  policy_arn = aws_iam_policy.read_all_tables.arn
}

data "aws_iam_policy_document" "feedback_reporting_policy" {
  statement {
    actions   = [
      "s3:PutObject",
      "s3:PutObjectAcl",]
    resources = [
      "${aws_s3_bucket.adaptive-feedback-reports-bucket.arn}/*",
    ]
  }

  statement {
    actions   = ["SNS:Publish"]
    resources = [aws_sns_topic.platform_notification.arn]
  }
  statement {
    actions   = ["lambda:InvokeFunction"]
    resources = [module.user_profile_lambda.function_arn]
  }

}
