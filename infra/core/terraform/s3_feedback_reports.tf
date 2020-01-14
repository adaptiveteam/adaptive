resource "aws_s3_bucket" "adaptive-feedback-reports-bucket" {
  bucket = "${local.client_id}-feedback-reports"
  acl    = "private"
  region = local.region

  versioning {
    enabled = true
  }

  tags = local.default_tags
}
