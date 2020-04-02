data "archive_file" "adaptive-lambda-zip" {
  type = "zip"
  # source_file = "../../../bin/adaptive"
  source_dir = "../../../target"
  output_path = "lambdas/adaptive.zip"
}
