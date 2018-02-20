.DEFAULT_GOAL := help
Owner := youyo
Name := zabbix-aws-integration
Repository := "github.com/$(Owner)/$(Name)"
GithubToken := ${GITHUB_TOKEN}
Version := $(shell git describe --tags --abbrev=0)
CommitHash := $(shell git rev-parse --verify HEAD)
BuildTime := $(shell date '+%Y/%m/%d %H:%M:%S %Z')
GoVersion := $(shell go version)

## Setup
setup:
	go get -u -v github.com/golang/dep/cmd/dep
	go get -u -v github.com/mitchellh/gox
	go get -u -v github.com/tcnksm/ghr
	go get -u -v github.com/jstemmer/go-junit-report

## Install dependencies
deps:
	dep ensure -v

## Execute `go run`
run:
	go run \
		-ldflags "\
			-X \"$(Repository)/cmd/$(Name)/cmd.Name=$(Name)\" \
			-X \"$(Repository)/cmd/$(Name)/cmd.Version=$(Version)\" \
			-X \"$(Repository)/cmd/$(Name)/cmd.CommitHash=$(CommitHash)\" \
			-X \"$(Repository)/cmd/$(Name)/cmd.BuildTime=$(BuildTime)\" \
			-X \"$(Repository)/cmd/$(Name)/cmd.GoVersion=$(GoVersion)\"\
		" \
		./cmd/$(Name)/main.go ${OPTION}

## Build
build:
	gox -osarch="darwin/amd64 linux/amd64" \
		-ldflags="\
			-X \"$(Repository)/cmd/$(Name)/cmd.Name=$(Name)\" \
			-X \"$(Repository)/cmd/$(Name)/cmd.Version=$(Version)\" \
			-X \"$(Repository)/cmd/$(Name)/cmd.CommitHash=$(CommitHash)\" \
			-X \"$(Repository)/cmd/$(Name)/cmd.BuildTime=$(BuildTime)\" \
			-X \"$(Repository)/cmd/$(Name)/cmd.GoVersion=$(GoVersion)\"\
		" \
		-output="pkg/$(Name)_{{.OS}}_{{.Arch}}" \
		./cmd/$(Name)/

## Packaging
package:
	for arch in darwin_amd64 linux_amd64; do \
		zip -j pkg/$(Name)_$$arch.zip pkg/$(Name)_$$arch; \
		done

## Release
release:
	ghr -t ${GithubToken} -u $(Owner) -r $(Name) --replace $(Version) pkg/

## Remove packages
clean:
	rm -rf pkg/

## Show help
help:
	@make2help $(MAKEFILE_LIST)

.PHONY: setup deps run build release clean help
