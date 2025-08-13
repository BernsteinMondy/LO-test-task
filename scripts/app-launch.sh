#!/bin/bash

script_dir="$(dirname "$0")"

cd "${script_dir}/.."

if [ -f .env ]; then
    export $(grep -v '^#' .env | xargs)
else
    echo ".env file not found"
    exit 1
fi

cd ./src/cmd/server && go run .