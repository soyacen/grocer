# Makefile for grocer project

# Variables
VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
COMMIT := $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")

# Generate protobuf files
protoc:
	@echo "Generating protobuf files..."
	bash ./scripts/protoc-all.sh

# Build the project
build:
	@echo "Building the project..."
	go build -ldflags="-X github.com/soyacen/grocer/cmd.Version=$(VERSION) -X github.com/soyacen/grocer/cmd.Commit=$(COMMIT)" -o bin/grocer .

# Install the project
install:
	@echo "Installing the project..."
	go install -ldflags="-X github.com/soyacen/grocer/cmd.Version=$(VERSION) -X github.com/soyacen/grocer/cmd.Commit=$(COMMIT)" .

help:
	@echo "Available targets:"
	@echo "  protoc  - Generate protobuf files"
	@echo "  build   - Build the project"
	@echo "  install - Install the project"
	@echo "  help    - Show this help message"

.PHONY: protoc build install help