#!/bin/bash
aws --profile mine dynamodb create-table \
   --table-name Records \
   --attribute-definitions AttributeName=pk,AttributeType=S \
   --key-schema AttributeName=pk,KeyType=HASH \
   --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 \
   --endpoint-url http://localhost:8000
