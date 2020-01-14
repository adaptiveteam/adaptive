resource "aws_lambda_function" "lambda" {
  count = !var.attach_vpc_config && !var.attach_dl_config ? 1 : 0

  filename                       = var.filename
  description                    = var.description
  function_name                  = "${var.client_id}_${var.function_name}"
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

  //  dynamic "environment" {
  //    for_each = var.environment_variables == null ? {} : var.environment_variables
  //    content {
  //      variables = var.environment_variables
  //    }
  //  }
}

resource "aws_lambda_function" "lambda_with_dl" {
  count = var.attach_dl_config && ! var.attach_vpc_config ? 1 : 0

  dead_letter_config {
    target_arn = var.dl_config["target_arn"]
  }

  filename                       = var.filename
  description                    = var.description
  function_name                  = "${var.client_id}_${var.function_name}"
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

  environment {
    variables = var.environment_variables
  }
}

resource "aws_lambda_function" "lambda_with_vpc" {
  count = var.attach_vpc_config && ! var.attach_dl_config ? 1 : 0

  vpc_config {
    security_group_ids = [
      var.vpc_config["security_group_ids"]]
    subnet_ids         = [
      var.vpc_config["subnet_ids"]]
  }

  filename                       = var.filename
  description                    = var.description
  function_name                  = "${var.client_id}_${var.function_name}"
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

  environment {
    variables = var.environment_variables
  }
}

resource "aws_lambda_function" "lambda_with_dl_vpc" {
  count = var.attach_dl_config && var.attach_vpc_config ? 1 : 0

  dead_letter_config {
    target_arn = var.dl_config["target_arn"]
  }

  vpc_config {
    security_group_ids = [
      var.vpc_config["security_group_ids"]]
    subnet_ids         = [
      var.vpc_config["subnet_ids"]]
  }

  filename                       = var.filename
  description                    = var.description
  function_name                  = "${var.client_id}_${var.function_name}"
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

  environment {
    variables = var.environment_variables
  }
}
