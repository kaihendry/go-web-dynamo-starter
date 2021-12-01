#!/bin/bash
aws --profile mine dynamodb list-tables \
   --endpoint-url http://localhost:8000
