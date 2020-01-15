resource "aws_api_gateway_rest_api" "emojiBot" {
  name = "emojiBot" //TODO: change to a workspace-related name in develop first
  endpoint_configuration {
    types = ["REGIONAL"]
  }
}


// custom path for recieving emojies from slack API
resource "aws_api_gateway_resource" "slackConnector" {
  path_part   = "slack"
  parent_id   = aws_api_gateway_rest_api.emojiBot.root_resource_id
  rest_api_id = aws_api_gateway_rest_api.emojiBot.id
}

resource "aws_api_gateway_method" "slackConnector" {
  rest_api_id   = aws_api_gateway_rest_api.emojiBot.id
  resource_id   = aws_api_gateway_resource.slackConnector.id
  http_method   = "POST"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "slackConnector" {
  rest_api_id             = aws_api_gateway_rest_api.emojiBot.id
  resource_id             = aws_api_gateway_resource.slackConnector.id
  http_method             = aws_api_gateway_method.slackConnector.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.slackConnector.invoke_arn
}

# Deployment
resource "aws_api_gateway_deployment" "botv1" {
  rest_api_id = aws_api_gateway_rest_api.emojiBot.id
  stage_name  = "prod"

  depends_on = [
    aws_api_gateway_integration.slackConnector,
  ]
}

# Create access to url https://emojibot.aws.ctrlok.dev/v1/slack
resource "aws_api_gateway_base_path_mapping" "botv1" {
  api_id      = aws_api_gateway_rest_api.emojiBot.id
  domain_name = aws_api_gateway_domain_name.emojibot.domain_name
  stage_name  = aws_api_gateway_deployment.botv1.stage_name
  base_path   = "v1"
}