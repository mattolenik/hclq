data "foo" {
  bin = [1, 2, 3]

  bar = "foo string"
}

data "baz" {
  bin = [4, 5, 6]

  bar = "baz string"
}

resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"
  id = "main-vpc"
}

resource "aws_vpc" "dev" {
  cidr_block = "10.0.0.0/16"
  id = "dev-vpc"
}

data "aws_vpc" "precreated1" {
  id = "some-vpc-id-1"
}

data "aws_vpc" "other2" {
  id = "some-vpc-id-2"
}
