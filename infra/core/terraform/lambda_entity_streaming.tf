locals {
  entity_streaming_function_name_suffix = "entity-streaming-lambda-go"
  entity_streaming_function_name = "${var.client_id}_${local.entity_streaming_function_name_suffix}"
}
module "entity_streaming" {
  source = "../../../terraform-modules/adaptive-lambda-s3"
  s3_bucket = aws_s3_bucket.binary_bucket.bucket
  s3_key = aws_s3_bucket_object.adaptive_zip.key
  source_hash = data.archive_file.adaptive_lambda_zip.output_md5

  client_id     = var.client_id
  handler       = "adaptive"
  function_name_suffix = local.entity_streaming_function_name_suffix
  runtime       = var.lambda_runtime
  timeout       = var.lambda_timeout

  reserved_concurrent_executions = -1

  // Add environment variables.
  environment_variables = merge(local.environment_variables, {
    LAMBDA_ROLE   = "entity-streaming"
    LOG_NAMESPACE = "entity-streaming"
    DB_USER = module.reporting_db.this_db_instance_username
    DB_PASS = module.reporting_db.this_db_instance_password
    DB_NAME = module.reporting_db.this_db_instance_name
    DB_HOST = module.reporting_db.this_db_instance_endpoint
    STREAM_EVENT_MAPPER_LAMBDA = module.stream_event_mapping.function_name
  })

  // Attach extra policy
  attach_policy = true
  policy        = data.aws_iam_policy_document.entity_streaming_policy.json

  tags = local.default_tags
}

data "aws_iam_policy_document" "entity_streaming_policy" {
  statement {
    actions = [
      "dynamodb:GetRecords",
      "dynamodb:GetShardIterator",
      "dynamodb:DescribeStream",
      "dynamodb:ListStreams",
    ]
    resources = [
      aws_dynamodb_table.client_config_dynamodb_table.stream_arn,
      aws_dynamodb_table.adaptive_users_dynamodb_table.stream_arn,
      aws_dynamodb_table.adaptive_value_dynamodb_table.stream_arn,
      aws_dynamodb_table.vision.stream_arn,
      aws_dynamodb_table.adaptive_user_feedback_dynamodb_table.stream_arn,
      aws_dynamodb_table.ad_hoc_holidays.stream_arn,
      aws_dynamodb_table.adaptive_user_engagements_dynamo_table.stream_arn,
      aws_dynamodb_table.user_objective_dynamodb_table.stream_arn,
      aws_dynamodb_table.user_objectives_progress.stream_arn,
      aws_dynamodb_table.user_communities.stream_arn,
      aws_dynamodb_table.community_users.stream_arn,
      aws_dynamodb_table.initiative_communities.stream_arn,
      aws_dynamodb_table.strategy_objectives.stream_arn,
      aws_dynamodb_table.strategy_communities.stream_arn,
      aws_dynamodb_table.accountability_partnership_rejections_table.stream_arn,
    ]
  }

  statement {
    actions = ["lambda:InvokeFunction"]
    resources = [module.stream_event_mapping.function_arn]
  }
}

resource "aws_lambda_event_source_mapping" "client_config_streaming_lambda_source_mapping" {
  event_source_arn = aws_dynamodb_table.client_config_dynamodb_table.stream_arn
  function_name = module.entity_streaming.function_arn
  starting_position = "LATEST"
}

resource "aws_lambda_event_source_mapping" "user_streaming_lambda_source_mapping" {
  event_source_arn = aws_dynamodb_table.adaptive_users_dynamodb_table.stream_arn
  function_name = module.entity_streaming.function_arn
  starting_position = "LATEST"
}

resource "aws_lambda_event_source_mapping" "competency_streaming_lambda_source_mapping" {
  event_source_arn = aws_dynamodb_table.adaptive_value_dynamodb_table.stream_arn
  function_name = module.entity_streaming.function_arn
  starting_position = "LATEST"
}

resource "aws_lambda_event_source_mapping" "vision_streaming_lambda_source_mapping" {
  event_source_arn = aws_dynamodb_table.vision.stream_arn
  function_name = module.entity_streaming.function_arn
  starting_position = "LATEST"
}

resource "aws_lambda_event_source_mapping" "user_feedback_streaming_lambda_source_mapping" {
  event_source_arn = aws_dynamodb_table.adaptive_user_feedback_dynamodb_table.stream_arn
  function_name = module.entity_streaming.function_arn
  starting_position = "LATEST"
}

resource "aws_lambda_event_source_mapping" "holidays_streaming_lambda_source_mapping" {
  event_source_arn = aws_dynamodb_table.ad_hoc_holidays.stream_arn
  function_name = module.entity_streaming.function_arn
  starting_position = "LATEST"
}

resource "aws_lambda_event_source_mapping" "engagements_streaming_lambda_source_mapping" {
  event_source_arn = aws_dynamodb_table.adaptive_user_engagements_dynamo_table.stream_arn
  function_name = module.entity_streaming.function_arn
  starting_position = "LATEST"
}

resource "aws_lambda_event_source_mapping" "user_objective_streaming_lambda_source_mapping" {
  event_source_arn = aws_dynamodb_table.user_objective_dynamodb_table.stream_arn
  function_name = module.entity_streaming.function_arn
  starting_position = "LATEST"
}

resource "aws_lambda_event_source_mapping" "user_objective_progress_streaming_lambda_source_mapping" {
  event_source_arn = aws_dynamodb_table.user_objectives_progress.stream_arn
  function_name = module.entity_streaming.function_arn
  starting_position = "LATEST"
}

resource "aws_lambda_event_source_mapping" "communities_streaming_lambda_source_mapping" {
  event_source_arn = aws_dynamodb_table.user_communities.stream_arn
  function_name = module.entity_streaming.function_arn
  starting_position = "LATEST"
}

resource "aws_lambda_event_source_mapping" "community_users_streaming_lambda_source_mapping" {
  event_source_arn = aws_dynamodb_table.community_users.stream_arn
  function_name = module.entity_streaming.function_arn
  starting_position = "LATEST"
}

resource "aws_lambda_event_source_mapping" "initiative_communities_streaming_lambda_source_mapping" {
  event_source_arn = aws_dynamodb_table.initiative_communities.stream_arn
  function_name = module.entity_streaming.function_arn
  starting_position = "LATEST"
}

resource "aws_lambda_event_source_mapping" "objectives_streaming_lambda_source_mapping" {
  event_source_arn = aws_dynamodb_table.strategy_objectives.stream_arn
  function_name = module.entity_streaming.function_arn
  starting_position = "LATEST"
}

resource "aws_lambda_event_source_mapping" "objective_communities_streaming_lambda_source_mapping" {
  event_source_arn = aws_dynamodb_table.strategy_communities.stream_arn
  function_name = module.entity_streaming.function_arn
  starting_position = "LATEST"
}

resource "aws_lambda_event_source_mapping" "partnership_rejections_streaming_lambda_source_mapping" {
  event_source_arn = aws_dynamodb_table.accountability_partnership_rejections_table.stream_arn
  function_name = module.entity_streaming.function_arn
  starting_position = "LATEST"
}
