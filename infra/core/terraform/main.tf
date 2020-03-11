# Specify the provider and access details
provider "aws" {
  region = var.aws_region
  version = "~> 2.52.0"
}

data "aws_caller_identity" "current" {
}

locals {
  //  adaptive_values_table_name = "${var.client_id}_adaptive_values"
  //  dynamo_adaptive_values_platform_id_index = "AdaptiveValuesPlatformIdIndex"

  default_tags = {
    Environment = var.environment
    ManagedBy = "Terraform"
  }

  platform_notification_topic_arn = aws_sns_topic.platform_notification.arn
  // TODO: Move these to core state
  reporting_lambda_name        = "${var.client_id}_feedback-reporting-lambda-go"
  report_posting_lambda_name   = "${var.client_id}_feedback-report-posting-lambda-go"
  # feedback_reports_bucket_name = aws_s3_bucket.adaptive-feedback-reports-bucket.bucket

  # reporting_lambda_name        = data.terraform_remote_state.core.outputs.reporting_lambda_name
  # report_posting_lambda_name   = data.terraform_remote_state.feedback.outputs.report_posting_lambda_name
  # feedback_reports_bucket_name = data.terraform_remote_state.feedback.outputs.feedback_reports_bucket_name

  slack_message_processor_suffix = "slack-message-processor-lambda-go"
  region = var.aws_region
  client_id = var.client_id
}

