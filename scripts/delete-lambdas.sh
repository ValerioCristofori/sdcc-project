#!/bin/bash

sudo aws lambda delete-function --function-name put
sudo aws lambda delete-function --function-name get 
sudo aws lambda delete-function --function-name append 
sudo aws lambda delete-function --function-name delete 
