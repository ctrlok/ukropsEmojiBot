resource "aws_iam_role" "slackConnector" {
  name = "ukrops_emojiBot_slackConnector"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow"
    }
  ]
}
EOF

  tags = {
    Terraform = "true"
  }
}

# Allow lambda to write logs
resource "aws_iam_role_policy" "slackConnector" {
  name = "slackConnector_cloudwatch_access"
  role = aws_iam_role.slackConnector.id

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Effect": "Allow",
      "Resource": "*"
    }
  ]
}
EOF
}

resource "aws_lambda_permission" "slackConnector_allow_APIGateway" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.slackConnector.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.emojiBot.execution_arn}/*"
}
