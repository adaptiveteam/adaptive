variable "aws_region" {
  description = "The AWS region to create things in."
}

variable "profile" {
  description = "Named AWS profile to use for deployment"
}

variable "environment" {
  description = "Environment"
  type = string
}

variable "remote_state_bucket" {
  description = "S3 bucket for terraform remote state storage"
  default = "adaptive-core-infra-remote-state"
}

variable "remote_state_table" {
  description = "Dynamo table for terraform remote state locking"
  default = "adaptive-core-infra-remote-state"
}

variable "remote_state_key" {
  description = "Key for remote state file"
  default = "terraform.tfstate"
}

variable "client_id" {
  description = "Unique id associated with a client"
}

variable "lambda_timeout" {
  description = "Lambda timeout"
  default = 30
}

variable "multi_core_memory_size" {
  description = "Amount of memory in MB to be allocated to lambda function to use multi-cores"
  type = string
  default = 1856
}

variable "lambda_runtime" {
  description = "Runtime for Lambda"
  default = "go1.x"
}

// https://github.com/terraform-providers/terraform-provider-aws/issues/6632#issuecomment-445323845
variable "dynamo_ondemand_read_capacity" {
  description = "GSI on-demadnd read capacity for dynamo table"
  default = 0
}

variable "dynamo_ondemand_write_capacity" {
  description = "GSI on-demand write capacity for dynamo table"
  default = 0
}

variable "dynamo_stream_view_type" {
  description = "View of dynamo streams"
  default = "NEW_AND_OLD_IMAGES"
}

variable "user_engagement_answered_dynamo_index" {
  description = "GSI for user with answered engagements"
  default = "UserIDAnsweredIndex"
}

# DAX

variable "dax_cluster_name" {
  description = "DAX cluster name"
  type = string
  default = "hoger-main"
}

variable "dax_end_point" {
  description = "DAX end point"
  default = "hoger-main.brhczd.clustercfg.dax.use1.cache.amazonaws.com:8111"
}

### API
variable "gateway_name" {
  description = "Name of the API gateway"
  default = "adaptive-gateway"
}

variable "gateway_stage" {
  description = "Stage for API gateway"
  default = "dev"
}

variable "gateway_http_method" {
  description = "HTTP method for the gateway"
  default = "POST" // We accept GET and POST. GET is used for OAuth redirect url
}

variable "gateway_global_cloudwatch_role" {
  description = "Allows API Gateway to push logs to CloudWatch Logs"
  default = "apigateway-logs"
}

## Coaching
variable "dynamo_users_platform_index" {
  description = "GSI for users on platform id"
  default = "PlatformIDIndex"
}

variable "dynamo_users_timezone_offset_index" {
  description = "GSI for users on platform with timezone offset"
  default = "PlatformIDTimezoneOffsetIndex"
}

variable "dynamo_users_scheduled_time_index" {
  description = "GSI for users on platform scheduled time with Adaptive"
  default = "PlatformIDAdaptiveScheduledTimeInUTCIndex"
}

variable "dynamo_coaching_relationship_coach_index" {
  description = "GSI for coaching relationship coach-quarter-year"
  default = "CoachQuarterYearIndex"
}

variable "dynamo_coaching_relationship_coachee_index" {
  description = "GSI for coaching relationship coachee-quarter-year"
  default = "CoacheeQuarterYearIndex"
}

variable "dynamo_coaching_relationship_quarter_year_index" {
  description = "GSI for coaching relationship quarter year"
  default = "QuarterYearIndex"
}

variable "dynamo_community_users_channel_index" {
  description = "GSI for community-users index by channel"
  default = "ChannelIDIndex"
}

variable "dynamo_community_users_user_community_index" {
  description = "GSI for community-users index by user and community"
  default = "UserIDCommunityIDIndex"
}

variable "dynamo_community_users_user_index" {
  description = "GSI for community-users index by user"
  default = "UserIDIndex"
}

variable "dynamo_community_users_community_index" {
  description = "GSI for community-users index by community"
  default = "PlatformIDCommunityIDIndex"
}

variable "dynamo_dialog_content_contect_subject_index" {
  description = "GSI for dailog content by context-subject"
  default = "ContextSubjectIndex"
}

variable "user_engagement_scheduler_lambda_prefix" {
  default = "user-engagement-scheduler-lambda-go"
}

## User objectives
# TODO: remove the following variable - dynamo_user_objectives_user_index. It's not used to create the index. The actual index name is defined in daos.
variable "dynamo_user_objectives_user_index" {
  description = "GSI for user objectives user index"
  default = "UserIDCompletedIndex"
}

# TODO: remove the following variable - dynamo_user_objectives_id_index. It's not used to create the index. There is no such index
variable "dynamo_user_objectives_id_index" {
  description = "GSI for user objectives id index"
  default = "IDIndex"
}

# TODO: remove the following variable - dynamo_user_objectives_partner_index. It's not used to create the index. The actual index name is defined in daos.
variable "dynamo_user_objectives_partner_index" {
  description = "GSI for user objectives partner index"
  default = "AccountabilityPartnerIndex"
}

# TODO: remove the following variable - dynamo_user_objectives_accepted_index. It's not used to create the index. The actual index name is defined in daos.
variable "dynamo_user_objectives_accepted_index" {
  description = "GSI for user objectives accepted index"
  default = "AcceptedIndex"
}

# TODO: remove the following variable - dynamo_user_objectives_type_index. It's not used to create the index. The actual index name is defined in daos.
variable "dynamo_user_objectives_type_index" {
  description = "GSI for user objectives based on type (individual/strategy)"
  default = "UserIDTypeIndex"
}
# This variable was updated accordingly to daos._user_objective_progress
variable "dynamo_user_objectives_progress_index" {
  description = "GSI for user objectives progress index"
  default = "IDIndex"
}
# This variable was updated accordingly to daos._user_objective_progress
variable "dynamo_user_objectives_progress_created_on_index" {
  description = "GSI for user objectives progress created on index"
  default = "CreatedOnIndex"
}

variable "dynamo_strategy_communities_platform_channel_created_index" {
  description = "Index for strategy community platform-channel index"
  default = "PlatformIDChannelCreatedIndex"
}

variable "dynamo_strategy_communities_platform_index" {
  default = "PlatformIDIndex"
}

variable "dynamo_strategy_communities_channel_index" {
  default = "ChannelIDIndex"
}

# strategy communities
variable "dynamo_capability_communities_platform_index" {
  description = "Index for strategy objectives based on platform id"
  default = "PlatformIDIndex"
}

variable "dynamo_strategy_initiative_communities_platform_index" {
  description = "Index for strategy initiatives based on platform id"
  default = "PlatformIDIndex"
}

# strategy initiatives
variable "dynamo_strategy_initiatives_platform_index" {
  description = "Index for strategy initiatives based on platform id"
  default = "PlatformIDIndex"
}

variable "dynamo_strategy_initiatives_initiative_community_index" {
  description = "Index for strategy initiatives based on initiative community id"
  default = "InitiativeCommunityIDIndex"
}

# strategy objectives
variable "dynamo_strategy_objectives_platform_index" {
  description = "Index for strategy objectives based on platform id"
  default = "PlatformIDIndex"
}

variable "dynamo_strategy_objectives_capability_community_index" {
  description = "Index for strategy objectives based on capabbility community id"
  default = "CapabilityCommunityIDsIndex"
}

# values
# TODO: remove the following variable - dynamo_adaptive_values_platform_id_index. It's not used to create the index. The actual index name is defined in daos.
variable "dynamo_adaptive_values_platform_id_index" {
  description = "GSI for adaptive values dynamo_adaptive_values_platform_id_index index"
  default = "PlatformIDIndex"
}

# User communities
variable "user_community_channel_dynamo_index" {
  description = "GSI for user communities with channel"
  default = "ChannelIndex"
}

variable "user_community_platform_dynamo_index" {
  description = "GSI for user communities with platform"
  default = "PlatformIDIndex"
}

# Holidays
variable "dynamo_holidays_date_index" {
  description = "GSI for holidays date index"
  default = "PlatformIDDateIndex"
}
# TODO: Remove this variable. The index does not exist and is not being used
variable "dynamo_holidays_id_index" {
  description = "GSI for HW id index"
  default = "HolidaysIdIndex"
}

# User feedback
variable "feedback_source_quarter_year_index" {
  description = "GSI for source with quarter year"
  default = "QuarterYearSourceIndex"
}
# TODO: remove this variable 'feedback_target_quarter_year_index'. It's not being referenced anywhere
variable "feedback_target_quarter_year_index" {
  description = "GSI for target with quarter year"
  default = "QuarterYearTargetIndex"
}

variable "RDS_USER" {
  default = "user"
}
variable "RDS_PASSWORD" {}
variable "RDS_PORT" {
  default = "3306"
}
variable "RDS_DB_NAME" {
  default = "test_report"
}

variable "SLACK_CLIENT_ID" {
  default = ""
}
variable "SLACK_CLIENT_SECRET" {
    default = ""
}
variable "SLACK_SIGNING_SECRET" {
    default = ""
}

# VPC
variable "vpc_cidr" {
  description = "VPC CIDR block"
  type = string
  default = "10.0.0.0/16"
}

variable "vpc_private_subnets" {
  description = "Private subnets"
  type = list
  default = [
    "10.0.1.0/24",
    "10.0.2.0/24",
  ]
}

variable "vpc_public_subnets" {
  description = "Public subnets"
  type = list
  default = [
    "10.0.101.0/24",
    "10.0.102.0/24",
  ]
}

variable "vpc_database_subnets" {
  default = [
    "10.0.3.0/24",
    "10.0.4.0/24",
  ]
}

variable "aws_availability_zones" {
  description = "AWS availability zones"
  type = map(list(string))
  default = {
    //  N. Virginia
    us-east-1 = [
      "us-east-1a",
      "us-east-1b",
      "us-east-1c",
    ]
    //  Ohio
    us-east-2 = [
      "us-east-2a",
      "us-east-2b",
      "us-east-2c",
    ]
    //  N. California
    us-west-1 = [
      "us-west-1a",
      "us-west-1b",
      "us-west-1c",
    ]
    //  Oregon
    us-west-2 = [
      "us-west-2a",
      "us-west-2b",
      "us-west-2c",
    ]
  }
}

locals {
    availability_zones = var.aws_availability_zones[var.aws_region]
}
