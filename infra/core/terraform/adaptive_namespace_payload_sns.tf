resource "aws_sns_topic" "namespace_payload" {
  name         = "${var.client_id}_namespace_payload"
  display_name = "${var.client_id}_namespace_payload"
}

