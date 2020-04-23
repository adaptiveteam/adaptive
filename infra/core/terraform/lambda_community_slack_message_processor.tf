locals {
  community_slack_lambda_function_name_suffix = "community-slack-message-processor-lambda-go"
  community_slack_lambda_function_name = "${var.client_id}_${local.community_slack_lambda_function_name_suffix}"
}

module "community_slack_message_processor_lambda" {
  source = "../../../terraform-modules/adaptive-lambda-s3"
  s3_bucket = aws_s3_bucket.binary_bucket.bucket
  s3_key = aws_s3_bucket_object.adaptive_zip.key
  source_hash = data.archive_file.adaptive_lambda_zip.output_md5

  client_id     = local.client_id
  handler       = "adaptive"
  function_name_suffix = local.community_slack_lambda_function_name_suffix
  runtime       = var.lambda_runtime
  timeout       = 600 // 600 seconds = 10 minutes; because Current/Next Quarter Events handling take too long.

  memory_size   = var.multi_core_memory_size

  reserved_concurrent_executions = -1

  // Add environment variables.
  environment_variables = merge(local.environment_variables, {
    LAMBDA_ROLE   = "community-slack-message-processor"
    LOG_NAMESPACE = "community-slack-message-processor"
  })

  // Schedule the lambda
  schedule             = true
  schedule_name        = "community_slack_message_processor_lambda_warmer"
  schedule_description = "Community Slack Message Processor Lambda Warmer for ${local.client_id}"
  schedule_expression  = "rate(5 minutes)"
  # schedule_invoke_json = data.local_file.sns_lambda_warmer_json.content

  // Attach extra policy
  attach_policy = true
  policy        = data.aws_iam_policy_document.community-slack-message-processor-policy.json

  tags = local.default_tags

}

data "aws_iam_policy_document" "community-slack-message-processor-policy" {
  statement {
    actions   = [
      "dynamodb:PutItem",
      "dynamodb:DeleteItem",
      "dynamodb:UpdateItem",
      ]
    resources = [
      aws_dynamodb_table.strategy_communities.arn,
      aws_dynamodb_table.strategy_initiatives.arn,
      aws_dynamodb_table.strategy_objectives.arn,

      aws_dynamodb_table.user_communities.arn,
      aws_dynamodb_table.community_users.arn,
      aws_dynamodb_table.adaptive_users_dynamodb_table.arn,
      aws_dynamodb_table.adaptive_user_engagements_dynamo_table.arn,

      aws_dynamodb_table.user_objective_dynamodb_table.arn,
      aws_dynamodb_table.user_objectives_progress.arn,
   
      aws_dynamodb_table.coaching_relationships.arn,
      aws_dynamodb_table.coaching_rejections.arn,
      ]
  }

  statement {
    actions   = ["lambda:InvokeFunction"]
    resources = [
      module.user_profile_lambda.function_arn,
      "arn:aws:lambda:${local.region}:${data.aws_caller_identity.current.account_id}:function:${local.report_posting_lambda_name}",
      "arn:aws:lambda:${local.region}:${data.aws_caller_identity.current.account_id}:function:${local.reporting_lambda_name}",

      # module.feedback_report_posting_lambda.function_arn,
      # module.feedback_reporting_lambda.function_arn,
      module.user_engagement_scripting_lambda.function_arn,
      module.user_setup_lambda.function_arn,
      module.user_engagement_scheduler_lambda.function_arn,
    ]
  }

  statement {
    actions   = [
      "s3:GetObject",
      "s3:GetObjectAcl",]
    resources = [
      "arn:aws:s3:::${aws_s3_bucket.adaptive-feedback-reports-bucket.bucket}/*",
    ]
  }

  // TEMP
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
    resources = [local.platform_notification_topic_arn]
  }
}

resource "aws_sns_topic_subscription" "feedback_slack_message_processor_lambda_sns" {
  topic_arn = aws_sns_topic.namespace_payload.arn
  protocol  = "lambda"
  endpoint  = module.community_slack_message_processor_lambda.function_arn
}

resource "aws_lambda_permission" "feedback_slack_message_processor_lambda_sns" {
  statement_id  = "AllowExecutionFromSNS"
  action        = "lambda:InvokeFunction"
  function_name = module.community_slack_message_processor_lambda.function_arn
  principal     = "sns.amazonaws.com"
  source_arn    = aws_sns_topic.namespace_payload.arn
}

resource "aws_iam_role_policy_attachment" "community_slack_message_processor_lambda_read_all_tables" {
  role       = module.community_slack_message_processor_lambda.role_name
  policy_arn = aws_iam_policy.read_all_tables.arn
}
