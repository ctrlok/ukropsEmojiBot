resource "aws_cloudwatch_log_group" "slackConnector" {
  name              = "/aws/lambda/ukrops/emojiBot-${terraform.workspace}"
  retention_in_days = 14
}

resource "aws_lambda_function" "slackConnector" {
  filename         = data.archive_file.lambda.output_path
  function_name    = "ukrops_slackConnector-${terraform.workspace}"
  role             = aws_iam_role.slackConnector.arn
  handler          = "slackConnectorBin"
  source_code_hash = data.archive_file.lambda.output_base64sha256
  timeout          = 600
  runtime          = "go1.x"
  environment {
    variables = {
      EMOJIBOT_BEST_CHANNEL_ID               = local.current.best_channel_id
      EMOJIBOT_SSM_SLACK_API_KEY_PATH        = local.ssm.slack_api_key_path
      EMOJIBOT_SSM_SLACK_API_LEGACY_KEY_PATH = local.ssm.legacy_slack_api_key_path
      EMOJIBOT_AWS_REGION                    = data.aws_region.current.name
    }
  }
}
