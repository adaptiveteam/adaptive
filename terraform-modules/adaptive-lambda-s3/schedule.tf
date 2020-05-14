resource "aws_cloudwatch_event_rule" "rule" {
  count               = var.schedule ? 1 : 0
  name                = "${var.client_id}_${var.schedule_name}"
  description         = var.schedule_description
  schedule_expression = var.schedule_expression
  is_enabled          = var.schedule_is_enabled
}

resource "aws_cloudwatch_event_target" "target" {
  count      = var.schedule ? 1 : 0
  arn        = concat(aws_lambda_function.lambda.*.arn, ["unavailable1"])[0]
  rule       = concat(aws_cloudwatch_event_rule.rule.*.name, ["rule is absent"])[0]
  input      = var.schedule_invoke_json
  target_id  = concat(aws_lambda_function.lambda.*.function_name, ["unavailable2"])[0]
  depends_on = [aws_lambda_function.lambda]
}

resource "aws_lambda_permission" "cloudwatch_invoke_lambda" {
  count         = var.schedule ? 1 : 0
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = concat(aws_lambda_function.lambda.*.arn, ["unavailable3"])[0]
  principal     = "events.amazonaws.com"
  source_arn    = concat(aws_cloudwatch_event_rule.rule.*.arn,["rule unavailable3"])[0]
}
