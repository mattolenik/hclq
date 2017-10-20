resource "aws_lambda_function" "function" {
  filename         = "deployment.zip"
  function_name    = "lambdafunction"
  description      = "Performs some action"
  role             = "somerole"
  handler          = "module.lambda_handler"
  runtime          = "python3.6"
  source_code_hash = "somehash"
  timeout          = "30"

  environment = {
    variables = {
      ENV = "qa"
      APP = "something"
    }
  }
}

variable "amis" {
  type = "map"

  default = {
    "us-east-1" = "ami-abd1234"
    "us-west-2" = "ami-abd1234"
  }
}

locals {
  testlist        = ["one", "two", "three"]
  testliteral     = "somethinghere"
  testlistnested  = [[1], [2], [3]]
  testlistnested2 = [[1], [2], 3, "4"]

  testmap = {
    a = 1
    b = 2
    c = [1, 2, 3]

    d = {
      anothermap = "here"
    }
  }

  testmap2 = "${var.amis}"

  testobj {
    testvar = 1
  }
}
