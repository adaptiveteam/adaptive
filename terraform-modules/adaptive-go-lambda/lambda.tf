resource "aws_lambda_function" "lambda" {
  # count = "${!var.attach_vpc_config && !var.attach_dl_config ? 1 : 0}"

  filename         = "${var.modpath}${var.name}.zip"
  description      = "${var.name} lambda"
  function_name    = "${local.prefix_name}_lambda"
  source_code_hash = base64sha256("${var.modpath}/${var.name}.zip")
  // triggers { rerun = "${base64sha256(file("${var.modpath}/${var.name}.zip"))}" }
  handler          = var.name

  memory_size = var.memory_size
  runtime     = var.runtime
  timeout     = var.timeout

  reserved_concurrent_executions = var.reserved_concurrent_executions
  tags                           = var.tags

  role = aws_iam_role.lambda.arn

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

locals {
  function_arn = aws_lambda_function.lambda.*.arn[0]
  #function_arn = "${element(concat(aws_lambda_function.lambda.*.arn, aws_lambda_function.lambda_with_dl.*.arn, aws_lambda_function.lambda_with_vpc.*.arn, aws_lambda_function.lambda_with_dl_vpc.*.arn), 0)}"

  function_name = aws_lambda_function.lambda.*.function_name[0]
  //  function_name = "${element(concat(aws_lambda_function.lambda.*.function_name, aws_lambda_function.lambda_with_dl.*.function_name, aws_lambda_function.lambda_with_vpc.*.function_name, aws_lambda_function.lambda_with_dl_vpc.*.function_name), 0)}"
  prefix_name   = "${var.id}_${var.name}"
}
