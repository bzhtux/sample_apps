#!/usr/bin/env bash

set -euo pipefail

HOST="app-pg.127.0.0.1.nip.io"
PORT="80"

echo "Adding a new Book: The Hitchhiker's Guide to the Galaxy by Douglas Adams"
bookID=$(curl -sL -X POST -d '{"title": "The Hitchhiker'\'s' Guide to the Galaxy", "author": "Douglas Adams"}'  http://"${HOST}":"${PORT}"/add | jq .data.ID)
echo "New book has ID ${bookID}"


echo ""

echo "Adding twice the same book, expecting a conflict"
curl -sL -X POST -d '{"title": "The Hitchhiker'\'s' Guide to the Galaxy", "author": "Douglas Adams"}'  http://"${HOST}":"${PORT}"/add | jq .


echo ""

echo "Getting the book with ID: ${bookID}"
curl -sL http://"${HOST}":"${PORT}"/get/"${bookID}" | jq .data

echo ""

echo "Deleting book with ID ${bookID}"
curl -sL -X DELETE http://"${HOST}":"${PORT}"/del/"${bookID}" | jq .
