# Logs policy
resource aws_cloudwatch_log_group "consolidated_log_group" {
  name              = "/aws/lambda/${var.client_id}_consolidated"
  retention_in_days = 14
}

data "aws_iam_policy_document" "lambda_log" {
  statement {
    actions   = [
      "logs:CreateLogStream",
      "logs:PutLogEvents",]
    resources = [
      "arn:${data.aws_partition.current.partition}:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/*:*",
    ]
  }
}

resource "aws_iam_policy" "lambda_log" {
  name   = "${var.client_id}_lambda_log"
  policy = data.aws_iam_policy_document.lambda_log.json
}

resource "aws_iam_policy_attachment" "logs" {
  name       = "${var.client_id}_lambda_log"
  roles      = [aws_iam_role.lambda.name]
  policy_arn = aws_iam_policy.lambda_log.arn
}
