#!/bin/bash

id=$(uuidgen)
randomWord=$(shuf -n1 /usr/share/dict/words)
date=$(date --iso-8601=seconds)

cat <<EOJSON > /tmp/add.json
{
"id":"${id}",
"issueDate":"${date}",
"randomWord":"${randomWord}"
}
EOJSON

cat /tmp/add.json

curl -X POST -d @/tmp/add.json http://localhost:3000/add
