# Logs policy
resource aws_cloudwatch_log_group "log_group" {
  name              = "/aws/lambda/${var.function_name}"
  retention_in_days = 14
}

data "aws_iam_policy_document" "lambda_log" {
  statement {
    actions   = ["logs:CreateLogStream"]
    resources = ["arn:${data.aws_partition.current.partition}:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${var.function_name}"]//"*"]
  }
  statement {
    actions   = ["logs:PutLogEvents"]
    resources = ["arn:${data.aws_partition.current.partition}:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${var.function_name}:log-stream:*"]
  }
  statement {
    actions   = ["logs:CreateLogGroup"]
    resources = ["*"]
  }
}

resource "aws_iam_policy" "lambda_log" {
  name   = "${var.function_name}_log"
  policy = data.aws_iam_policy_document.lambda_log.json
}

resource "aws_iam_policy_attachment" "logs" {
  name       = "${var.function_name}_log"
  roles      = [var.role_name]
  policy_arn = aws_iam_policy.lambda_log.arn
}

module "error_alarm" {
  // TODO: Pin to released version once this repo is released
  source = "github.com/dwp/terraform-aws-metric-filter-alarm?ref=master"
  log_group_name = aws_cloudwatch_log_group.log_group.name// module.slack_message_processor_lambda.log_group_name
  metric_namespace = "${var.client_id}-AWS/Lambda"
  pattern = "ERROR"
  alarm_name = "${var.function_name}_errors"
  alarm_action_arns = [var.errors_sns_topic_arn]
  period = "60"
  threshold = "1"
  statistic = "SampleCount"
}


## AWS X-ray policy
data "aws_iam_policy_document" "xray" {
  statement {
    actions = [
      "xray:PutTraceSegments",
      "xray:PutTelemetryRecords"
    ]
    resources = ["*"]
  }
}

resource "aws_iam_policy" "xray" {
  name   = "${var.function_name}-xray"
  policy = data.aws_iam_policy_document.xray.json
}

resource "aws_iam_policy_attachment" "xray" {
  name       = "${var.function_name}-xray"
  roles      = [var.role_name]
  policy_arn = aws_iam_policy.xray.arn
}
