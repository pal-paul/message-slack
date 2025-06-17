SERVICE		?= $(shell basename `go list`)
VERSION		?= $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || cat $(PWD)/.version 2> /dev/null || echo v0)
PACKAGE		?= $(shell go list)
PACKAGES	?= $(shell go list ./...)
FILES		?= $(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: help clean fmt lint vet test build all

default: help

help:   ## show this help
	@echo 'usage: make [target] ...'
	@echo ''
	@echo 'targets:'
	@egrep '^(.+)\:\ .*##\ (.+)' ${MAKEFILE_LIST} | sed 's/:.*##/#/' | column -t -c 2 -s '#'

all:    ## clean, format, build and unit test
	make clean-all
	make build
	make test

install:    ## build and install go application executable
	go install -v ./...
	go install github.com/golang/mock/mockgen@v1.6.0

env:    ## Print useful environment variables to stdout
	echo $(CURDIR)
	echo $(SERVICE)
	echo $(PACKAGE)
	echo $(VERSION)

clean:  ## go clean
	go clean

clean-all:  ## remove all generated artifacts and clean all build artifacts
	go clean -i ./...

tools:  ## fetch and install all required tools

vet:    ## run go vet on the source files
	go vet ./...

doc:    ## generate godocs and start a local documentation webserver on port 8085

update-dependencies:    ## update golang dependencies
	dep ensure

generate-mocks:     ## generate mock code
	go generate ./...

build: generate-mocks ## generate all mocks and build the go code
	go build -o cmd/cmd  cmd/cmd.go

deploy: install build

test: ## run unit tests
	go test -v ./cmd/

test-integration: ## run integration tests
	go test -v -tags=integration ./cmd/

test-all: ## run all tests (unit + integration)
	go test -v ./cmd/
	go test -v -tags=integration ./cmd/

test-coverage: ## run tests with coverage
	go test -v -coverprofile=coverage.out ./cmd/
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-bench: ## run benchmark tests
	go test -bench=. -benchmem ./cmd/

test-race: ## run tests with race detection
	go test -race -v ./cmd/

test-clean: ## clean test artifacts
	rm -f coverage.out coverage.html
	rm -f cmd/cmd_test cmd/cmd_timeout_test

tidy:
	go get -u ./...
	go mod tidy