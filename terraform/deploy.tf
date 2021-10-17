 provider "aws" {
   region = "us-east-1"
}

resource "aws_dynamodb_table" "ddbtable" {
  name             = "Sensors"
  hash_key         = "Key"
  billing_mode   = "PROVISIONED"
  read_capacity  = 5
  write_capacity = 5
  attribute {
    name = "Key"
    type = "S"
  }
}

resource "aws_iam_role_policy" "lambda_policy" {
  name = "lambda_policy"
  role = aws_iam_role.role_for_LDC.id

  policy = file("privilege-policy.json")
}


resource "aws_iam_role" "role_for_LDC" {
  name = "myrole"

  assume_role_policy = file("trust-policy.json")

}

resource "aws_lambda_function" "put" {

  function_name = "put"
  s3_bucket     = "mybucket-sdcc-lambdas"
  s3_key        = "put.zip"
  role          = aws_iam_role.role_for_LDC.arn
  handler       = "put"
  runtime       = "go1.x"
  memory_size   = "512"
}

 resource "aws_lambda_function" "get" {

   function_name = "get"
   s3_bucket     = "mybucket-sdcc-lambdas"
   s3_key        = "get.zip"
   role          = aws_iam_role.role_for_LDC.arn
   handler       = "get"
   runtime       = "go1.x"
   memory_size   = "512"
 }

 resource "aws_lambda_function" "delete" {

   function_name = "delete"
   s3_bucket     = "mybucket-sdcc-lambdas"
   s3_key        = "delete.zip"
   role          = aws_iam_role.role_for_LDC.arn
   handler       = "delete"
   runtime       = "go1.x"
   memory_size   = "512"
 }

 resource "aws_lambda_function" "append" {

   function_name = "append"
   s3_bucket     = "mybucket-sdcc-lambdas"
   s3_key        = "append.zip"
   role          = aws_iam_role.role_for_LDC.arn
   handler       = "append"
   runtime       = "go1.x"
   memory_size   = "512"
 }
