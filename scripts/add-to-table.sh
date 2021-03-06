#!/bin/bash

id=$(uuidgen)
randomWord=$(shuf -n1 /usr/share/dict/words)
date=$(date +%s)

cat <<EOJSON > /tmp/add.json
[
{
"Put": {
"TableName" : "Records",
"Item" : {
"id":{"S":"${id}"},
"issueDate":{"S":"${date}"},
"randomWord":{"S":"${randomWord}"}
}
}
}
]
EOJSON

cat /tmp/add.json

aws --profile mine dynamodb transact-write-items --transact-items file:///tmp/add.json \
    --endpoint-url http://localhost:8000
