SHELL := $(shell which bash)
ENV = /usr/bin/env

.SHELLFLAGS = -c

.ONESHELL:
.NOTPARALLEL:
.EXPORT_ALL_VARIABLES:

.PHONY: all
.DEFAULT_GOAL := help

VERSION = `git describe --tags --always`
BUILD   = `date +%FT%T%z`

LDFLAGS = -w -s -X main.version=${VERSION} -X main.build=${BUILD}

build: clean ## Build binary
	go build -ldflags "${LDFLAGS}"

coverage: test ## Create coverage report
	go tool cover -func=coverage.txt
	go tool cover -html=coverage.txt

clean: ## Delete binary
	rm -f gnc-api-d

dependencies: ## Show dependencies state
	@go list -m -u all | column -t

help: ## Show Help
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-16s\033[0m %s\n", $$1, $$2}'

install: clean ## Install binary in GOPATH
	go install -ldflags "${LDFLAGS}"

test: ## Run tests
	go test -coverprofile=coverage.txt -ldflags "${LDFLAGS}" ./...
