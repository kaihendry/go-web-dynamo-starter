#!/bin/bash

id=$(uuidgen)
randomWord=$(shuf -n1 /usr/share/dict/words)
date=$(date --iso-8601=seconds)

cat <<EOJSON > /tmp/add.json
[
{
"Put": {
"TableName" : "Records",
"Item" : {
"Id":{"S":"${id}"},
"IssueDate":{"S":"${date}"},
"RandomWord":{"S":"${randomWord}"}
}
}
}
]
EOJSON

cat /tmp/add.json

aws --profile mine dynamodb transact-write-items --transact-items file:///tmp/add.json \
   --endpoint-url http://localhost:8000
