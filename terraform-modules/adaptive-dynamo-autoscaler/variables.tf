variable "table_name" {
  description = "Name of the table for auto-scaling"
  type        = string
}

variable "table_arn" {
  description = "ARN of the table for auto-scaling"
  type        = string
}

variable "min_read_capacity" {
  description = "Autoscaling minimum read capacity for the table"
  type        = string
  default     = 5
}

variable "max_read_capacity" {
  description = "Autoscaling maximum read capacity for the table"
  type        = string
  default     = 10
}

variable "min_write_capacity" {
  description = "Autoscaling minimum write capacity for the table"
  type        = string
  default     = 5
}

variable "max_write_capacity" {
  description = "Autoscaling maximum write capacity for the table"
  type        = string
  default     = 10
}

variable "target_tracking_read_threshold" {
  description = "Threshold value of read capacity beyind which scaling should be triggered"
  type        = string
  default     = 70
}

variable "target_tracking_write_threshold" {
  description = "Threshold value of write capacity beyind which scaling should be triggered"
  type        = string
  default     = 70
}

variable "client_id" {
  description = "Unique id for the infrastructure"
  type        = string
}