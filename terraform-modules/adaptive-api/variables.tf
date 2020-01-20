variable "name" {
  description = "Name of the REST API"
}

variable "stage" {
  description = "Stage name for API deployment (dev/test/staging/prod)"
}

variable "method" {
  description = "HTTP method"
  default     = "GET"
}

variable "lambda_name" {
  description = "Name of the lambda function to invoke"
}

variable "lambda_arn" {
  description = "ARN of the lambda function to invoke"
}

variable "region" {
  description = "AWS region"
}

variable "client_id" {
  description = "Unique id for the client"
}

variable "cloudwatch_role_arn" {
  description = "Cloudwatch role ARN for API gateway"
}