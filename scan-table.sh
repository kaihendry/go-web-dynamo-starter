#!/bin/bash
aws --profile mine dynamodb scan --table-name Records \
   --endpoint-url http://localhost:8000
