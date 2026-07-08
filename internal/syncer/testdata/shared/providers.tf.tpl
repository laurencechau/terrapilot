terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "${provider_version}"
    }
  }
}

provider "aws" {
  region = "${region}"
}
