.PHONY: help build test lint cover bench clean check
.DEFAULT_GOAL := help
.SILENT:

help:
	@echo "\nAvailable Commands:\n"
	@echo " - help    :  Show this help message."
	@echo " - build   :  Compile the project."
	@echo " - test    :  Run all unit tests."
	@echo " - lint    :  Run all linting & style checks."
	@echo " - cover   :  Generate test coverage."
	@echo " - bench   :  Run all benchmarks."
	@echo " - clean   :  Clear build artifacts and project clutter."
	@echo " - check   :  Run all checking targets in one go."

build:
	@echo "\nBuilding: goignore\n"
	go mod tidy
	go build .

test:
	@echo "\nRunning Unit Tests\n"
	go test -race -cover ./...

lint:
	@echo "\nRunning: golangci-lint\n"
	golangci-lint run

cover:
	@echo "\nTest Coverage\n"
	go test -race -cover -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

bench:
	@echo "\nBenchmarks\n"
	go test -bench=.

clean:
	@echo "\nCleaning Project Clutter\n"
	go clean ./...
	rm -f coverage.out coverage.html

check: test bench lint cover