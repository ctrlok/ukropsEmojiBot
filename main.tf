terraform {
  backend "remote" {
    organization = "ukrops"
    workspaces {
      name = "bestbot"
    }
  }
}

provider aws {
  region = "us-east-1"
}

