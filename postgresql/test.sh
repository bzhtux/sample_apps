#!/usr/bin/env bash

set -euo pipefail
HOST="${APP_NAME}.127.0.0.1.nip.io"
PORT="80"

addBook(){
    echo -ne "Adding a new Book: The Hitchhiker's Guide to the Galaxy by Douglas Adams\n"
    bookID=$(curl -sL -X POST -d '{"title": "The Hitchhiker'\'s' Guide to the Galaxy", "author": "Douglas Adams"}'  http://"${HOST}":"${PORT}"/add | jq .data.ID)
    echo -ne "New book has ID ${bookID}\n"
}

addTwiceBook(){
    echo -ne "Adding twice the same book, expecting a conflict\n"
    curl -sL -X POST -d '{"title": "The Hitchhiker'\'s' Guide to the Galaxy", "author": "Douglas Adams"}'  http://"${HOST}":"${PORT}"/add | jq .
}

getBook(){
    echo -ne "Getting the book with ID: ${bookID}\n"
    curl -sL http://"${HOST}":"${PORT}"/get/"${bookID}" | jq .data >&2
}

deleteBook(){
    echo -ne "Deleting book with ID ${bookID}\n"
    curl -sL -X DELETE http://"${HOST}":"${PORT}"/del/"${bookID}" | jq . >&2
}

runTests(){
    addBook
    echo ""
    sleep 1

    addTwiceBook
    echo ""
    sleep 1

    getBook
    echo ""
    sleep 1

    deleteBook
}
