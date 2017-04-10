# SHELL := /bin/bash
# ROOT := $(shell pwd)

# export GOPATH ?= $(ROOT)
# export GOBIN ?= $(ROOT)/bin

# .DEFAULT_GOAL = help

# .PHONY: help
# help:
# 	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# .PHONY: init
# init: ## a helper for setting up golang projects
# 	mkdir bin && mkdir pkg && mkdir src; \
# 	go version;

# .PHONY: get
# get: ## a wrapper around `go get`
# 	cd src; \
# 	go get;
