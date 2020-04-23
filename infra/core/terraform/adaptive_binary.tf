data "archive_file" "adaptive_lambda_zip" {
  type = "zip"
  # source_file = "../../../bin/adaptive"
  source_dir = "../../../target"
  output_path = "lambdas/adaptive.zip"
}

resource "aws_s3_bucket" "binary_bucket" {
  bucket = "${local.client_id}-adaptive-binary"
  acl    = "private"
  region = local.region

  versioning {
    enabled = true
  }

  tags = local.default_tags

}

resource "aws_s3_bucket_object" "adaptive_zip" {
  bucket = aws_s3_bucket.binary_bucket.bucket
  key    = "adaptive.zip"
  source = data.archive_file.adaptive_lambda_zip.output_path

  # The filemd5() function is available in Terraform 0.11.12 and later
  # For Terraform 0.11.11 and earlier, use the md5() function and the file() function:
  # etag = "${md5(file("path/to/file"))}"
  etag = data.archive_file.adaptive_lambda_zip.output_md5//base64sha256
}
