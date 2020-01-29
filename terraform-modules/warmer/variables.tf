
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

variable "target_lambda_arn" {
  description = "The arn of lambda to trigger"
  type        = string
}
