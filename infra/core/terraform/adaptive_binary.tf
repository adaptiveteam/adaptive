data "archive_file" "adaptive-lambda-zip" {
  type = "zip"
  source_file = "../../../bin/adaptive"
  output_path = "lambdas/adaptive.zip"
}
