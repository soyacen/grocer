# Makefile for grocer project

# Generate protobuf files for all API subdirectories
protoc:
	@echo "Generating protobuf files for all API subdirectories..."
	@for dir in api/*; do \
		if [ -d "$$dir" ] && [ -f "$$dir/protoc.sh" ]; then \
			echo "Running $$dir/protoc.sh"; \
			bash "$$dir/protoc.sh"; \
		fi \
	done

help:
	@echo "Available targets:"
	@echo "  protoc                    - Generate protobuf files for all API subdirectories"
	@echo "  help                      - Show this help message"

.PHONY: protoc help protoc-%