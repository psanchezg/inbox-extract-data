# Helper variables and Make settings
.PHONY: help clean build proto-link proto-vendor run
.DEFAULT_GOAL := help
.ONESHELL :
.SHELLFLAGS := -ec
SHELL := /bin/bash

help:                                  ## Print list of tasks
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_%-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' Makefile

build: clean                ## Build go project
	go build -o build/inbox-extract-data main.go

build-win: clean                ## Build go project
	GOARCH=386 GOOS=windows go build -o build/inbox-extract-data-386.exe main.go

build-win64: clean                ## Build go project
	GOARCH=amd64 GOOS=windows go build -o build/inbox-extract-data-amd64.exe main.go

build-arm: clean                ## Build go project
	GOARCH=arm GOARM=7 GOOS=linux go build -o build/inbox-extract-data-armv7 main.go

build-mac-arm64: clean
	GOOS=darwin GOARCH=arm64 go build -o build/inbox-extract-data-arm64-darwin main.go

build-mac-amd64: clean
	GOOS=darwin GOARCH=amd64 go build -o build/inbox-extract-data-amd64-darwin main.go

build-linux-arm64: clean
	GOOS=linux GOARCH=arm64 go build -o build/inbox-extract-data-arm64-linux main.go

build-linux-amd64: clean
	GOOS=linux GOARCH=amd64 go build -o build/inbox-extract-data-amd64-linux main.go

build-linux-386: clean
	GOOS=linux GOARCH=386 go build -o build/inbox-extract-data-386-linux main.go

clean:
	rm -rf vendor
	rm -rf proto
	rm -rf build/*

run:                                   ## Runs the demo server
	./build/inbox-extract-data-arm64-darwin

dev:
	go run main.go

test:
	go test -v

clean-token:
	rm -rf ~/.credentials/google.json

update:
	go mod tidy
	go mod vendor
	go build -mod=vendor