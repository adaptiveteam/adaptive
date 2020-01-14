resource "aws_sns_topic" "platform_notification" {
  name            = "${var.client_id}_platform_notification"
  display_name    = "${var.client_id}_platform_notification"
  delivery_policy = data.template_file.delivery_policy.rendered
}

data "template_file" "delivery_policy" {
  template = file("templates/sns_delivery_policy.tpl.json")
}

