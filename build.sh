#!/bin/bash

set -ux

NAME=${1}

go get -v github.com/golang/dep/cmd/dep
go get -v github.com/Songmu/make2help/cmd/make2help
rm -rf vendor/ artifacts/*
dep ensure -v
go build -o artifacts/${NAME} -x
rm -rf vendor/
