locals {
  feedback_setup_function_name_suffix = "feedback-setup-lambda-go"
  feedback_setup_function_name = "${var.client_id}_${local.feedback_setup_function_name_suffix}"
}
module "feedback_setup_lambda" {
  source = "../../../terraform-modules/adaptive-lambda"

  client_id     = var.client_id
  filename      = data.archive_file.adaptive-lambda-zip.output_path
  source_hash   = data.archive_file.adaptive-lambda-zip.output_base64sha256
  handler       = "adaptive"
  function_name_suffix = local.feedback_setup_function_name_suffix
  runtime       = var.lambda_runtime
  timeout       = var.lambda_timeout

  reserved_concurrent_executions = -1

  // Add environment variables.
  environment_variables = merge(local.environment_variables, {
    LAMBDA_ROLE   = "feedback-setup"
    LOG_NAMESPACE = "feedback-setup"
  })

  // Attach extra policy
  attach_policy = true
  policy        = data.aws_iam_policy_document.feedback_setup_policy.json

  tags = local.default_tags
  // Schedule the lambda
  schedule             = true
  schedule_name        = "feedback_setup_lambda_warmer"
  schedule_description = "Feedback setup lambda warmer for ${local.client_id}"
  schedule_expression  = "rate(15 minutes)"
  schedule_invoke_json = "{\"payload\" : \"warmup\"}"

}

data "aws_iam_policy_document" "feedback_setup_policy" {
  statement {
    actions   = [
      "dynamodb:DescribeTable",
      "dynamodb:PutItem",
      "dynamodb:GetItem",
      "dynamodb:DeleteItem",
      "dynamodb:UpdateItem",]
    resources = [aws_dynamodb_table.adaptive_user_engagements_dynamo_table.arn]
  }

  statement {
    actions   = ["dynamodb:PutItem", 
      "dynamodb:DeleteItem",
      "dynamodb:UpdateItem",]
    resources = [aws_dynamodb_table.adaptive_user_feedback_dynamodb_table.arn,]
  }

  statement {
    actions   = ["lambda:InvokeFunction",]
    resources = [
      module.user_profile_lambda.function_arn,
      module.feedback_analysis_lambda.function_arn,
      module.feedback_report_posting_lambda.function_arn,
    ]
  }
  statement {
    actions   = ["lambda:InvokeFunction",]
    resources = [module.user_profile_lambda.function_arn]
  }


  statement {
    actions   = ["SNS:Publish"]
    resources = [aws_sns_topic.platform_notification.arn,]
  }

}

resource "aws_iam_role_policy_attachment" "feedback_setup_lambda_read_all_tables" {
  role       = module.feedback_setup_lambda.role_name
  policy_arn = aws_iam_policy.read_all_tables.arn
}
