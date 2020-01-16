terraform {
  backend "remote" {
    organization = "ukrops"
    workspaces {
      prefix = "emojibot-"
    }
  }
}

provider aws {
  region = "us-east-1"
}

data "aws_region" "current" {}
data "aws_caller_identity" "current" {}

locals {
  ssm = {
    slack_api_key_path        = "/ukrops/emojiBot/${terraform.workspace}/slackAPIkey"
    legacy_slack_api_key_path = "/ukrops/emojiBot/${terraform.workspace}/LegacySlackAPIkey"
  }
  current = local.workspace[terraform.workspace]
  workspace = {
    production = {
      dns_name = "emojibot.aws.ctrlok.dev"
    }
    development = {
      dns_name        = "emojibot-development.aws.ctrlok.dev"
      best_channel_id = ""
    }
  }
}
