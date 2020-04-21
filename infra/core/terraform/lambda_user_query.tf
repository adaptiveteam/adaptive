module "user_query_lambda" {
  source = "../../../terraform-modules/adaptive-lambda-s3"
  s3_bucket = aws_s3_bucket.binary_bucket.bucket
  s3_key = aws_s3_bucket_object.adaptive_zip.key
  source_hash = data.archive_file.adaptive_lambda_zip.output_md5

  client_id     = var.client_id
  handler       = "adaptive"
  function_name_suffix = "user-query-lambda-go"
  runtime       = var.lambda_runtime
  timeout       = var.lambda_timeout
  memory_size   = var.multi_core_memory_size

  reserved_concurrent_executions = -1

  // Add environment variables.
  environment_variables = merge(local.environment_variables, {
    LAMBDA_ROLE   = "user-query"
    LOG_NAMESPACE = "user-query"
  })

  // Attach extra policy
  attach_policy = true
  policy        = data.aws_iam_policy_document.user_query_policy.json

  // Schedule the lambda
  schedule             = true
  schedule_name        = "1pm_UTC_rule"
  schedule_description = "Cloudwatch event rule for 1PM UTC/ 8AM EST/ 6:30PM IST"
  schedule_expression  = "cron(0 13 * * ? *)"
  # schedule_expression  = "cron(0/15 * * * ? *)" 

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

