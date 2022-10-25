#!/usr/bin/env bash

set -euo pipefail

HOST="gomongo.127.0.0.1.nip.io"
PORT="80"

# shellcheck disable=SC1073
while [ "$(curl -sL http://${HOST}:${PORT}/ping | jq -r .status)" != "Ok" ];
do
    sleep 3
done

echo "*** Create a new doc:"
echo -n '{"Title": "Hello world", "Author":"bzhtux"}' | jq .

docID=$(curl -sL -X POST -d '{"Title": "Hello world", "Author":"bzhtux"}' http://${HOST}:${PORT}/add | jq -r .data.ID)    
echo "*** New doc has id ${docID}"

echo "*** Create twice the same doc, expecting a conflict:"
curl -sL -X POST -d '{"Title": "Hello world", "Author":"bzhtux"}' http://${HOST}:${PORT}/add | jq .

echo "*** Create new doc with missing values, expecting an error"
curl -sL -X POST -d '{"Title": "Hello world"}' http://${HOST}:${PORT}/add | jq .

echo "*** Get one doc by Name: Hello world"
curl -sL http://${HOST}:${PORT}/get/byName/Hello%20world | jq .

echo "*** Get one doc by ID: $docID"
curl -sL http://${HOST}:${PORT}/get/byID/"${docID}" | jq .

echo "*** Delete one doc by Name: Hello world"
curl -sL -X DELETE http://${HOST}:${PORT}/del/byName/Hello%20world | jq .
