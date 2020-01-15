# Specify the provider and access details
provider "aws" {
  region = var.aws_region
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
  reporting_lambda_name = "${var.client_id}_feedback-reporting-lambda-go"

  slack_message_processor_suffix = "slack-message-processor-lambda-go"
  region = var.aws_region
  client_id = var.client_id
}

