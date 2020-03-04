locals {
  user_objectives_lambda_function_name_suffix = "user-objectives-lambda-go"
  user_objectives_lambda_function_name = "${var.client_id}_${local.user_objectives_lambda_function_name_suffix}"
}

module "user_objectives_lambda" {
  source = "../../../terraform-modules/adaptive-lambda"

  client_id     = local.client_id
  filename      = data.archive_file.adaptive-lambda-zip.output_path
  source_hash   = data.archive_file.adaptive-lambda-zip.output_base64sha256
  handler       = "adaptive"
  function_name_suffix = local.user_objectives_lambda_function_name_suffix
  runtime       = var.lambda_runtime
  timeout       = var.lambda_timeout

  memory_size   = var.multi_core_memory_size

  // Add environment variables.
  environment_variables = merge(local.environment_variables, {
    LAMBDA_ROLE   = "user-objectives"
    LOG_NAMESPACE = "user-objectives"
    # REPORTS_BUCKET_NAME = aws_s3_bucket.adaptive-feedback-reports-bucket.bucket
    # USER_OBJECTIVES_LEARN_MORE_PATH          = "user-objectives"
    # USER_OBJECTIVES_CLOSEOUT_LEARN_MORE_PATH = "user-objectives"

  })

  reserved_concurrent_executions = -1

  // Schedule the lambda
  schedule             = true
  schedule_name        = "user_objectives_lambda_warmer"
  schedule_description = "User Objectives Lambda Warmer for ${local.client_id}"
  schedule_expression  = "rate(5 minutes)"
  # schedule_invoke_json = data.local_file.sns_lambda_warmer_json.content

  // Attach extra policy
  attach_policy = true
  policy        = data.aws_iam_policy_document.user_objectives_policy.json

  tags = local.default_tags

}

data "aws_iam_policy_document" "user_objectives_policy" {
  # statement {
  #   actions   = [
  #     "dynamodb:GetItem",
  #   ]
  #   resources = [
  #     local.user_communities_table_arn,
  #     local.initiative_communities_table_arn,
  #     local.capability_communities_table_arn,
  #     local.strategy_initiatives_table_arn,
  #     local.strategy_objectives_table_arn,
  #     local.strategy_communities_table_arn,
  #     local.adaptive_vision_table_arn,
  #     local.postponed_event_table_arn,
  #   ]
  # }

  # statement {
  #   actions   = [
  #     "dynamodb:UpdateItem",
  #     "dynamodb:PutItem",]
  #   resources = [
  #     local.user_objectives_table_arn,
  #     local.coaching_relationships_dynamo_table_arn,
  #   ]
  # }
  # statement {
  #   actions   = [
  #     "dynamodb:PutItem",
  #     "dynamodb:GetItem",
  #     "dynamodb:Scan",
  #     "dynamodb:UpdateItem",]
  #   resources = [
  #     aws_dynamodb_table.user_objectives_progress.arn,]
  # }

  # statement {
  #   actions   = [
  #     "dynamodb:PutItem",
  #     "dynamodb:DeleteItem",
  #     "dynamodb:Scan",]
  #   resources = [
  #     "arn:aws:dynamodb:${local.region}:${data.aws_caller_identity.current.account_id}:table/${local.users_table_name}",]
  # }

  # statement {
  #   actions   = [
  #     "dynamodb:PutItem",
  #   ]
  #   resources = [
  #     aws_dynamodb_table.user_objectives_progress.arn,
  #     local.accountability_partnership_rejections_table_arn,
  #     local.postponed_event_table_arn,
  #   ]
  # }

  # statement {
  #   actions   = [
  #     "dynamodb:Query",]
  #   resources = [
  #     "arn:aws:dynamodb:${local.region}:${data.aws_caller_identity.current.account_id}:table/${local.users_table_name}/index/${local.users_platform_index}",
  #     local.user_objectives_table_arn,
  #     aws_dynamodb_table.user_objectives_progress.arn,
  #     local.capability_communities_table_arn,
  #     "${local.capability_communities_table_arn}/index/*",
  #     "${local.user_objectives_table_arn}/index/*",
  #     "${aws_dynamodb_table.user_objectives_progress.arn}/index/*",
  #     "arn:aws:dynamodb:${local.region}:${data.aws_caller_identity.current.account_id}:table/${local.adaptive_dialog_table_name}/index/*",
  #     "${local.dynamo_community_users_table_arn}/index/${local.dynamo_community_users_community_index}",
  #     "${local.strategy_objectives_table_arn}/index/*",
  #     "${local.strategy_initiatives_table_arn}/index/*",
  #     "${local.dynamo_community_users_table_arn}/index/*",
  #     "${local.values_table_arn}/index/${local.dynamo_adaptive_values_platform_id_index}",
  #     # "${local.dialog_table_arn}",
  #     # "arn:aws:dynamodb:${local.region}:${data.aws_caller_identity.current.account_id}:table/${local.adaptive_dialog_table_name}/index/${local.dynamo_dialog_content_contect_subject_index}",
  #     # "${local.adaptive_dialog_table_arn}/index/*",
  #   ]
  # }

  # statement {
  #   actions   = [
  #     "dynamodb:GetItem",
  #   ]
  #   resources = [
  #     "${local.values_table_arn}",
  #     "${local.client_config_table_arn}",
  #   ]
  # }

  # statement {
  #   actions   = [
  #     "dynamodb:PutItem",
  #     "dynamodb:GetItem",
  #     "dynamodb:UpdateItem",]
  #   resources = [
  #     "arn:aws:dynamodb:${local.region}:${data.aws_caller_identity.current.account_id}:table/${local.user_engagements_table_name}",
  #     local.strategy_objectives_table_arn,
  #     local.strategy_initiatives_table_arn,
  #   ]
  # }

  # statement {
  #   actions   = ["dynamodb:GetItem"]
  #   resources = [
  #     aws_dynamodb_table.adaptive_users_dynamodb_table.arn,
  #     aws_dynamodb_table.user_objective_dynamodb_table.arn,
  #     aws_dynamodb_table.user_objectives_progress.arn,
  #   ]
  # }

  statement {
    actions   = ["lambda:InvokeFunction"]
    resources = [module.user_profile_lambda.function_arn]
  }

  statement {
    actions   = ["SNS:Publish"]
    resources = [local.platform_notification_topic_arn,]
  }

  statement {
    actions   = [
      "comprehend:DetectSyntax",
      "translate:TranslateText"]
    resources = ["*"]
  }
}

resource "aws_sns_topic_subscription" "user_objectives_lambda_sns" {
  topic_arn = aws_sns_topic.namespace_payload.arn
  protocol  = "lambda"
  endpoint  = module.user_objectives_lambda.function_arn
}

resource "aws_lambda_permission" "user_objectives_lambda_sns" {
  statement_id  = "AllowExecutionFromSNS"
  action        = "lambda:InvokeFunction"
  function_name = module.user_objectives_lambda.function_arn
  principal     = "sns.amazonaws.com"
  source_arn    = aws_sns_topic.namespace_payload.arn
}

resource "aws_iam_role_policy_attachment" "user_objectives_lambda_read_all_tables" {
  role       = module.user_objectives_lambda.role_name
  policy_arn = aws_iam_policy.read_all_tables.arn
}

resource "aws_iam_role_policy_attachment" "user_objectives_lambda_write_issues_policy_attachment" {
  role       = module.user_objectives_lambda.role_name
  policy_arn = aws_iam_policy.write_issues.arn
}
