data "aws_iam_policy_document" "assume_role" {
  statement {
    effect  = "Allow"
    actions = [
      "sts:AssumeRole"]

    principals {
      identifiers = [
        "lambda.amazonaws.com"]
      type        = "Service"
    }
  }

  statement {
    effect  = "Allow"
    actions = [
      "sts:AssumeRole"
    ]

    principals {
      identifiers = [
        "events.amazonaws.com"]
      type        = "Service"
    }
  }
}

# Lambda policy
resource "aws_iam_role" "lambda" {
  name               = local.function_name
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
}

# Logs policy
resource aws_cloudwatch_log_group "log_group" {
  name              = "/aws/lambda/${element(concat(aws_lambda_function.lambda.*.function_name, ["unavailable"]), 0)}"
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
      "arn:${data.aws_partition.current.partition}:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${local.function_name}:*",
    ]
  }
}

resource "aws_iam_policy" "logs" {
  count  = var.enable_cloudwatch_logs ? 1 : 0
  name   = "${local.function_name}-logs"
  policy = data.aws_iam_policy_document.logs[0].json
}

resource "aws_iam_policy_attachment" "logs" {
  count      = var.enable_cloudwatch_logs ? 1 : 0
  name       = "${local.function_name}-logs"
  roles      = [aws_iam_role.lambda.name]
  policy_arn = concat(aws_iam_policy.logs.*.arn,["absent logs"])[0]
}

# Attach additional policy if required
resource "aws_iam_policy" "additional" {
  count  = var.attach_policy ? 1 : 0
  name   = local.function_name
  policy = var.policy
}

resource "aws_iam_policy_attachment" "additional" {
  count = var.attach_policy ? 1 : 0
  name       = local.function_name
  roles      = [aws_iam_role.lambda.name]
  policy_arn = concat(aws_iam_policy.additional.*.arn,["policy not found"])[0]
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
  name   = "${local.function_name}-xray"
  policy = data.aws_iam_policy_document.xray.json
}

resource "aws_iam_policy_attachment" "xray" {
  name       = "${local.function_name}-xray"
  roles      = [aws_iam_role.lambda.name]
  policy_arn = aws_iam_policy.xray.arn
}
