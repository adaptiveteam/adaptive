data "archive_file" "platform_engagement_scheduler_lambda_zip" {
  type        = "zip"
  source_file = "../../../bin/platform-engagement-scheduler-lambda-go"
  output_path = "lambdas/platform-engagement-scheduler-lambda-go.zip"
}

module "platform_engagement_scheduler_lambda" {
  source = "../../../terraform-modules/adaptive-lambda"

  client_id     = var.client_id
  filename      = data.archive_file.platform_engagement_scheduler_lambda_zip.output_path
  source_hash   = data.archive_file.platform_engagement_scheduler_lambda_zip.output_base64sha256
  function_name = "platform-engagement-scheduler-lambda-go"
  handler       = "platform-engagement-scheduler-lambda-go"
  runtime       = var.lambda_runtime
  timeout       = var.lambda_timeout
  memory_size   = var.multi_core_memory_size

  reserved_concurrent_executions = -1

  // Attach extra policy
  attach_policy = true
  policy        = data.aws_iam_policy_document.platform_engagement_scheduler_lambda_policy.json

  tags = local.default_tags

  environment_variables = {
    VISION_TABLE_NAME                  = aws_dynamodb_table.vision.name
    STRATEGY_OBJECTIVES_TABLE_NAME     = aws_dynamodb_table.strategy_objectives.name
    USER_OBJECTIVES_TABLE_NAME         = aws_dynamodb_table.user_objective_dynamodb_table.name
    STRATEGY_OBJECTIVES_PLATFORM_INDEX = var.dynamo_strategy_objectives_platform_index
    USER_OBJECTIVES_ID_INDEX           = var.dynamo_user_objectives_id_index
    CLIENT_CONFIG_TABLE_NAME           = aws_dynamodb_table.client_config_dynamodb_table.name
    CLIENT_ID                          = var.client_id
    ADAPTIVE_COMMUNITIES_TABLE         = aws_dynamodb_table.user_communities.name
    LOG_NAMESPACE                      = "platform-notification"
  }
}

data "aws_iam_policy_document" "platform_engagement_scheduler_lambda_policy" {
  statement {
    effect    = "Allow"
    actions   = [
      "dynamodb:Scan",
      "dynamodb:GetItem",
    ]
    resources = [
      aws_dynamodb_table.client_config_dynamodb_table.arn,
    ]
  }
  statement {
    effect    = "Allow"
    actions   = [
      "dynamodb:GetItem",
    ]
    resources = [
      aws_dynamodb_table.vision.arn,
      aws_dynamodb_table.user_communities.arn,
    ]
  }
  statement {
    effect    = "Allow"
    actions   = [
      "dynamodb:Query",
    ]
    resources = [
      "${aws_dynamodb_table.strategy_objectives.arn}/index/${var.dynamo_strategy_objectives_platform_index}",
      "${aws_dynamodb_table.user_objective_dynamodb_table.arn}/index/${var.dynamo_user_objectives_id_index}",
    ]
  }
}

