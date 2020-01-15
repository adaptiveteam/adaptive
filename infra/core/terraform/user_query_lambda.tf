data "archive_file" "user-query-lambda-zip" {
  type        = "zip"
  source_file = "../../../bin/user-query-lambda-go"
  output_path = "lambdas/user-query-lambda-go.zip"
}

module "user_query_lambda" {
  source = "../../../terraform-modules/adaptive-lambda"

  client_id     = var.client_id
  filename      = data.archive_file.user-query-lambda-zip.output_path
  source_hash   = data.archive_file.user-query-lambda-zip.output_base64sha256
  function_name = "user-query-lambda-go"
  handler       = "user-query-lambda-go"
  runtime       = var.lambda_runtime
  timeout       = var.lambda_timeout
  memory_size   = var.multi_core_memory_size

  reserved_concurrent_executions = -1

  environment_variables = {
    SLACK_LAMBDA_FUNCTION_NAME = module.slack_user_query_lambda.function_name
    CLIENT_CONFIG_TABLE_NAME   = aws_dynamodb_table.client_config_dynamodb_table.name
    LOG_NAMESPACE              = "user-query"
  }

  // Attach extra policy
  attach_policy = true
  policy        = data.aws_iam_policy_document.user_query_policy.json

  // Schedule the lambda
  schedule             = true
  schedule_name        = "1pm_UTC_rule"
  schedule_description = "Cloudwatch event rule for 1PM UTC/ 8AM EST/ 6:30PM IST"
  schedule_expression  = "cron(0 13 * * ? *)"

  tags = local.default_tags
}

data "aws_iam_policy_document" "user_query_policy" {
  statement {
    effect    = "Allow"
    actions   = [
      "dynamodb:DescribeTable",
      "dynamodb:Query",
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
      "lambda:InvokeFunction",
    ]
    # TF-UPGRADE-TODO: In Terraform v0.10 and earlier, it was sometimes necessary to
    # force an interpolation expression to be interpreted as a list by wrapping it
    # in an extra set of list brackets. That form was supported for compatibilty in
    # v0.11, but is no longer supported in Terraform v0.12.
    #
    # If the expression in the following list itself returns a list, remove the
    # brackets to avoid interpretation as a list of lists. If the expression
    # returns a single list item then leave it as-is and remove this TODO comment.
    resources = [
      module.slack_user_query_lambda.function_arn,
    ]
  }
}
