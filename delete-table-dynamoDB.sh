#!/bin/bash
sudo aws configure set region us-east-1
sudo aws dynamodb delete-table --table-name Sensors