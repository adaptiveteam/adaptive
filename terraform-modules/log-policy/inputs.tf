variable "function_name" {
    description = "The name of lambda function"
    type        = string
}

variable "client_id" {
  description = "Unique id for the infrastructure"
  type        = string
}

variable "errors_sns_topic_arn" {
  description = "aws_sns_topic.errors.arn"
  type        = string
}

variable "role_name" {
  description = "role name"
  type        = string
}
