#!/bin/bash
cd ..
sudo docker-compose down

sh ./scripts/delete-table-dynamoDB.sh
cd ./terraform
sudo terraform destroy