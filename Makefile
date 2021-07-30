.PHONY: build build-plugin check fmt lint test test-race vet test-cover-html help build-local
.DEFAULT_GOAL := build-neo
GOVERSION := $(shell go version | cut -d ' ' -f 3 | cut -d '.' -f 2)
SHELL := /usr/bin/env bash
ROOT := $(shell pwd)
TASKS := $(shell [ -d ${ROOT}/task ] && ls ${ROOT}/task)
HOOKS := $(shell [ -d ${ROOT}/hook ] && ls ${ROOT}/hook)

build: test build-gorelease  ## build all
	@echo " > build finished"

build-gorelease:
	goreleaser --snapshot --skip-publish --rm-dist

build-neo:
	@echo " > building binaries"
	@mkdir -p ./dist
	cd ./task/neo && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ../../dist/optimus-neo_linux_amd64 ./main.go
	cd ./task/neo && GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o ../../dist/optimus-neo_darwin_amd64 ./main.go
	@echo " > build complete"

install: ## install plugin to optimus directory
	mkdir -p ~/.optimus/plugins
	cp ./dist/neo_darwin_amd64/* ~/.optimus/plugins/

clean: ## clean binaries
	rm -rf ./dist

fmt: ## Run FMT
	@for target in ${TASKS}; do \
	  cd ${ROOT}/task/$${target}; go fmt . ; \
	done
	@for target in ${HOOKS}; do \
	  cd ${ROOT}/hook/$${target}; go fmt . ; \
	done

test: ## Run tests
	@if [ -d ${ROOT}/task ]; then \
	for target in ${TASKS}; do \
	  cd ${ROOT}/task/$${target}; go vet . ; go test . -race; \
	done \
	fi
	@if [ -d ${ROOT}/hook ]; then \
	for target in ${HOOKS}; do \
	  cd ${ROOT}/hook/$${target}; go vet . ; go test . -race; \
	done \
	fi

help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'