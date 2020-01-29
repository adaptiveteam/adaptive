resource "aws_cloudwatch_event_rule" "rule" {
  name                = "${var.client_id}_${var.schedule_name}"
  description         = var.schedule_description
  schedule_expression = var.schedule_expression
  is_enabled          = var.schedule_is_enabled
}

resource "aws_cloudwatch_event_target" "target" {
  arn        = var.target_lambda_arn
  rule       = aws_cloudwatch_event_rule.rule.name
  input      = var.schedule_invoke_json
  target_id  = "${var.client_id}_${var.schedule_name}"
}

resource "aws_lambda_permission" "cloudwatch_invoke_lambda" {
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = var.target_lambda_arn
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.rule.arn
}
