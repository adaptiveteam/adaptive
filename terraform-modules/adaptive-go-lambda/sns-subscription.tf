resource "aws_sns_topic_subscription" "lambda_sns" {
  topic_arn = var.input_sns_topic_arn
  protocol  = "lambda"
  endpoint  = local.function_arn
}

resource "aws_lambda_permission" "lambda_sns" {
  statement_id  = "AllowExecutionFromSNS"
  action        = "lambda:InvokeFunction"
  function_name = local.function_arn
  principal     = "sns.amazonaws.com"
  source_arn    = var.input_sns_topic_arn
}
