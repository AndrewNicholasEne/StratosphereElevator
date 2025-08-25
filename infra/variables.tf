provider "aws" {
  region = "eu-west-1"
  default_tags {
    tags = {
      Project   = "StratosphereElevator"
      Env       = "dev"
      ManagedBy = "Terraform"
    }
  }
}

variable "aws_region" {
  type    = string
  default = "eu-west-1"
}

variable "image_tag" {
  type    = string
  default = "dev"
}
