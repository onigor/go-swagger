GOPATH=$(pwd)
BIN_NAME="swag.exe"
all:build

build:
	go GOPATH=$(GOPATH) build src/swagit/swagit.go -o $(BIN_NAME)