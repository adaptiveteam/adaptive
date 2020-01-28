locals {
  user_profile_lambda_function_name_suffix = "user-profile-lambda-go"
  user_profile_lambda_function_name = "${var.client_id}_${local.user_profile_lambda_function_name_suffix}"
}

module "user_profile_lambda" {
  source = "../../../terraform-modules/adaptive-lambda"

  client_id     = var.client_id
  filename      = data.archive_file.adaptive-lambda-zip.output_path
  source_hash   = data.archive_file.adaptive-lambda-zip.output_base64sha256
  handler       = "adaptive"
  function_name_suffix = local.user_profile_lambda_function_name_suffix
  runtime       = var.lambda_runtime
  timeout       = var.lambda_timeout
  memory_size   = 1536

  reserved_concurrent_executions = -1

  // Add environment variables.
  environment_variables = merge(local.environment_variables, {
    LAMBDA_ROLE   = "user-profile"
    LOG_NAMESPACE = "user-profile"
  })

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

