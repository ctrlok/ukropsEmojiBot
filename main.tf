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
