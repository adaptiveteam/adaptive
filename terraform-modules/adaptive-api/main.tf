resource "aws_api_gateway_rest_api" "api" {
  name = "${var.client_id}_${var.name}"
}

resource "aws_api_gateway_account" "demo" {
  cloudwatch_role_arn = var.cloudwatch_role_arn
}

resource "aws_api_gateway_resource" "proxy" {
  parent_id   = aws_api_gateway_rest_api.api.root_resource_id
  path_part   = "{proxy+}"
  rest_api_id = aws_api_gateway_rest_api.api.id
}

resource "aws_api_gateway_method" "request_method" {
  authorization = "NONE"
  http_method   = var.method
  resource_id   = aws_api_gateway_resource.proxy.id
  rest_api_id   = aws_api_gateway_rest_api.api.id
}

resource "aws_api_gateway_integration" "request_method_integration" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  resource_id = aws_api_gateway_resource.proxy.id
  http_method = aws_api_gateway_method.request_method.http_method
  type        = "AWS_PROXY"
  uri         = "arn:aws:apigateway:${var.region}:lambda:path/2015-03-31/functions/${var.lambda_arn}/invocations"

  # AWS lambdas can only be invoked with the POST method
  integration_http_method = "POST"
}

resource "aws_api_gateway_deployment" "deployment" {
  rest_api_id       = aws_api_gateway_rest_api.api.id
  stage_name        = var.stage
  stage_description = "Deployment"
  description       = "Deployment at ${timestamp()}"
  depends_on        = [
    aws_api_gateway_integration.request_method_integration,
    aws_api_gateway_integration_response.response_method_integration
  ]
}

## GET response from lambda
resource "aws_api_gateway_method_response" "response_method" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  resource_id = aws_api_gateway_resource.proxy.id
  http_method = aws_api_gateway_integration.request_method_integration.http_method
  status_code = "200"

  response_models = {
    "application/json" = "Empty"
  }
}

resource "aws_api_gateway_integration_response" "response_method_integration" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  resource_id = aws_api_gateway_resource.proxy.id
  http_method = aws_api_gateway_method_response.response_method.http_method
  status_code = aws_api_gateway_method_response.response_method.status_code

  response_templates = {
    "application/json" = ""
  }
}

resource "aws_api_gateway_method_settings" "settings" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  stage_name  = var.stage
  method_path = "${aws_api_gateway_resource.proxy.path_part}/${aws_api_gateway_method.request_method.http_method}"

  settings {
    metrics_enabled = true
    logging_level   = "INFO"
  }

  depends_on = [aws_api_gateway_deployment.deployment]
}

resource "aws_lambda_permission" "allow_api_gateway" {
  function_name = var.lambda_arn
  statement_id  = "${aws_api_gateway_rest_api.api.name}-AllowExecutionFromApiGateway"
  action        = "lambda:InvokeFunction"
  principal     = "apigateway.amazonaws.com"
  source_arn    = "arn:aws:execute-api:${var.region}:${data.aws_caller_identity.current.account_id}:${aws_api_gateway_rest_api.api.id}/*/${var.method}${aws_api_gateway_resource.proxy.path}"
  depends_on    = [aws_api_gateway_rest_api.api, aws_api_gateway_resource.proxy]
}
