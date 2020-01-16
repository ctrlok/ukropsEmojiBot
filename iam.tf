resource "aws_iam_role" "slackConnector" {
  name = "ukrops_emojiBot_slackConnector-${terraform.workspace}"

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
resource "aws_iam_role_policy" "slackConnectorLogs" {
  name = "slackConnector_cloudwatch_access-${terraform.workspace}"
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

resource "aws_iam_role_policy" "slackConnectorSecrets" {
  name = "slackConnector_secrets_access-${terraform.workspace}"
  role = aws_iam_role.slackConnector.id

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
        "Effect": "Allow",
        "Action": [
            "ssm:GetParameter*"
        ],
        "Resource": "arn:aws:ssm:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:parameter/ukrops/emojiBot/${terraform.workspace}/*"
    },
    {
        "Effect": "Allow",
        "Action": [
            "kms:Decrypt"
        ],
        "Resource": "${aws_kms_key.emojiBot_slackApi.arn}"
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
