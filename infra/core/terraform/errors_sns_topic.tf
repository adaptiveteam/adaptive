resource "aws_sns_topic" "errors" {
  name         = "${var.client_id}_errors"
  display_name = "${var.client_id}_errors"
}

