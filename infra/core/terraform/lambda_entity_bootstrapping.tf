locals {
  entity_bootstrapping_function_name_suffix = "entity-bootstrapping-lambda-go"
  entity_bootstrapping_function_name = "${var.client_id}_${local.entity_bootstrapping_function_name_suffix}"
}
module "entity_bootstrapping" {
  source = "../../../terraform-modules/adaptive-lambda-s3"
  s3_bucket = aws_s3_bucket.binary_bucket.bucket
  s3_key = aws_s3_bucket_object.adaptive_zip.key
  source_hash = data.archive_file.adaptive_lambda_zip.output_md5

  client_id     = var.client_id
  handler       = "adaptive"
  function_name_suffix = local.entity_bootstrapping_function_name_suffix
  runtime       = var.lambda_runtime
  timeout       = 300 // 5 minutes

  reserved_concurrent_executions = -1

  // Add environment variables.
  environment_variables = merge(local.environment_variables, {
    LAMBDA_ROLE   = "entity-bootstrapping"
    LOG_NAMESPACE = "entity-bootstrapping"
    STREAM_EVENT_MAPPER_LAMBDA = module.stream_event_mapping.function_name
  })

  // Attach extra policy
  attach_policy = true
  policy        = data.aws_iam_policy_document.entity_bootstrapping_policy.json

  tags = local.default_tags
}

data "aws_iam_policy_document" "entity_bootstrapping_policy" {
  statement {
    actions = [
      "dynamodb:ListTables",
      "dynamodb:Scan",
    ]
    resources = ["*"]
  }

  statement {
    actions = ["lambda:InvokeFunction"]
    resources = [module.stream_event_mapping.function_arn]
  }
}
