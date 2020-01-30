module "platform_engagement_scheduler_lambda" {
  source = "../../../terraform-modules/adaptive-lambda"

  client_id     = var.client_id
  filename = data.archive_file.adaptive-lambda-zip.output_path
  source_hash = data.archive_file.adaptive-lambda-zip.output_base64sha256
  handler = "adaptive"
  function_name_suffix = "platform-engagement-scheduler-lambda-go"
  runtime       = var.lambda_runtime
  timeout       = var.lambda_timeout
  memory_size   = var.multi_core_memory_size

  reserved_concurrent_executions = -1

  // Attach extra policy
  attach_policy = true
  policy        = data.aws_iam_policy_document.platform_engagement_scheduler_lambda_policy.json

  tags = local.default_tags

  // Add environment variables.
  environment_variables = merge(local.environment_variables, {
    LAMBDA_ROLE   = "platform-engagement-scheduler"
    LOG_NAMESPACE = "platform-engagement-scheduler"
  })

}

data "aws_iam_policy_document" "platform_engagement_scheduler_lambda_policy" {
  statement {
    resources = [aws_dynamodb_table.client_config_dynamodb_table.arn]
    actions   = [
      "dynamodb:Scan",
      "dynamodb:GetItem",
    ]
  }
  statement {
    actions   = [
      "dynamodb:GetItem",
    ]
    resources = [
      aws_dynamodb_table.vision.arn,
      aws_dynamodb_table.user_communities.arn,
    ]
  }
  statement {
    actions   = [
      "dynamodb:Query",
    ]
    resources = [
      "${aws_dynamodb_table.strategy_objectives.arn}/index/${var.dynamo_strategy_objectives_platform_index}",
      "${aws_dynamodb_table.user_objective_dynamodb_table.arn}/index/${var.dynamo_user_objectives_id_index}",
    ]
  }
}

