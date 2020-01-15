variable "id" {
  description = "Module instance id. It'll be prefixed to all global names"
}

variable "name" {
  description = "Name of lambda"
}

variable "modpath" {
  description = "Local relative path to a folder that contains golang module"
}

variable "memory_size" {
  description = "Amount of memory in MB to be allocated to lambda function"
  type        = string
  default     = 128
}

variable "runtime" {
  description = "Runtime for lambda"
  default     = "go1.x"
}

variable "timeout" {
  description = "Timeout for lambda function"
  default     = 30
}

variable "reserved_concurrent_executions" {
  description = "The amount of reserved concurrent executions for the Lambda function"
  type        = string
  default     = -1
}

variable "tags" {
  description = "A mapping of tags"
  type        = map
  default     = {}
}

variable "input_sns_topic_arn" {
  description = "ARN of SNS topic with input data for lambda"
  type        = string
}

variable "environment_variables" {
  description = "Environment configuration for lambda function"
  type        = map
  default     = {}
}

variable "go_get_url" {
  description = "Environment configuration for lambda function"
  type        = map
  default     = {}
}

variable "enable_cloudwatch_logs" {
  description = "Flag to enable creation of cloudwatch logs"
  type        = string
  default     = true
}

