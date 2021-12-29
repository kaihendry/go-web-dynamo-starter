#!/bin/bash
echo Use template.yml
exit 
aws --profile mine dynamodb create-table \
	--table-name Records \
	--attribute-definitions AttributeName=id,AttributeType=S \
	--key-schema AttributeName=id,KeyType=HASH \
	--billing-mode PAY_PER_REQUEST
