#!/bin/bash
aws --profile mine dynamodb describe-table \
   --table-name Records \
   --endpoint-url http://localhost:8000
