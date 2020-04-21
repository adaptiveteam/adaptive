resource "aws_lambda_function" "lambda" {
  count = 1

  s3_bucket                      = var.s3_bucket
  s3_key                         = var.s3_key
  description                    = var.description
  function_name                  = local.function_name
  role                           = aws_iam_role.lambda.arn
  memory_size                    = var.memory_size
  handler                        = var.handler
  source_code_hash               = var.source_hash
  runtime                        = var.runtime
  timeout                        = var.timeout
  reserved_concurrent_executions = var.reserved_concurrent_executions
  tags                           = var.tags

  tracing_config {
    mode = var.lambda_tracing_mode
  }

  # The aws_lambda_function resource has a schema for the environment
  # variable, where the only acceptable values are:
  #   a. Undefined
  #   b. An empty list
  #   c. A list containing 1 element: a map with a specific schema
  # Use slice to get option "b" or "c" depending on whether a non-empty
  # value was passed into this module.

  environment {
    variables = var.environment_variables
  }

}
