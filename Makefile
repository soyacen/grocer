# Makefile for grocer project

# Variables
VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")

# Copy proto files from pkg to internal/layout/third_party
copy-proto:
	@echo "Copying proto files from pkg to internal/layout/third_party..."
	@mkdir -p internal/layout/third_party
	@find pkg -name "*.proto" -type f -exec cp {} internal/layout/third_party/ \;
	@echo "Proto files copied successfully!"

# Generate protobuf files
protoc:
	@echo "Generating protobuf files..."
	bash ./scripts/protoc-all.sh

# Build the project
build:
	@echo "Building the project..."
	go build -ldflags="-X github.com/soyacen/grocer/cmd.Version=$(VERSION)" -o bin/grocer .

# Install the project
install:
	@echo "Installing the project..."
	go install -ldflags="-X github.com/soyacen/grocer/cmd.Version=$(VERSION)" .

help:
	@echo "Available targets:"
	@echo "  copy-proto - Copy proto files from pkg to internal/layout/third_party"
	@echo "  protoc     - Generate protobuf files"
	@echo "  build      - Build the project"
	@echo "  install    - Install the project"
	@echo "  help       - Show this help message"

.PHONY: copy-proto protoc build install help