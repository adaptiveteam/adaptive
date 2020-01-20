output "http_method" {
  value = aws_api_gateway_integration_response.response_method_integration.http_method
}

output "api_url" {
  value = aws_api_gateway_deployment.deployment.invoke_url
}

output "api_id" {
  value = aws_api_gateway_rest_api.api.id
}

output "api_stage_name" {
  value = aws_api_gateway_deployment.deployment.stage_name
}