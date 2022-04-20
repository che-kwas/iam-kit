.DEFAULT_GOAL := help

# ==============================================================================
# Build options

OUTPUT_DIR   := ./_output
MAKEFLAGS    += --no-print-directory

# ==============================================================================
# Includes

include make-rules/tools.mk

# ==============================================================================
# Targets

## all: Build all.
.PHONY: all
all: tools lint test

## lint: Check syntax and styling of go sources.
.PHONY: lint
lint: tools.verify.golangci-lint
	go mod tidy -compat=1.17
	golangci-lint run ./...

## test: Run unit test.
.PHONY: test
test:
	@-mkdir -p $(OUTPUT_DIR)
	go test -race -cover -coverprofile=$(OUTPUT_DIR)/coverage.out ./...

## cover: Run unit test and get test coverage.
.PHONY: cover
cover: test
	sed -i '/.*_mock.go/d' $(OUTPUT_DIR)/coverage.out
	go tool cover -html=$(OUTPUT_DIR)/coverage.out -o $(OUTPUT_DIR)/coverage.html

## clean: Remove all files that are created by building.
.PHONY: clean
clean:
	-rm -vrf $(OUTPUT_DIR)

## tools: Install dependent tools.
.PHONY: tools
tools:
	make tools.verify

## help: Show help info.
.PHONY: help
help: Makefile
	@echo "\nUsage: make <TARGETS> <OPTIONS> ...\n\nTargets:"
	@sed -n 's/^##//p' $< | column -t -s ':' | sed -e 's/^/ /'
