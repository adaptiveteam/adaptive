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
  name               = "${var.client_id}_${var.function_name}"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
}

# Logs policy
resource aws_cloudwatch_log_group "log_group" {
  name              = "/aws/lambda/${aws_lambda_function.lambda[0].function_name}"
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
      "arn:${data.aws_partition.current.partition}:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${var.client_id}_${var.function_name}:*",
    ]
  }
}

resource "aws_iam_policy" "logs" {
  count  = var.enable_cloudwatch_logs ? 1 : 0
  name   = "${var.client_id}_${var.function_name}-logs"
  policy = data.aws_iam_policy_document.logs[0].json
}

resource "aws_iam_policy_attachment" "logs" {
  count      = var.enable_cloudwatch_logs ? 1 : 0
  name       = "${var.client_id}_${var.function_name}-logs"
  roles      = [
    aws_iam_role.lambda.name]
  policy_arn = aws_iam_policy.logs[0].arn
}

# Dead letter config policy
data "aws_iam_policy_document" "dl" {
  count = var.attach_dl_config ? 1 : 0

  statement {
    effect = "Allow"

    actions = [
      "sns:Publish",
      "sqs:SendMessage",
    ]

    resources = [
      lookup(var.dl_config, "target_arn", ""),
    ]
  }
}

resource "aws_iam_policy" "dl" {
  count  = var.attach_dl_config ? 1 : 0
  name   = "${var.client_id}_${var.function_name}-dl"
  policy = data.aws_iam_policy_document.dl[0].json
}

resource "aws_iam_policy_attachment" "dead_letter" {
  count      = var.attach_dl_config ? 1 : 0
  name       = "${var.client_id}_${var.function_name}-dl"
  roles      = [
    aws_iam_role.lambda.name]
  policy_arn = aws_iam_policy.dl[0].arn
}

# VPC config policy
data "aws_iam_policy_document" "network" {
  statement {
    effect = "Allow"

    actions = [
      "ec2:CreateNetworkInterface",
      "ec2:DescribeNetworkInterfaces",
      "ec2:DeleteNetworkInterface",
    ]

    resources = [
      "*",
    ]
  }
}

resource "aws_iam_policy" "network" {
  count = var.attach_vpc_config ? 1 : 0

  name   = "${var.client_id}_${var.function_name}-network"
  policy = data.aws_iam_policy_document.network.json
}

resource "aws_iam_policy_attachment" "network" {
  count = var.attach_vpc_config ? 1 : 0

  name       = "${var.client_id}_${var.function_name}-network"
  roles      = [
    aws_iam_role.lambda.name]
  policy_arn = aws_iam_policy.network[0].arn
}

# Attach additional policy if required
resource "aws_iam_policy" "additional" {
  count  = var.attach_policy ? 1 : 0
  name   = "${var.client_id}_${var.function_name}"
  policy = var.policy
}

resource "aws_iam_policy_attachment" "additional" {
  count = var.attach_policy ? 1 : 0

  name       = "${var.client_id}_${var.function_name}"
  roles      = [
    aws_iam_role.lambda.name]
  policy_arn = aws_iam_policy.additional[0].arn
}

## AWS X-ray policy
data "aws_iam_policy_document" "xray" {
  statement {
    effect = "Allow"

    actions = [
      "xray:PutTraceSegments",
      "xray:PutTelemetryRecords"
    ]

    resources = [
      "*",
    ]
  }
}

resource "aws_iam_policy" "xray" {
  name   = "${var.client_id}_${var.function_name}-xray"
  policy = data.aws_iam_policy_document.xray.json
}

resource "aws_iam_policy_attachment" "xray" {
  name       = "${var.client_id}_${var.function_name}-xray"
  roles      = [
    aws_iam_role.lambda.name]
  policy_arn = aws_iam_policy.xray.arn
}
