# Makefile for grocer project

# Generate protobuf files
protoc:
	@echo "Generating protobuf files..."
	bash ./scripts/protoc-all.sh

# Build the project
build:
	@echo "Building the project..."
	go build -o grocer .

# Install the project
install:
	@echo "Installing the project..."
	go install .

help:
	@echo "Available targets:"
	@echo "  protoc  - Generate protobuf files"
	@echo "  build   - Build the project"
	@echo "  install - Install the project"
	@echo "  help    - Show this help message"

.PHONY: protoc build install help