#!/bin/bash
aws --profile mine dynamodb create-table \
	--table-name Records \
	--attribute-definitions AttributeName=id,AttributeType=S AttributeName=created,AttributeType=N \
	--key-schema AttributeName=id,KeyType=HASH AttributeName=created,KeyType=RANGE \
	--billing-mode PAY_PER_REQUEST \
	--endpoint-url http://localhost:8000

# Not working and I don't know why
# --time-to-live-specification Enabled=true,AttributeName=expires \
