locals {
  slack_message_processor_lambda_function_name_suffix = "slack-message-processor-lambda-go"
  slack_message_processor_lambda_function_name = "${var.client_id}_${local.slack_message_processor_lambda_function_name_suffix}"
}

resource "aws_lambda_function" "slack_message_processor_lambda" {
  s3_bucket = aws_s3_bucket_object.adaptive_binary.bucket
  s3_key = aws_s3_bucket_object.adaptive_binary.key
  # s3_object_version = "latest"
  handler          = "adaptive"

  description      = "${var.client_id} slack_message_processor_lambda"
  function_name    = local.slack_message_processor_lambda_function_name
  source_code_hash = data.archive_file.adaptive-lambda-zip.output_base64sha256


  runtime = var.lambda_runtime
  timeout = var.lambda_timeout
  memory_size = var.multi_core_memory_size

  reserved_concurrent_executions = -1

  role = aws_iam_role.lambda.arn

  environment {
    variables =  merge(local.environment_variables, {
      LAMBDA_ROLE   = "slack-message-processor"
      LOG_NAMESPACE = "slack-message-processor"
    })
  }

  tags = local.default_tags

}

module "slack_message_processor_lambda_warmer" {
  source = "../../../terraform-modules/warmer"

  client_id = var.client_id
  target_lambda_arn = aws_lambda_function.slack_message_processor_lambda.arn
  schedule_name = "slack_message_processor_lambda_warmer_every_5_minutes"
}

module "slack_message_processor_lambda_log_policy" {
  source = "../../../terraform-modules/log-policy"

  client_id = var.client_id
  function_name = aws_lambda_function.slack_message_processor_lambda.function_name
  errors_sns_topic_arn = aws_sns_topic.errors.arn
  role_name = aws_iam_role.lambda.name
}

# module "slack_message_processor_lambda" {
#   source = "../../../terraform-modules/adaptive-lambda"

#   client_id = var.client_id
#   filename = data.archive_file.adaptive-lambda-zip.output_path
#   source_hash = data.archive_file.adaptive-lambda-zip.output_base64sha256
#   handler = "adaptive"
#   function_name_suffix = local.slack_message_processor_lambda_function_name_suffix
#   runtime = var.lambda_runtime
#   timeout = var.lambda_timeout
#   memory_size = var.multi_core_memory_size

#   reserved_concurrent_executions = -1

#   // Schedule the lambda
#   schedule = true
#   schedule_name = "slack_message_processor_lambda-warmer"
#   schedule_description = "Slack Message Processor lambda warmer for ${var.client_id}"
#   schedule_expression = "rate(5 minutes)"
#   schedule_invoke_json = data.local_file.slack_message_processor_lambda_warmer_json.content

#   // Add environment variables.
#   environment_variables = merge(local.environment_variables, {
#     LAMBDA_ROLE   = "slack-message-processor"
#     LOG_NAMESPACE = "slack-message-processor"
#   })

#   // Attach extra policy
#   attach_policy = true
#   policy = data.aws_iam_policy_document.slack_message_processor_dynamo_write_policy.json

#   tags = local.default_tags
# }
data "local_file" "slack_message_processor_lambda_warmer_json" {
  filename = "${path.module}/templates/api_slack_warmup.json"
}
