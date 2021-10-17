#!/bin/bash
cd ..
# build and run the cluster
sudo docker-compose build
sudo docker-compose up
sudo docker-compose down
