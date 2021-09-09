resource "aws_instance" "web" {
  ami           = "ami-00a2249139ac35088" 
  instance_type = "t2.micro"

}
