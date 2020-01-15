resource "aws_s3_bucket" "backup_bucket" {
  bucket = "${local.client_id}-adaptive-backup"
  acl    = "private"
  region = local.region

  versioning {
    enabled = true
  }

  tags = local.default_tags

}
