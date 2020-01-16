resource "aws_cloudwatch_event_rule" "rule" {
  count               = var.schedule ? 1 : 0
  name                = "${var.client_id}_${var.schedule_name}"
  description         = var.schedule_description
  schedule_expression = var.schedule_expression
  is_enabled          = var.schedule_is_enabled
}

resource "aws_cloudwatch_event_target" "target" {
  count      = var.schedule ? 1 : 0
  arn        = aws_lambda_function.lambda.0.arn
  rule       = aws_cloudwatch_event_rule.rule[0].name
  input      = var.schedule_invoke_json
  target_id  = aws_lambda_function.lambda.0.function_name
  depends_on = [aws_lambda_function.lambda]
}

resource "aws_lambda_permission" "cloudwatch_invoke_lambda" {
  count         = var.schedule ? 1 : 0
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.lambda.0.arn
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.rule[0].arn
}
