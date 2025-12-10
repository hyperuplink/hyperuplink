.PHONY: all help test build db\:drop db\:dump db\:restore
PWD := $(shell pwd)
GOPATH := $(shell go env GOPATH)

GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

NAME := hyperuplink
PREFIX := github.com/mrusme/
PROJECT := $(PREFIX)$(NAME)
VERSION := $(shell git describe --tags 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --verify HEAD)
DATE := $(shell date)

all: build

help: ## print this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' Makefile | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-40s\033[0m %s\n", $$1, $$2}'

test: ## test
	go test -v ./...

build: ## build
	@echo "Building with the following parameters:"
	@echo "VERSION = $(VERSION)"
	@echo "COMMIT  = $(COMMIT)"
	@echo "DATE    = $(DATE)"
	@CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags "-s -w -X \"${PROJECT}/runtime.Version=${VERSION}\" -X \"${PROJECT}/runtime.Commit=${COMMIT}\" -X \"${PROJECT}/runtime.Date=${DATE}\"" -o $(PWD)/build/$(NAME)

db\:drop: ## drop (and recreate) development database
	psql -h localhost -p 5432 -U postgres -c "DROP DATABASE hyperuplink_dev;" -c "CREATE DATABASE hyperuplink_dev;"

db\:dump: ## dump current data from development database into dummy.sql
	pg_dump -h localhost -p 5432 -U postgres -d hyperuplink_dev -t categories -t forums --data-only --inserts -f dummy.sql

db\:restore: ## restore dummy.sql data
	psql -h localhost -p 5432 -U postgres -d hyperuplink_dev < dummy.sql

run: build ## build and run
	./build/hyperuplink -c "file://$(PWD)/hyperuplink.toml"
