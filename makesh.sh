#!/bin/bash
export GOPATH="$(pwd)"
BIN_NAME="bin/swag"

if [[ "$OSTYPE" == "msys" ]]; then
  BIN_NAME="$BIN_NAME.exe"
elif [[ "$OSTYPE" == "win32" ]]; then
  BIN_NAME="$BIN_NAME.exe"
fi

go build -o $BIN_NAME src/swagit/swagit.go 