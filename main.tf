terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.24"
    }
  }

  required_version = ">= 1.2.0"

  cloud {
    organization = "flat35hd99"

    workspaces {
      name = "todo-go"
    }
  }
}

provider "aws" {
  region = "us-west-2"
}
