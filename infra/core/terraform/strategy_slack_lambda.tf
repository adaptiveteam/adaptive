locals {
  strategy_slack_lambda_function_name_suffix = "strategy-slack-message-processor-lambda-go"
  strategy_slack_lambda_function_name = "${var.client_id}_${local.strategy_slack_lambda_function_name_suffix}"
}
module "strategy_slack_lambda" {
  source = "../../../terraform-modules/adaptive-lambda"

  client_id     = var.client_id
  filename      = data.archive_file.adaptive-lambda-zip.output_path
  source_hash   = data.archive_file.adaptive-lambda-zip.output_base64sha256
  handler       = "adaptive"
  function_name_suffix = local.strategy_slack_lambda_function_name_suffix
  runtime       = var.lambda_runtime
  timeout       = var.lambda_timeout

  reserved_concurrent_executions = -1

  // Add environment variables.
  environment_variables = merge(local.environment_variables, {
    LAMBDA_ROLE   = "strategy-slack-message-processor"
    LOG_NAMESPACE = "strategy-slack-message-processor"
  })

  // Attach extra policy
  attach_policy = true
  policy        = data.aws_iam_policy_document.strategy_slack_policy.json

  tags = local.default_tags


  // Schedule the lambda
  schedule             = true
  schedule_name        = "strategy_objectives_lambda_warmer"
  schedule_description = "Strategy Objectives Lambda Warmer for ${local.client_id}"
  schedule_expression  = "rate(5 minutes)"
  # schedule_invoke_json = data.local_file.strategy_objectives_lambda_warmer_json.content

}

data "aws_iam_policy_document" "strategy_slack_policy" {

  statement {
    actions   = [
      "dynamodb:PutItem",
      "dynamodb:UpdateItem",
    ]
    resources = [
      aws_dynamodb_table.adaptive_user_engagements_dynamo_table.arn,
      aws_dynamodb_table.strategy_communities.arn,
    ]
  }

  statement {
    actions   = [
      "dynamodb:PutItem",
      "dynamodb:GetItem",
      "dynamodb:UpdateItem",
    ]
    resources = [
      aws_dynamodb_table.strategy_objectives.arn,
      aws_dynamodb_table.user_objectives_progress.arn,
      ]
  }

  statement {
    actions   = [
      "dynamodb:PutItem",
      "dynamodb:GetItem",
    ]
    resources = [
      aws_dynamodb_table.vision.arn,
      aws_dynamodb_table.capability_communities.arn,
      aws_dynamodb_table.strategy_initiatives.arn,
      aws_dynamodb_table.initiative_communities.arn,
      aws_dynamodb_table.user_objective_dynamodb_table.arn,
      aws_dynamodb_table.user_objectives_progress.arn,
    ]
  }

  statement {
    actions   = [
      "dynamodb:GetItem",
    ]
    resources = [
      aws_dynamodb_table.user_communities.arn,
      aws_dynamodb_table.strategy_communities.arn,
      aws_dynamodb_table.adaptive_users_dynamodb_table.arn,
      aws_dynamodb_table.client_config_dynamodb_table.arn,
    ]
  }

  statement {
    actions   = [
      "dynamodb:Query",]
    resources = [
      "${aws_dynamodb_table.adaptive_users_dynamodb_table.arn}/index/*",
      "${aws_dynamodb_table.adaptive_dialog_content.arn}/index/*",
      "${aws_dynamodb_table.user_objective_dynamodb_table.arn}/index/*",
      "${aws_dynamodb_table.community_users.arn}/index/*",
      "${aws_dynamodb_table.strategy_objectives.arn}/index/*",
      "${aws_dynamodb_table.capability_communities.arn}/index/*",
      "${aws_dynamodb_table.initiative_communities.arn}/index/*",
      "${aws_dynamodb_table.strategy_initiatives.arn}/index/*",
      "${aws_dynamodb_table.user_communities.arn}/index/*",
      "${aws_dynamodb_table.strategy_communities.arn}/index/*",
      "${aws_dynamodb_table.user_objective_dynamodb_table.arn}/index/*",
      "${aws_dynamodb_table.user_objectives_progress.arn}/index/*",
      aws_dynamodb_table.user_objectives_progress.arn,
    ]
  }

  statement {
    actions   = ["lambda:InvokeFunction",]
    resources = [module.user_profile_lambda.function_arn]
      # "arn:aws:lambda:${local.region}:${data.aws_caller_identity.current.account_id}:function:${local.user_profile_lambda_name}",]
  }

  statement {
    actions   = ["dynamodb:Scan",]
    resources = [aws_dynamodb_table.adaptive_users_dynamodb_table.arn]
  }

  // required for Adaptive NLP
  statement {
    actions   = [
      "comprehend:DetectSyntax",
      "translate:TranslateText"]
    resources = [
      "*",
    ]
  }

  statement {
    actions   = ["SNS:Publish"]
    resources = [local.platform_notification_topic_arn,]
  }
}

