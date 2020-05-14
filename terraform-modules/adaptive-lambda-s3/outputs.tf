output "function_arn" {
  description = "ARN of the lambda function"
  value       = aws_lambda_function.lambda[0].arn
}

output "function_name" {
  description = "The name of the Lambda function"
  value       = aws_lambda_function.lambda[0].function_name
}

output "role_arn" {
  description = "ARN of the IAM role created for the Lambda function"
  value       = aws_iam_role.lambda.arn
}

output "role_name" {
  description = "Name of the IAM role created for the Lambda function"
  value       = aws_iam_role.lambda.name
}

output "log_group_name" {
  description = "Log group name for the lambda"
  value       = aws_cloudwatch_log_group.log_group.name
}

output "log_group_arn" {
  description = "Log group ARN for the lambda"
  value       = aws_cloudwatch_log_group.log_group.arn
}