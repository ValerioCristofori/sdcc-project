#!/bin/bash

# delete all the backups
sudo docker volume rm sdcc-project_backup1
sudo docker volume rm sdcc-project_backup2
sudo docker volume rm sdcc-project_backup3
sudo docker volume rm sdcc-project_backup4

# build and run the cluster
sudo docker-compose build
sudo docker-compose up
sudo docker-compose down

