resource "aws_kms_key" "emojiBot_slackApi" {
  description = "Slack API keys for EmojiBots"
}

output "kms" {
  value = aws_kms_key.emojiBot_slackApi
}