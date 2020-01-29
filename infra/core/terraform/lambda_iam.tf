data "aws_iam_policy_document" "assume_lambda_role" {
  statement {
    principals {
      type        = "Service"
      identifiers = [
        "lambda.amazonaws.com",
        "events.amazonaws.com"]
    }
    actions = ["sts:AssumeRole"]
  }
}

# Lambda policy
resource "aws_iam_role" "lambda" {
  name               = "lambda.role.${var.client_id}"
  assume_role_policy = data.aws_iam_policy_document.assume_lambda_role.json
}
