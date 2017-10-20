terraform {
  backend "s3" {
    bucket  = "curalate-configuration"
    key     = "terraform/testdeleteifseen/us-east-1/qa/terraform.tfstate"
    region  = "us-east-1"
    profile = "qa"
    encrypt = "true"
  }
}
