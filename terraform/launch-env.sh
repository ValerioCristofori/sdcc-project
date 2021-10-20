#!/bin/bash

cd ../lambda
# build lambda func and create the zip
env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /tmp/put ./put
zip -j ./put.zip /tmp/put

env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /tmp/get ./get
zip -j ./get.zip /tmp/get

env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /tmp/delete ./delete
zip -j ./delete.zip /tmp/delete

env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /tmp/append ./append
zip -j ./append.zip /tmp/append

sudo aws s3api create-bucket --bucket=mybucket-sdcc-lambda --region=us-east-1

sudo aws s3 cp ./put.zip s3://mybucket-sdcc-lambda/put.zip
sudo aws s3 cp ./get.zip s3://mybucket-sdcc-lambda/get.zip
sudo aws s3 cp ./delete.zip s3://mybucket-sdcc-lambda/delete.zip
sudo aws s3 cp ./append.zip s3://mybucket-sdcc-lambda/append.zip

cd ../terraform
sudo terraform init
sudo terraform apply
