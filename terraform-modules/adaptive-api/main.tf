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

# resource "aws_api_gateway_method" "POST_request_method" {
#   authorization = "NONE"
#   http_method   = "POST"
#   resource_id   = aws_api_gateway_resource.proxy.id
#   rest_api_id   = aws_api_gateway_rest_api.api.id
# }

resource "aws_api_gateway_method" "ANY_request_method" {
  authorization = "NONE"
  http_method   = "ANY"
  resource_id   = aws_api_gateway_resource.proxy.id
  rest_api_id   = aws_api_gateway_rest_api.api.id
}

# resource "aws_api_gateway_method" "GET_request_method" {
#   authorization = "NONE"
#   http_method   = "GET"
#   resource_id   = aws_api_gateway_resource.proxy.id
#   rest_api_id   = aws_api_gateway_rest_api.api.id
# }

# resource "aws_api_gateway_integration" "POST_request_method_integration" {
#   rest_api_id = aws_api_gateway_rest_api.api.id
#   resource_id = aws_api_gateway_resource.proxy.id
#   http_method = aws_api_gateway_method.POST_request_method.http_method
#   type        = "AWS_PROXY"
#   uri         = "arn:aws:apigateway:${var.region}:lambda:path/2015-03-31/functions/${var.lambda_arn}/invocations"

#   # AWS lambdas can only be invoked with the POST method
#   integration_http_method = "POST"
#}

resource "aws_api_gateway_integration" "ANY_request_method_integration" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  resource_id = aws_api_gateway_resource.proxy.id
  http_method = aws_api_gateway_method.ANY_request_method.http_method
  type        = "AWS_PROXY"
  uri         = "arn:aws:apigateway:${var.region}:lambda:path/2015-03-31/functions/${var.lambda_arn}/invocations"

  # AWS lambdas can only be invoked with the POST method
  integration_http_method = "POST"
}

# resource "aws_api_gateway_integration" "GET_request_method_integration" {
#   rest_api_id = aws_api_gateway_rest_api.api.id
#   resource_id = aws_api_gateway_resource.proxy.id
#   http_method = aws_api_gateway_method.GET_request_method.http_method
#   type        = "AWS_PROXY"
#   uri         = "arn:aws:apigateway:${var.region}:lambda:path/2015-03-31/functions/${var.lambda_arn}/invocations"

#   # AWS lambdas can only be invoked with the POST method
#   integration_http_method = "POST"
# }

resource "aws_api_gateway_deployment" "deployment" {
  rest_api_id       = aws_api_gateway_rest_api.api.id
  stage_name        = var.stage
  stage_description = "Deployment"
  description       = "Deployment at ${timestamp()}"
  depends_on        = [
    aws_api_gateway_integration.ANY_request_method_integration,
    # aws_api_gateway_integration.POST_request_method_integration,
    # aws_api_gateway_integration.GET_request_method_integration,
    aws_api_gateway_integration_response.response_method_integration_ANY
    # aws_api_gateway_integration_response.response_method_integration_POST,
    # aws_api_gateway_integration_response.response_method_integration_GET
  ]
}

resource "aws_api_gateway_method_response" "response_method_ANY" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  resource_id = aws_api_gateway_resource.proxy.id
  http_method = aws_api_gateway_integration.ANY_request_method_integration.http_method
  status_code = "200"

  response_models = {
    "application/json" = "Empty"
  }
}

# resource "aws_api_gateway_method_response" "response_method_POST" {
#   rest_api_id = aws_api_gateway_rest_api.api.id
#   resource_id = aws_api_gateway_resource.proxy.id
#   http_method = aws_api_gateway_integration.POST_request_method_integration.http_method
#   status_code = "200"

#   response_models = {
#     "application/json" = "Empty"
#   }
# }

# resource "aws_api_gateway_method_response" "response_method_GET" {
#   rest_api_id = aws_api_gateway_rest_api.api.id
#   resource_id = aws_api_gateway_resource.proxy.id
#   http_method = aws_api_gateway_integration.GET_request_method_integration.http_method
#   status_code = "200"

#   response_models = {
#     "application/json" = "Empty"
#   }
# }

resource "aws_api_gateway_integration_response" "response_method_integration_ANY" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  resource_id = aws_api_gateway_resource.proxy.id
  http_method = aws_api_gateway_method_response.response_method_ANY.http_method
  status_code = aws_api_gateway_method_response.response_method_ANY.status_code

  response_templates = {
    "application/json" = ""
  }
}

# resource "aws_api_gateway_integration_response" "response_method_integration_POST" {
#   rest_api_id = aws_api_gateway_rest_api.api.id
#   resource_id = aws_api_gateway_resource.proxy.id
#   http_method = aws_api_gateway_method_response.response_method_POST.http_method
#   status_code = aws_api_gateway_method_response.response_method_POST.status_code

#   response_templates = {
#     "application/json" = ""
#   }
# }

# resource "aws_api_gateway_integration_response" "response_method_integration_GET" {
#   rest_api_id = aws_api_gateway_rest_api.api.id
#   resource_id = aws_api_gateway_resource.proxy.id
#   http_method = aws_api_gateway_method_response.response_method_GET.http_method
#   status_code = aws_api_gateway_method_response.response_method_GET.status_code

#   response_templates = {
#     "application/json" = ""
#   }
# }

resource "aws_api_gateway_method_settings" "settings_ANY" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  stage_name  = var.stage
  method_path = "${aws_api_gateway_resource.proxy.path_part}/*" #${aws_api_gateway_method.ANY_request_method.http_method}"

  settings {
    metrics_enabled = true
    logging_level   = "INFO"
  }

  depends_on = [aws_api_gateway_deployment.deployment]
}

# resource "aws_api_gateway_method_settings" "settings_POST" {
#   rest_api_id = aws_api_gateway_rest_api.api.id
#   stage_name  = var.stage
#   method_path = "${aws_api_gateway_resource.proxy.path_part}/${aws_api_gateway_method.POST_request_method.http_method}"

#   settings {
#     metrics_enabled = true
#     logging_level   = "INFO"
#   }

#   depends_on = [aws_api_gateway_deployment.deployment]
# }

# resource "aws_api_gateway_method_settings" "settings_GET" {
#   rest_api_id = aws_api_gateway_rest_api.api.id
#   stage_name  = var.stage
#   method_path = "${aws_api_gateway_resource.proxy.path_part}/${aws_api_gateway_method.GET_request_method.http_method}"

#   settings {
#     metrics_enabled = true
#     logging_level   = "INFO"
#   }

#   depends_on = [aws_api_gateway_deployment.deployment]
# }

resource "aws_lambda_permission" "allow_api_gateway_ANY" {
  function_name = var.lambda_arn
  statement_id  = "${aws_api_gateway_rest_api.api.name}-AllowExecutionFromApiGatewayPOST"
  action        = "lambda:InvokeFunction"
  principal     = "apigateway.amazonaws.com"
  source_arn    = "arn:aws:execute-api:${var.region}:${data.aws_caller_identity.current.account_id}:${aws_api_gateway_rest_api.api.id}/*/ANY${aws_api_gateway_resource.proxy.path}"
  depends_on    = [aws_api_gateway_rest_api.api, aws_api_gateway_resource.proxy]
}
# resource "aws_lambda_permission" "allow_api_gateway_POST" {
#   function_name = var.lambda_arn
#   statement_id  = "${aws_api_gateway_rest_api.api.name}-AllowExecutionFromApiGatewayPOST"
#   action        = "lambda:InvokeFunction"
#   principal     = "apigateway.amazonaws.com"
#   source_arn    = "arn:aws:execute-api:${var.region}:${data.aws_caller_identity.current.account_id}:${aws_api_gateway_rest_api.api.id}/*/${var.method}${aws_api_gateway_resource.proxy.path}"
#   depends_on    = [aws_api_gateway_rest_api.api, aws_api_gateway_resource.proxy]
# }

# resource "aws_lambda_permission" "allow_api_gateway_GET" {
#   function_name = var.lambda_arn
#   statement_id  = "${aws_api_gateway_rest_api.api.name}-AllowExecutionFromApiGatewayGET"
#   action        = "lambda:InvokeFunction"
#   principal     = "apigateway.amazonaws.com"
#   source_arn    = "arn:aws:execute-api:${var.region}:${data.aws_caller_identity.current.account_id}:${aws_api_gateway_rest_api.api.id}/*/GET${aws_api_gateway_resource.proxy.path}"
#   depends_on    = [aws_api_gateway_rest_api.api, aws_api_gateway_resource.proxy]
# }
