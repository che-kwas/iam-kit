.DEFAULT_GOAL := help

# ==============================================================================
# Build options

MAKEFLAGS    += --no-print-directory

# ==============================================================================
# Includes

include make-rules/tools.mk

# ==============================================================================
# Targets

## all: Build all.
.PHONY: all
all: lint gen test

## lint: Check syntax and styling of go sources.
.PHONY: lint
lint: tools.verify.golangci-lint
	go mod tidy -compat=1.18
	golangci-lint run ./...

## gen: Generate error code and document.
.PHONY: gen
gen: tools.verify.codegen
	codegen ./code
	codegen -doc -output ./error_code_base.md ./code

## test: Run unit test.
.PHONY: test
test:
	go test -race -cover -coverprofile=coverage.out ./...

## cover: Run unit test and get test coverage.
.PHONY: cover
cover: test
	sed -i '/mock_.*.go/d' coverage.out
	go tool cover -html=coverage.out -o coverage.html

## update: Update all modules.
.PHONY: update
update:
	go get -u ./...
	go mod tidy -compat=1.18

## clean: Remove all files that are created by building.
.PHONY: clean
clean:
	-rm -vrf coverage.*

## help: Show help info.
.PHONY: help
help: Makefile
	@echo "\nUsage: make <TARGETS> <OPTIONS> ...\n\nTargets:"
	@sed -n 's/^##//p' $< | column -t -s ':' | sed -e 's/^/ /'
