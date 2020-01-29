data "archive_file" "adaptive-lambda-zip" {
  type = "zip"
  source_file = "../../../bin/adaptive"
  output_path = "lambdas/adaptive.zip"
}

resource "aws_s3_bucket" "adaptive-binary-bucket" {
  bucket = "${local.client_id}-binary"
  acl    = "private"
  region = local.region

  versioning {
    enabled = true
  }

  tags = local.default_tags
}

resource "aws_s3_bucket_object" "adaptive_binary" {
  bucket = aws_s3_bucket.adaptive-binary-bucket.bucket
  key    = "adaptive.zip"
  source = data.archive_file.adaptive-lambda-zip.output_path

  # The filemd5() function is available in Terraform 0.11.12 and later
  # For Terraform 0.11.11 and earlier, use the md5() function and the file() function:
  # etag = "${md5(file("path/to/file"))}"
  etag = filemd5(data.archive_file.adaptive-lambda-zip.output_path)
}
