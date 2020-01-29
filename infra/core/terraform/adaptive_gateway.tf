## API gateway
module "api" {
  source = "../../../terraform-modules/adaptive-api"

  client_id           = var.client_id
  name                = var.gateway_name
  stage               = var.gateway_stage
  method              = var.gateway_http_method
  lambda_name         = aws_lambda_function.slack_message_processor_lambda.function_name
  lambda_arn          = aws_lambda_function.slack_message_processor_lambda.arn
  region              = var.aws_region
  cloudwatch_role_arn = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/${var.gateway_global_cloudwatch_role}"
}
