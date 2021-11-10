#!/bin/bash
cd ..
sudo docker-compose down

sh ./scripts/delete-table-dynamoDB.sh
cd ./terraform
sudo terraform destroy

sudo aws s3 rm s3://mybucket-sdcc-lambda1 --recursive
sudo aws s3api delete-bucket --bucket mybucket-sdcc-lambda1 --region us-east-1
