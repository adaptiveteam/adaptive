data "aws_iam_policy_document" "competencies_policy" {
  statement {
    actions   = ["dynamodb:GetItem"]
    resources = [
      aws_dynamodb_table.client_config_dynamodb_table.arn,
      aws_dynamodb_table.adaptive_users_dynamodb_table.arn,
    ]
  }
  statement {
    actions   = [
      "dynamodb:PutItem",
      "dynamodb:GetItem",
      "dynamodb:Scan",
      "dynamodb:UpdateItem",
    ]
    resources = [aws_dynamodb_table.adaptive_value_dynamodb_table.arn]
  }

  statement {
    actions   = [
      "dynamodb:PutItem",
      "dynamodb:GetItem",
      "dynamodb:Query",]
    resources = [
      "${aws_dynamodb_table.adaptive_value_dynamodb_table.arn}/index/*",
      "${aws_dynamodb_table.community_users.arn}/index/*",
      "${aws_dynamodb_table.adaptive_dialog_content.arn}/index/*",
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

resource "aws_iam_policy" "competencies_additional_policy" {
  name   = "${var.client_id}_competencies_additional_policy"
  policy = data.aws_iam_policy_document.competencies_policy.json
}
