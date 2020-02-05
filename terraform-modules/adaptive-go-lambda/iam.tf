data "aws_iam_policy_document" "assume_role" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      identifiers = ["lambda.amazonaws.com"]
      type        = "Service"
    }
  }

  statement {
    actions = ["sts:AssumeRole"]
    principals {
      identifiers = ["events.amazonaws.com"]
      type        = "Service"
    }
  }
}

# Lambda policy
resource "aws_iam_role" "lambda" {
  name               = "role.${var.id}.${var.name}"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
}
