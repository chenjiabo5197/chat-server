#!/bin/bash
workspace=$(cd $(dirname $0) && pwd -P)
export GOPATH=$workspace

gofmt -l -w -s src/
go build -o server src/main/main.go

