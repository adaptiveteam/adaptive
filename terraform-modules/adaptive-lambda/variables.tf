variable "description" {
  description = "Description of the Lambda function"
  type        = string
  default     = "Managed by Terraform"
}

variable "filename" {
  description = "Path to the lambda zip file"
}

variable "source_hash" {
  description = "Hash of the source lambda zip"
}

variable "function_name_suffix" {
  description = "Function name of the lambda"
}

locals {
  function_name = "${var.client_id}_${var.function_name_suffix}"
}

variable "handler" {
  description = "Handler for the lambda function"
}

variable "runtime" {
  description = "Runtime for lambda"
}

variable "timeout" {
  description = "Timeout for lambda function"
}

variable "environment_variables" {
  description = "Environment variables configuration for lambda function"
  type        = map
  default     = null
}

variable "memory_size" {
  description = "Amount of memory in MB to be allocated to lambda function"
  type        = string
  default     = 128
}

variable "attach_vpc_config" {
  description = "VPC configuration for lambda function"
  type        = string
  default     = false
}

variable "attach_dl_config" {
  description = "Flag to configure dead letter for lambda function"
  type        = string
  default     = false
}

variable "reserved_concurrent_executions" {
  description = "The amount of reserved concurrent executions for the Lambda function"
  type        = string
  default     = 0
}

variable "tags" {
  description = "A mapping of tags"
  type        = map
  default     = {}
}

variable "dl_config" {
  description = "Dead letter configuration for the Lambda function"
  type        = map
  default     = {}
}

variable "vpc_config" {
  description = "VPC configuration for the Lambda function"
  type        = map
  default     = {}
}

variable "enable_cloudwatch_logs" {
  description = "Flag to enable creation of cloudwatch logs"
  type        = string
  default     = true
}

variable "attach_policy" {
  description = "Flag to attach additional policy"
  type        = string
  default     = false
}

variable "policy" {
  description = "An addional policy to attach to the Lambda function"
  type        = string
  default     = ""
}

variable "schedule" {
  description = "Flag indicating whether lambda run should be scheduled"
  type        = string
  default     = false
}

variable "schedule_name" {
  description = "Name of the schedule rule"
  type        = string
  default     = ""
}

variable "schedule_description" {
  description = "Description of the schedule rule"
  type        = string
  default     = ""
}

variable "schedule_expression" {
  description = "Schedule expression for scheduled lambda"
  type        = string
  default     = ""
}

variable "schedule_invoke_json" {
  description = "JSON that should be passed when cloudwatch invokes the lambda based on a rule"
  type        = string
  default     = "{}"
}

variable "schedule_is_enabled" {
  type        = string
  description = "Schedule Rule is enabled (true/false)"
  default     = "true"
}

variable "client_id" {
  description = "Unique id for the infrastructure"
  type        = string
}

variable "lambda_tracing_mode" {
  description = "AWS X-ray tracing mode for lambda functions. Can be either PassThrough or Active."
  type        = string
  default     = "Active"
}