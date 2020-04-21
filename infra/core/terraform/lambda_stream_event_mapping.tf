locals {
  stream_event_mapping_function_name_suffix = "stream-event-mapping-lambda-go"
  stream_event_mapping_function_name = "${var.client_id}_${local.stream_event_mapping_function_name_suffix}"
}
module "stream_event_mapping" {
  source = "../../../terraform-modules/adaptive-lambda-s3"
  s3_bucket = aws_s3_bucket.binary_bucket.bucket
  s3_key = aws_s3_bucket_object.adaptive_zip.key
  source_hash = data.archive_file.adaptive_lambda_zip.output_md5

  client_id     = var.client_id
  handler       = "adaptive"
  function_name_suffix = local.stream_event_mapping_function_name_suffix
  runtime       = var.lambda_runtime
  timeout       = var.lambda_timeout

  reserved_concurrent_executions = -1

  // Add environment variables.
  environment_variables = merge(local.environment_variables, {
    LAMBDA_ROLE   = "stream-event-mapping"
    LOG_NAMESPACE = "stream-event-mapping"
    DB_USER = module.reporting_db.this_db_instance_username
    DB_PASS = module.reporting_db.this_db_instance_password
    DB_NAME = module.reporting_db.this_db_instance_name
    DB_HOST = module.reporting_db.this_db_instance_endpoint
    # STREAM_EVENT_MAPPER_LAMBDA = module.stream_event_mapping.function_name
  })

  tags = local.default_tags
}
