.PHONY: all build run gotool clean help

all: gotool build

BINARY="bluebell"

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/bluebell



build_run:
	@go run ./main.go conf/config.yaml

run:
	./${BINARY} conf/config.yaml