provider "aws" { region = "eu-west-1" }

resource "aws_s3_bucket" "tf_state" {
  bucket = "stratosphereelevator-tfstate-andrew"

  lifecycle { prevent_destroy = true }

  tags = {
    Name      = "tfstate"
    Project   = "StratosphereElevator"
    ManagedBy = "Terraform"
  }
}

resource "aws_s3_bucket_ownership_controls" "tf_state" {
  bucket = aws_s3_bucket.tf_state.id
  rule { object_ownership = "BucketOwnerEnforced" }
}

resource "aws_s3_bucket_public_access_block" "tf_state" {
  bucket                  = aws_s3_bucket.tf_state.id
  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_s3_bucket_versioning" "tf_state" {
  bucket = aws_s3_bucket.tf_state.id
  versioning_configuration { status = "Enabled" }
}

resource "aws_s3_bucket_server_side_encryption_configuration" "tf_state" {
  bucket = aws_s3_bucket.tf_state.id
  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256" # or "aws:kms" + kms_master_key_id for CMK
    }
  }
}

resource "aws_dynamodb_table" "tf_lock" {
  name         = "stratosphereelevator-tf-lock"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "LockID"

  attribute {
    name = "LockID"
    type = "S"
  }

  tags = {
    Name      = "tfstate-lock"
    Project   = "StratosphereElevator"
    ManagedBy = "Terraform"
  }
}
