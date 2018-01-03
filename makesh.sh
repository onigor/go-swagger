#!/bin/bash
export GOPATH="$(pwd)"
BIN_NAME="swag.exe"
go build -o $BIN_NAME src/swagit/swagit.go 