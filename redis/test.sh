#!/usr/bin/env bash

set -euo pipefail

HOST="app-redis.127.0.0.1.nip.io"
PORT="80"

echo "Adding a new key: key1=val1"
curl -sL -X POST -d '{"key": "key1", "value": "val1"}'  http://"${HOST}":"${PORT}"/add | jq .
# curl -sL -X POST -d '{"key": "Nathan", "value": "Le BG"}'  http://"${HOST}":"${PORT}"/add | jq .

echo ""

echo "Adding twice the same key, expecting a conflict"
curl -sL -X POST -d '{"key": "key1", "value": "val1"}'  http://"${HOST}":"${PORT}"/add | jq .
# curl -sL -X POST -d '{"key": "Nathan", "value": "Le BG"}'  http://"${HOST}":"${PORT}"/add | jq .

echo ""

echo "Getting the previous key: key1"
curl -sL http://"${HOST}":"${PORT}"/get/key1 | jq .
# curl -sL http://"${HOST}":"${PORT}"/get/Nathan | jq .

echo ""

echo "Deleting key1"
curl -sL -X DELETE http://"${HOST}":"${PORT}"/del/key1 | jq .
# curl -sL -X DELETE http://"${HOST}":"${PORT}"/del/Nathan | jq .
