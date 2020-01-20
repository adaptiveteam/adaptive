# Logs policy
resource aws_cloudwatch_log_group "log_group" {
  name              = "/aws/lambda/${element(aws_lambda_function.lambda.*.function_name, 0)}"
  retention_in_days = 14
}

data "aws_iam_policy_document" "logs" {
  count = var.enable_cloudwatch_logs ? 1 : 0
  statement {
    effect    = "Allow"
    actions   = [
      "logs:CreateLogStream",
      "logs:PutLogEvents",]
    resources = [
      "arn:${data.aws_partition.current.partition}:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${local.prefix_name}_lambda:*",
    ]
  }
}

resource "aws_iam_policy" "logs" {
  count  = var.enable_cloudwatch_logs ? 1 : 0
  name   = "${local.prefix_name}_logs"
  policy = data.aws_iam_policy_document.logs[0].json
}

resource "aws_iam_policy_attachment" "logs" {
  count      = var.enable_cloudwatch_logs ? 1 : 0
  name       = "${local.prefix_name}_logs"
  roles      = [
    aws_iam_role.lambda.name]
  policy_arn = aws_iam_policy.logs[0].arn
}
