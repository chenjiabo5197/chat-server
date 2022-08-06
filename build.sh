#!/bin/bash
workspace=$(cd $(dirname $0) && pwd -P)
export GOPATH=$workspace

gofmt -l -w -s src/
go build -o client src/main/main.go

