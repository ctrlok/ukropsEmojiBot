resource "aws_cloudwatch_log_group" "slackConnector" {
  name              = "/aws/lambda/ukrops/emojiBot"
  retention_in_days = 14
}

resource "aws_lambda_function" "slackConnector" {
  filename         = data.archive_file.lambda.output_path
  function_name    = "ukrops_slackConnector"
  role             = aws_iam_role.slackConnector.arn
  handler          = "slackConnectorBin"
  source_code_hash = data.archive_file.lambda.output_base64sha256
  timeout          = 600
  runtime          = "go1.x"
}
