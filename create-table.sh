#!/bin/bash
aws --profile mine dynamodb create-table \
   --table-name Records \
   --attribute-definitions AttributeName=Id,AttributeType=S \
   --key-schema AttributeName=Id,KeyType=HASH \
   --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 \
   --endpoint-url http://localhost:8000
