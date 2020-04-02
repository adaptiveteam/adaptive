locals {
  feedback_slack_message_processor_function_name_suffix = "feedback-slack-message-processor-lambda-go"
  feedback_slack_message_processor_function_name = "${var.client_id}_${local.feedback_slack_message_processor_function_name_suffix}"
}
module "feedback_slack_message_processor_lambda" {
  source = "../../../terraform-modules/adaptive-lambda"

  client_id     = var.client_id
  filename      = data.archive_file.adaptive-lambda-zip.output_path
  source_hash   = data.archive_file.adaptive-lambda-zip.output_base64sha256
  handler       = "adaptive"
  function_name_suffix = local.feedback_slack_message_processor_function_name_suffix
  runtime       = var.lambda_runtime
  timeout       = var.lambda_timeout

  reserved_concurrent_executions = -1

  // Add environment variables.
  environment_variables = merge(local.environment_variables, {
    LAMBDA_ROLE   = "feedback-slack-message-processor"
    LOG_NAMESPACE = "feedback-slack-message-processor"
  })

  // Attach extra policy
  attach_policy = true
  policy        = data.aws_iam_policy_document.feedback_slack_message_processor_policy.json

  tags = local.default_tags

  // Schedule the lambda
  schedule             = true
  schedule_name        = "feedback_slack_message_processor_lambda_warmer"
  schedule_description = "Feedback slack_message_processor lambda warmer for ${local.client_id}"
  schedule_expression  = "rate(5 minutes)"
  schedule_invoke_json = data.local_file.feedback_slack_message_processor_lambda_warmer_json.content
}

data "local_file" "feedback_slack_message_processor_lambda_warmer_json" {
  filename = "${path.module}/templates/feedback-slack-message-processor-warmup.json"
}

data "aws_iam_policy_document" "feedback_slack_message_processor_policy" {
  statement {
    actions   = ["dynamodb:PutItem"]
    resources = [aws_dynamodb_table.adaptive_user_engagements_dynamo_table.arn]
  }

  statement {
    actions   = ["lambda:InvokeFunction"]
    resources = [
      module.feedback_setup_lambda.function_arn,
      module.feedback_reporting_lambda.function_arn,
      module.feedback_report_posting_lambda.function_arn,
    ]
  }

  statement {
    actions   = ["SNS:Publish"]
    resources = [aws_sns_topic.platform_notification.arn,]
  }

  statement {
    actions = ["dynamodb:*"]
    resources = [
      aws_dynamodb_table.user_objective_dynamodb_table.arn,
      "${aws_dynamodb_table.user_objective_dynamodb_table.arn}/index/*"]
  }

}

resource "aws_iam_role_policy_attachment" "feedback_slack_message_processor_lambda_read_all_tables" {
  role       = module.feedback_slack_message_processor_lambda.role_name
  policy_arn = aws_iam_policy.read_all_tables.arn
}
