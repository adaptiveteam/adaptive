data "archive_file" "user-profile-lambda-zip" {
  type        = "zip"
  source_file = "../../../bin/user-profile-lambda-go"
  output_path = "lambdas/user-profile-lambda-go.zip"
}

module "user_profile_lambda" {
  source = "../../../terraform-modules/adaptive-lambda"

  client_id     = var.client_id
  filename      = data.archive_file.user-profile-lambda-zip.output_path
  source_hash   = data.archive_file.user-profile-lambda-zip.output_base64sha256
  function_name = "user-profile-lambda-go"
  handler       = "user-profile-lambda-go"
  runtime       = var.lambda_runtime
  timeout       = var.lambda_timeout
  memory_size   = 1536

  reserved_concurrent_executions = -1

  environment_variables = {
    CLIENT_ID                = var.client_id
    USER_TABLE_NAME          = aws_dynamodb_table.adaptive_users_dynamodb_table.name
    CLIENT_CONFIG_TABLE_NAME = aws_dynamodb_table.client_config_dynamodb_table.name
    DAX_ENDPOINT             = var.dax_end_point
    LOG_NAMESPACE            = "user-profile"
  }

  // Attach extra policy
  attach_policy = true
  policy        = data.aws_iam_policy_document.user_profile_policy.json

  // Schedule the lambda
  schedule             = true
  schedule_name        = "user_profile_every_5_min"
  schedule_description = "User Profile lambda cloudwatch event rule for every 5 min"
  schedule_expression  = "rate(5 minutes)"
  schedule_invoke_json = "{\"user_id\" : \"\"}"

  tags = local.default_tags
}

data "aws_iam_policy_document" "user_profile_policy" {
  statement {
    effect    = "Allow"
    actions   = [
      "dynamodb:GetItem",
      "dynamodb:PutItem",
    ]
    resources = [
      aws_dynamodb_table.adaptive_users_dynamodb_table.arn,
    ]
  }

  statement {
    effect    = "Allow"
    actions   = [
      "dynamodb:GetItem",
    ]
    resources = [
      aws_dynamodb_table.client_config_dynamodb_table.arn,
    ]
  }

  # DAX access
  statement {
    effect    = "Allow"
    actions   = [
      "dax:GetItem",
      "dax:Query",
      "dax:Endpoints",
      "dax:*",
    ]
    resources = [
      "arn:aws:dax:${var.aws_region}:${data.aws_caller_identity.current.account_id}:cache/${var.dax_cluster_name}",
    ]
  }
}

