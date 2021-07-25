variable "profile" {
  type = string
  default = "default"
}

variable "region" {
  type = string
  default = "eu-west-2"
}

variable "domain" {
  type = string
}

variable "cert_domain" {
  type = string
}

variable "zone" {
  type = string
}

variable "notification_sns_queue_name" {
  type = string
  description = "The name of the SNS queue to send ok/error alarms if the lambda stops working."
}

provider "aws" {
  region = var.region
  profile = var.profile
}

terraform {
  backend "s3" {
    bucket = "ana-terraform-state-prod"
    key = "ddns-lambda/terraform.tfstate"
    region = "eu-west-2"
  }
}