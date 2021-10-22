#!/bin/bash

sudo aws iam delete-role-policy --role-name myrole --policy-name lambda_policy
sudo aws iam delete-role --role-name myrole
