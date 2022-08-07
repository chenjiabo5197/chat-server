#!/bin/bash
workspace=$(cd $(dirname $0) && pwd -P)
export GOPATH=$workspace
export GO111MODULE=off

gofmt -l -w -s src/
go build -o chat-service src/main/main.go

