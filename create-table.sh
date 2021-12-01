#!/bin/bash
aws --profile mine dynamodb create-table \
   --table-name Records \
   --attribute-definitions AttributeName=id,AttributeType=S AttributeName=issueDate,AttributeType=S \
   --key-schema AttributeName=id,KeyType=HASH AttributeName=issueDate,KeyType=RANGE \
   --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 \
   --endpoint-url http://localhost:8000
