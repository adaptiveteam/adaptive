locals {
  environment_variables = {

    ADAPTIVE_HELP_PAGE = "https://adaptiveteam.github.io/docs/general/commands"
    NAMESPACE_PAYLOAD_TOPIC_ARN = aws_sns_topic.namespace_payload.arn
    PLATFORM_NOTIFICATION_TOPIC = aws_sns_topic.platform_notification.arn
    USER_COMMUNITIES_TABLE = aws_dynamodb_table.user_communities.name
    COMMUNITY_USERS_TABLE_NAME = aws_dynamodb_table.community_users.name
    COMMUNITY_USERS_USER_COMMUNITY_INDEX = var.dynamo_community_users_user_community_index
    COMMUNITY_USERS_USER_INDEX = var.dynamo_community_users_user_index
    USER_OBJECTIVES_TABLE_NAME = aws_dynamodb_table.user_objective_dynamodb_table.name
    USER_OBJECTIVES_PARTNER_INDEX = var.dynamo_user_objectives_partner_index
    USER_OBJECTIVES_ID_INDEX = var.dynamo_user_objectives_id_index
    USER_OBJECTIVES_TYPE_INDEX = var.dynamo_user_objectives_type_index
    DIALOG_TABLE = aws_dynamodb_table.adaptive_dialog_content.name
    VISION_TABLE_NAME = aws_dynamodb_table.vision.name
    CAPABILITY_COMMUNITIES_TABLE_NAME = aws_dynamodb_table.capability_communities.name
    CAPABILITY_COMMUNITIES_PLATFORM_INDEX = var.dynamo_capability_communities_platform_index
    INITIATIVE_COMMUNITIES_TABLE_NAME = aws_dynamodb_table.initiative_communities.name
    INITIATIVE_COMMUNITIES_PLATFORM_INDEX = var.dynamo_strategy_initiative_communities_platform_index
    STRATEGY_COMMUNITIES_TABLE_NAME = aws_dynamodb_table.strategy_communities.name
    STRATEGY_COMMUNITIES_PLATFORM_CHANNEL_CREATED_INDEX = var.dynamo_strategy_communities_platform_channel_created_index
    // ADM
    USER_OBJECTIVES_TABLE = aws_dynamodb_table.user_objective_dynamodb_table.name
    USER_OBJECTIVES_USER_ID_INDEX = var.dynamo_user_objectives_user_index
    USER_OBJECTIVES_PROGRESS_TABLE = aws_dynamodb_table.user_objectives_progress.name
    USER_OBJECTIVES_PROGRESS_ID_INDEX = var.dynamo_user_objectives_progress_index
    CLIENT_ID = var.client_id
    ADAPTIVE_VALUES_TABLE = aws_dynamodb_table.adaptive_value_dynamodb_table.name
    USER_ENGAGEMENT_SCRIPTING_LAMBDA_NAME = module.user_engagement_scripting_lambda.function_name
    USER_ENGAGEMENTS_TABLE_NAME = aws_dynamodb_table.adaptive_user_engagements_dynamo_table.name
    USER_ANSWERED_INDEX = var.user_engagement_answered_dynamo_index
    STRATEGY_INITIATIVES_TABLE_NAME = aws_dynamodb_table.strategy_initiatives.name
    STRATEGY_INITIATIVES_PLATFORM_INDEX = var.dynamo_strategy_initiatives_platform_index
    STRATEGY_OBJECTIVES_TABLE_NAME = aws_dynamodb_table.strategy_objectives.name
    STRATEGY_OBJECTIVES_PLATFORM_INDEX = var.dynamo_strategy_objectives_platform_index
    STRATEGY_OBJECTIVES_CAPABILITY_COMMUNITY_INDEX = var.dynamo_strategy_objectives_capability_community_index
    STRATEGY_INITIATIVES_INITIATIVE_COMMUNITY_ID_INDEX = var.dynamo_strategy_initiatives_initiative_community_index

    USERS_TABLE_NAME = aws_dynamodb_table.adaptive_users_dynamodb_table.name
    USER_PROFILE_LAMBDA_NAME = module.user_profile_lambda.function_name

    // for adaptive-utils-go
    USERS_PLATFORM_INDEX = var.dynamo_users_platform_index
    COMMUNITY_USERS_COMMUNITY_INDEX = var.dynamo_community_users_community_index
    COMMUNITY_USERS_COMMUNITY_INDEX = var.dynamo_community_users_community_index
    COACHING_RELATIONSHIPS_TABLE_NAME = aws_dynamodb_table.coaching_relationships.name
    COACHING_RELATIONSHIPS_COACHEE_QUARTER_YEAR_INDEX = var.dynamo_coaching_relationship_coachee_index
    COACHING_RELATIONSHIPS_COACH_QUARTER_YEAR_INDEX = var.dynamo_coaching_relationship_coach_index
    COACHING_RELATIONSHIPS_QUARTER_YEAR_INDEX = var.dynamo_coaching_relationship_quarter_year_index
    ADAPTIVE_COMMUNITIES_TABLE = aws_dynamodb_table.user_communities.name
    USER_FEEDBACK_TABLE_NAME = aws_dynamodb_table.adaptive_user_feedback_dynamodb_table.name
    USER_FEEDBACK_SOURCE_QUARTER_YEAR_INDEX = var.feedback_source_quarter_year_index
    FEEDBACK_REPORTING_LAMBDA_NAME = local.reporting_lambda_name
    HOLIDAYS_AD_HOC_TABLE = aws_dynamodb_table.ad_hoc_holidays.name
    HOLIDAYS_PLATFORM_DATE_INDEX = var.dynamo_holidays_date_index
    STRATEGY_INITIATIVES_TABLE = aws_dynamodb_table.strategy_initiatives.name
    STRATEGY_INITIATIVES_PLATFORM_INDEX = var.dynamo_strategy_initiatives_platform_index
    STRATEGY_OBJECTIVES_TABLE = aws_dynamodb_table.strategy_objectives.name
    STRATEGY_OBJECTIVES_PLATFORM_INDEX = var.dynamo_strategy_objectives_platform_index
    VISION_TABLE_NAME = aws_dynamodb_table.vision.name
    SLACK_MESSAGE_PROCESSOR_SUFFIX = local.slack_message_processor_suffix

    # Reporting
    REPORTS_BUCKET_NAME = aws_s3_bucket.adaptive-feedback-reports-bucket.bucket

    RDS_HOST     = var.RDS_HOST
    RDS_USER     = var.RDS_USER
    RDS_PASSWORD = var.RDS_PASSWORD
    RDS_PORT     = var.RDS_PORT
    RDS_DB_NAME  = var.RDS_DB_NAME


    CLIENT_CONFIG_TABLE_NAME                = aws_dynamodb_table.client_config_dynamodb_table.name
    USER_COMMUNITY_TABLE_NAME               = aws_dynamodb_table.user_communities.name
    USER_COMMUNITY_PLATFORM_INDEX           = var.user_community_platform_dynamo_index

    COMMUNITY_USERS_CHANNEL_INDEX           = var.dynamo_community_users_channel_index

    USER_ENGAGEMENT_SCHEDULER_LAMBDA_PREFIX = var.user_engagement_scheduler_lambda_prefix

    USER_TABLE_NAME                         = aws_dynamodb_table.adaptive_users_dynamodb_table.name

    SLACK_LAMBDA_FUNCTION_NAME              = module.slack_user_query_lambda.function_name
    USER_SETUP_LAMBDA_NAME                  = module.user_setup_lambda.function_name


  }
}
