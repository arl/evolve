SHELL=/bin/bash

test: ## Run all tests, excluding those in '_example'.
	go test -race ./...

test-all: test ## Run all tests, including examples.
	cd _example && go test -race ./...

.PHONY: test test-all
