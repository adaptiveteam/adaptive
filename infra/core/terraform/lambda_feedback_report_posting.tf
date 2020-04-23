locals {
  feedback_report_posting_function_name_suffix = "feedback-report-posting-lambda-go"
  feedback_report_posting_function_name = "${var.client_id}_${local.feedback_report_posting_function_name_suffix}"
}
module "feedback_report_posting_lambda" {
  source = "../../../terraform-modules/adaptive-lambda-s3"
  s3_bucket = aws_s3_bucket.binary_bucket.bucket
  s3_key = aws_s3_bucket_object.adaptive_zip.key
  source_hash = data.archive_file.adaptive_lambda_zip.output_md5

  client_id     = var.client_id
  handler       = "adaptive"
  function_name_suffix = local.feedback_report_posting_function_name_suffix
  runtime       = var.lambda_runtime
  timeout       = var.lambda_timeout

  reserved_concurrent_executions = -1

  // Add environment variables.
  environment_variables = merge(local.environment_variables, {
    LAMBDA_ROLE   = "feedback-report-posting"
    LOG_NAMESPACE = "feedback-report-posting"
  })

  // Attach extra policy
  attach_policy = true
  policy        = data.aws_iam_policy_document.feedback_report_posting_policy.json

  tags = local.default_tags
}

data "aws_iam_policy_document" "feedback_report_posting_policy" {
  statement {
    actions   = [
      "s3:GetObject",
      "s3:GetObjectAcl",]
    resources = [
      "${aws_s3_bucket.adaptive-feedback-reports-bucket.arn}/*",
    ]
  }

  statement {
    actions   = ["lambda:InvokeFunction"]
    resources = [
      module.feedback_reporting_lambda.function_arn,
      module.user_profile_lambda.function_arn,
    ]
  }

  statement {
    actions   = ["SNS:Publish"]
    resources = [aws_sns_topic.platform_notification.arn]
  }
}

resource "aws_iam_role_policy_attachment" "feedback_report_posting_lambda_read_all_tables" {
  role       = module.feedback_report_posting_lambda.role_name
  policy_arn = aws_iam_policy.read_all_tables.arn
}
