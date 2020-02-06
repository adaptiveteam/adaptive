output "client_id" {
  description = "Unique id for the client"
  value = var.client_id
}

output "region" {
  description = "AWS region for the deployment"
  value = var.aws_region
}

output "client_config_table_name" {
  description = "Name of the client configuration table"
  value = aws_dynamodb_table.client_config_dynamodb_table.name
}

output "client_config_table_arn" {
  description = "Client configuration table ARN"
  value = aws_dynamodb_table.client_config_dynamodb_table.arn
}

output "user_query_lambda_name" {
  description = "Name of user query lambda"
  value = module.user_query_lambda.function_name
}

output "users_table_name" {
  description = "Name of adaptive users table"
  value = aws_dynamodb_table.adaptive_users_dynamodb_table.name
}

output "users_table_arn" {
  description = "ARN of adaptive users table"
  value = aws_dynamodb_table.adaptive_users_dynamodb_table.arn
}

output "dynamo_users_platform_index" {
  value = var.dynamo_users_platform_index
}

output "user_engagements_table_name" {
  description = "Name of user engagements table"
  value = aws_dynamodb_table.adaptive_user_engagements_dynamo_table.name
}

output "user_engagements_table_arn" {
  description = "Name of user engagements table"
  value = aws_dynamodb_table.adaptive_user_engagements_dynamo_table.arn
}

output "user_profile_lambda_name" {
  description = "Name of the user profile lambda"
  value = module.user_profile_lambda.function_name
}

output "user_setup_lambda_name" {
  description = "Name of the user setup lambda"
  value = module.user_setup_lambda.function_name
}

output "user_engagement_scripting_lambda_name" {
  description = "Name of the user engagement scripting lambda"
  value = module.user_engagement_scripting_lambda.function_name
}

output "errors_sns_topic_name" {
  description = "SNS topic name for errors"
  value = aws_sns_topic.errors.name
}

output "namespace_payload_sns_topic_name" {
  description = "SNS topic ARN for namespace payload"
  value = aws_sns_topic.namespace_payload.name
}

output "platform_notification_topic_arn" {
  description = "SNS topic ARN for platform notification"
  value = aws_sns_topic.platform_notification.arn
}

output "api_url" {
  description = "Gateway URL for Adaptive"
  value = module.api.api_url
}

output "api_gateway_id" {
  description = "ID of the API Gateway"
  value = module.api.api_id
}

output "api_gateway_stage" {
  description = "Stage name of the API Gateway"
  value = module.api.api_stage_name
}

## Coaching
output "coaching_relationships_dynamo_table_name" {
  description = "Coaching relationships dynamo table name"
  value = aws_dynamodb_table.coaching_relationships.name
}

output "coaching_relationships_dynamo_table_arn" {
  description = "Coaching relationships dynamo table ARN"
  value = aws_dynamodb_table.coaching_relationships.arn
}

output "coaching_rejections_dynamo_table_name" {
  description = "Coaching rejections dynamo table name"
  value = aws_dynamodb_table.coaching_rejections.name
}

output "coaching_rejections_dynamo_table_arn" {
  description = "Coaching rejections dynamo table ARN"
  value = aws_dynamodb_table.coaching_rejections.arn
}

output "dynamo_coaching_relationship_coach_index" {
  description = "GSI for coaching relationship coach"
  value = var.dynamo_coaching_relationship_coach_index
}

output "dynamo_coaching_relationship_coachee_index" {
  description = "GSI for coaching relationship coachee"
  value = var.dynamo_coaching_relationship_coachee_index
}

output "dynamo_coaching_relationship_quarter_year_index" {
  description = "GSI for coaching relationship quarter year"
  value = var.dynamo_coaching_relationship_quarter_year_index
}

output "dynamo_community_users_table_name" {
  description = "Community users table name"
  value = aws_dynamodb_table.community_users.name
}

output "dynamo_community_users_table_arn" {
  description = "Community users table ARN"
  value = aws_dynamodb_table.community_users.arn
}

output "dynamo_community_users_channel_index" {
  description = "GSI for community-users index by channel"
  value = var.dynamo_community_users_channel_index
}

output "dynamo_community_users_community_index" {
  description = "GSI for community-users index by channel"
  value = var.dynamo_community_users_community_index
}

output "environment_stage" {
  description = "Stage environment for Adaptive"
  value = var.environment
}

output "adaptive_dialog_table_name" {
  description = "Adaptive dialog content table"
  value = aws_dynamodb_table.adaptive_dialog_content.name
}

output "adaptive_dialog_table_arn" {
  description = "Adaptive dialog content table"
  value = aws_dynamodb_table.adaptive_dialog_content.arn
}

output "dynamo_dialog_content_contect_subject_index" {
  description = "GSI for dailog content by context-subject"
  value = var.dynamo_dialog_content_contect_subject_index
}

## Objectives
output "user_objectives_table_name" {
  value = aws_dynamodb_table.user_objective_dynamodb_table.name
}

output "user_objectives_table_arn" {
  value = aws_dynamodb_table.user_objective_dynamodb_table.arn
}

output "user_objectives_progress_table_name" {
  value = aws_dynamodb_table.user_objectives_progress.name
}

output "user_objectives_progress_table_arn" {
  value = aws_dynamodb_table.user_objectives_progress.arn
}

output "accountability_partnership_rejections_table_name" {
  value = aws_dynamodb_table.accountability_partnership_rejections_table.name
}

output "accountability_partnership_rejections_table_arn" {
  value = aws_dynamodb_table.accountability_partnership_rejections_table.arn
}

output "dynamo_user_objectives_user_index" {
  description = "GSI for user objectives user index"
  value = var.dynamo_user_objectives_user_index
}

output "dynamo_user_objectives_id_index" {
  description = "GSI for user objectives id index"
  value = var.dynamo_user_objectives_id_index
}

output "dynamo_user_objectives_accepted_index" {
  description = "GSI for user obejetives accepted index"
  value = var.dynamo_user_objectives_accepted_index
}

output "dynamo_user_objectives_partner_index" {
  description = "GSI for user obejetives partner index"
  value = var.dynamo_user_objectives_partner_index
}

output "dynamo_user_objectives_progress_index" {
  description = "GSI for user objctives progress index"
  value = var.dynamo_user_objectives_progress_index
}

output "dynamo_user_objectives_progress_created_on_index" {
  description = "GSI for user objctives progress created on index"
  value = var.dynamo_user_objectives_progress_created_on_index
}

output "dynamo_user_objectives_completed_index" {
  description = "GSI for user objectives index by completed"
  value = var.dynamo_community_users_user_index
}

output "dynamo_community_users_user_index" {
  description = "GSI for community-users index by user"
  value = var.dynamo_community_users_user_index
}

output "adaptive_vision_table_name" {
  value = aws_dynamodb_table.vision.name
}

output "adaptive_vision_table_arn" {
  value = aws_dynamodb_table.vision.arn
}

output "strategy_communities_table_name" {
  value = aws_dynamodb_table.strategy_communities.name
}

output "strategy_communities_table_arn" {
  value = aws_dynamodb_table.strategy_communities.arn
}

output "dynamo_strategy_communities_platform_channel_created_index" {
  description = "Index for strategy community platform-channel index"
  value = var.dynamo_strategy_communities_platform_channel_created_index
}

output "dynamo_strategy_initiative_communities_platform_index" {
  description = "Index for strategy initiatives based on platform id"
  value = "PlatformIDIndex"
}

output "dynamo_strategy_communities_channel_index" {
  value = "ChannelIDIndex"
}

# strategy communities
output "capability_communities_table_name" {
  value = aws_dynamodb_table.capability_communities.name
}

output "capability_communities_table_arn" {
  value = aws_dynamodb_table.capability_communities.arn
}

output "initiative_communities_table_name" {
  value = aws_dynamodb_table.initiative_communities.name
}

output "initiative_communities_table_arn" {
  value = aws_dynamodb_table.initiative_communities.arn
}

output "dynamo_capability_communities_platform_index" {
  description = "Index for strategy objectives based on platform id"
  value = var.dynamo_capability_communities_platform_index
}

output "dynamo_initiative_communities_platform_index" {
  description = "Index for strategy objectives based on platform id"
  value = var.dynamo_strategy_initiative_communities_platform_index
}

output "dynamo_strategy_communities_platform_index" {
  value = "PlatformIDIndex"
}

# strategy initiatives
output "strategy_initiatives_table_name" {
  value = aws_dynamodb_table.strategy_initiatives.name
}

output "strategy_initiatives_table_arn" {
  value = aws_dynamodb_table.strategy_initiatives.arn
}

output "dynamo_strategy_initiatives_platform_index" {
  description = "Index for strategy initiatives based on platform id"
  value = var.dynamo_strategy_initiatives_platform_index
}

output "dynamo_strategy_initiatives_initiative_community_index" {
  description = "Index for strategy initiatives based on initiative community id"
  value = var.dynamo_strategy_initiatives_initiative_community_index
}

# strategy objectives
output "strategy_objectives_table_name" {
  value = aws_dynamodb_table.strategy_objectives.name
}

output "strategy_objectives_table_arn" {
  value = aws_dynamodb_table.strategy_objectives.arn
}

output "dynamo_strategy_objectives_platform_index" {
  description = "Index for strategy objectives based on platform id"
  value = var.dynamo_strategy_objectives_platform_index
}

output "dynamo_strategy_objectives_capability_community_index" {
  description = "Index for strategy objectives based on capabbility community id"
  value = var.dynamo_strategy_objectives_capability_community_index
}

output "dynamo_user_objectives_type_index" {
  description = "GSI for user objectives based on type (individual/strategy)"
  value = var.dynamo_user_objectives_type_index
}

# values
output "values_table_name" {
  description = "Adaptive Values table name"
  value = aws_dynamodb_table.adaptive_value_dynamodb_table.name
}

output "values_table_arn" {
  description = "Adaptive Values table ARN"
  value = aws_dynamodb_table.adaptive_value_dynamodb_table.arn
}

output "dynamo_adaptive_values_platform_id_index" {
  description = "Adaptive Values table index on platform id"
  value = var.dynamo_adaptive_values_platform_id_index
}

output "dynamo_community_users_user_community_index" {
  description = "GSI for community-users index by user and community"
  value = var.dynamo_community_users_user_community_index
}

output "user_engagement_answered_dynamo_index" {
  description = "GSI for user with answered engagements"
  value = var.user_engagement_answered_dynamo_index
}

# User communities
output "user_communities_table_name" {
  value = aws_dynamodb_table.user_communities.name
}
output "user_communities_table_arn" {
  value = aws_dynamodb_table.user_communities.arn
}
output "user_community_channel_dynamo_index" {
  value = var.user_community_channel_dynamo_index
}
output "user_community_platform_dynamo_index" {
  value = var.user_community_platform_dynamo_index
}

#Holidays
output "holidays_table_name" {
  value = aws_dynamodb_table.ad_hoc_holidays.name
}
output "holidays_table_arn" {
  value = aws_dynamodb_table.ad_hoc_holidays.arn
}
output "dynamo_holidays_date_index" {
  value = var.dynamo_holidays_date_index
}
output "dynamo_holidays_id_index" {
  value = var.dynamo_holidays_id_index
}

# User Feedback
output "user_feedback_table_name" {
  value = aws_dynamodb_table.adaptive_user_feedback_dynamodb_table.name
}
output "user_feedback_table_arn" {
  value = aws_dynamodb_table.adaptive_user_feedback_dynamodb_table.arn
}
output "feedback_source_quarter_year_index" {
  value = var.feedback_source_quarter_year_index
}
output "feedback_target_quarter_year_index" {
  value = var.feedback_target_quarter_year_index
}

#Reporting
output "reporting_lambda_name" {
  value = local.reporting_lambda_name
}

output "engagement_scheduler_lambda_name" {
  value = module.user_engagement_scheduler_lambda.function_name
}

output "feedback_reports_bucket_name" {
  description = "Bucket name containing feedback reports"
  value = aws_s3_bucket.adaptive-feedback-reports-bucket.bucket
}

output "feedback_reports_bucket_arn" {
  description = "Bucket ARN for feedback reports"
  value = aws_s3_bucket.adaptive-feedback-reports-bucket.arn
}

output "adaptive_client_config_stream_arn" {
  value = aws_dynamodb_table.client_config_dynamodb_table.stream_arn
}

output "adaptive_users_stream_arn" {
  value = aws_dynamodb_table.adaptive_users_dynamodb_table.stream_arn
}

output "adaptive_vision_stream_arn" {
  value = aws_dynamodb_table.vision.stream_arn
}

output "adaptive_competency_stream_arn" {
  value = aws_dynamodb_table.adaptive_value_dynamodb_table.stream_arn
}

output "adaptive_user_feedback_stream_arn" {
  value = aws_dynamodb_table.adaptive_user_feedback_dynamodb_table.stream_arn
}

output "adaptive_user_objective_stream_arn" {
  value = aws_dynamodb_table.user_objective_dynamodb_table.stream_arn
}

output "adaptive_user_objective_progress_stream_arn" {
  value = aws_dynamodb_table.user_objectives_progress.stream_arn
}

output "adaptive_holidays_stream_arn" {
  value = aws_dynamodb_table.ad_hoc_holidays.stream_arn
}

output "adaptive_engagements_stream_arn" {
  value = aws_dynamodb_table.adaptive_user_engagements_dynamo_table.stream_arn
}

output "adaptive_communities_stream_arn" {
  value = aws_dynamodb_table.user_communities.stream_arn
}

output "adaptive_community_users_stream_arn" {
  value = aws_dynamodb_table.community_users.stream_arn
}

output "adaptive_initiative_communities_stream_arn" {
  value = aws_dynamodb_table.initiative_communities.stream_arn
}

output "adaptive_objectives_stream_arn" {
  value = aws_dynamodb_table.strategy_objectives.stream_arn
}

output "adaptive_objective_communities_stream_arn" {
  value = aws_dynamodb_table.strategy_communities.stream_arn
}

output "partnership_rejections_stream_arn" {
  value = aws_dynamodb_table.accountability_partnership_rejections_table.stream_arn
}

output "postponed_event_table_arn" {
	description = "ARN of the postponed_event table"
	value = aws_dynamodb_table.postponed_event_dynamodb_table.arn
}

output "postponed_event_table_name" {
	description = "Name of the postponed_event table"
	value = aws_dynamodb_table.postponed_event_dynamodb_table.name
}
