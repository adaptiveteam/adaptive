# module "holidays-lambda-go" {
#   source = "../../../terraform-modules/adaptive-go-lambda"
#   // Relative path is convenient for fast development roundtrip
#   modpath             = "lambdas/"
#   input_sns_topic_arn = aws_sns_topic.namespace_payload.arn
#   id                  = var.client_id
#   name                = "adaptive-holidays-lambda-go"
#   memory_size         = var.multi_core_memory_size
#   # // Attach extra policy
#   # attach_policy = true
#   # policy = "${data.aws_iam_policy_document.holidays_policy.json}"

#   tags = local.default_tags

#   // Add environment variables.
#   environment_variables = merge(local.environment_variables, {
#     LAMBDA_ROLE   = "holidays"
#     LOG_NAMESPACE = "holidays"
#   })
#   # environment_variables = {
#   #   USER_ENGAGEMENTS_TABLE_NAME           = local.user_engagements_table_name
#   #   PLATFORM_NOTIFICATION_TOPIC           = local.platform_notification_topic_arn
#   #   USER_PROFILE_LAMBDA_NAME              = module.user_profile_lambda.function_name
#   #   HOLIDAYS_AD_HOC_TABLE                 = local.holidays_table_name
#   #   HOLIDAYS_PLATFORM_DATE_INDEX          = local.dynamo_holidays_date_index
#   #   DIALOG_TABLE                          = local.adaptive_dialog_table_name
#   #   ADAPTIVE_DIALOG_CONTEXT_SUBJECT_INDEX = local.dynamo_dialog_content_contect_subject_index
#   #   HOLIDAYS_LEARN_MORE_PATH              = "holidays"
#   #   LOG_NAMESPACE                         = "holidays"
#   #   COMMUNITY_USERS_TABLE_NAME            = local.dynamo_community_users_table_name
#   #   COMMUNITY_USERS_COMMUNITY_INDEX       = local.dynamo_community_users_community_index
#   #   COMMUNITY_USERS_USER_COMMUNITY_INDEX  = local.dynamo_community_users_user_community_index
#   #   COMMUNITY_USERS_USER_INDEX            = local.dynamo_community_users_user_index
#   # }
# }

data "aws_iam_policy_document" "holidays_policy" {
  statement {
    actions   = [
      "dynamodb:PutItem",
      "dynamodb:GetItem",
      "dynamodb:Scan",
      "dynamodb:UpdateItem",
    ]
    resources = [aws_dynamodb_table.ad_hoc_holidays.arn]
  }

  statement {
    actions   = ["dynamodb:Query",]
    resources = [
      "${aws_dynamodb_table.ad_hoc_holidays.arn}/index/*",
      "${aws_dynamodb_table.community_users.arn}/index/*",
      "${aws_dynamodb_table.adaptive_dialog_content.arn}/index/*",
      "${aws_dynamodb_table.coaching_relationships.arn}/index/*",
    ]
  }

  statement {
    actions   = [
      "dynamodb:PutItem",
      "dynamodb:GetItem",
      "dynamodb:UpdateItem",]
    resources = [aws_dynamodb_table.adaptive_user_engagements_dynamo_table.arn]
  }

  statement {
    actions   = ["lambda:InvokeFunction",]
    resources = [module.user_profile_lambda.function_arn]
  }

  statement {
    actions   = ["SNS:Publish",]
    resources = [local.platform_notification_topic_arn,]
  }

  statement {
    actions   = [
      "comprehend:DetectSyntax",
      "translate:TranslateText"]
    resources = ["*"]
  }

}

resource "aws_iam_policy" "holidays_additional_policy" {
  name   = "${var.client_id}_holidays_additional_policy"
  policy = data.aws_iam_policy_document.holidays_policy.json
}
