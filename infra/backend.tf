terraform {
  backend "s3" {
    bucket         = "stratosphereelevator-tfstate-andrew"
    key            = "envs/dev/terraform.tfstate"
    region         = "eu-west-1"
    dynamodb_table = "stratosphereelevator-tf-lock"
    encrypt        = true
  }
}
