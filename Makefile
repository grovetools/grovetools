# Root Makefile for Grove Ecosystem

PACKAGES = agentlogs core cx docgen flow grove grove-anthropic grove-gemini grove.nvim hooks nb notify project-tmpl-go skills tend nav
# GROVE-META:ADD-REPO:PACKAGES - Do not remove this comment, used by grove meta add-repo

BINARIES = aglogs core cx docgen flow grove grove-anthropic grove-gemini grove-nvim hooks nb notify skills tend nav
# GROVE-META:ADD-REPO:BINARIES - Do not remove this comment, used by grove meta add-repo

.PHONY: all build clean test fmt vet lint check test-e2e test-ecosystem dev build-all register schema setup help

# Build all packages
all: build

# Generate JSON schemas for all packages
schema:
	@echo "Generating JSON schemas..."
	@echo "Generating core schemas..."
	@$(MAKE) -C core schema
	@echo "Schema generation complete!"

build:
	@echo "Building all Grove packages..."
	@for pkg in $(PACKAGES); do \
		echo "Building $$pkg..."; \
		$(MAKE) -C $$pkg build; \
	done
	@echo "Build complete!"

# Clean all packages
clean:
	@echo "Cleaning all packages..."
	@for pkg in $(PACKAGES); do \
		echo "Cleaning $$pkg..."; \
		$(MAKE) -C $$pkg clean; \
	done
	@echo "Clean complete!"

# Run tests for all packages
test:
	@echo "Testing all packages..."
	@for pkg in $(PACKAGES); do \
		echo "Testing $$pkg..."; \
		$(MAKE) -C $$pkg test; \
	done
	@echo "All tests complete!"

# Format all packages
fmt:
	@echo "Formatting all packages..."
	@for pkg in $(PACKAGES); do \
		$(MAKE) -C $$pkg fmt; \
	done
	@echo "Formatting complete!"

# Run vet on all packages
vet:
	@echo "Running vet on all packages..."
	@for pkg in $(PACKAGES); do \
		$(MAKE) -C $$pkg vet; \
	done
	@echo "Vet complete!"

# Run linter on all packages
lint:
	@echo "Running linter on all packages..."
	@for pkg in $(PACKAGES); do \
		$(MAKE) -C $$pkg lint; \
	done
	@echo "Lint complete!"

# Run all checks
check: fmt vet lint test

# Run E2E tests for all packages that have them
test-e2e:
	@echo "Running E2E tests for all packages..."
	@for pkg in $(PACKAGES); do \
		if [ -f $$pkg/Makefile ] && grep -q "^test-e2e:" $$pkg/Makefile; then \
			echo "Running E2E tests for $$pkg..."; \
			$(MAKE) -C $$pkg test-e2e || { echo "E2E tests failed for $$pkg"; exit 1; }; \
		fi; \
	done
	@echo "All E2E tests complete!"

# Quick test of the ecosystem
test-ecosystem: build
	@echo "Running ecosystem test..."
	@./test-ecosystem.sh

# Development builds with race detector
dev:
	@echo "Building all packages with race detector..."
	@for pkg in grove cx flow tend nb notify agentlogs nav hooks; do \
		echo "Building $$pkg (dev)..."; \
		$(MAKE) -C $$pkg dev; \
	done

# Build for all platforms
build-all:
	@echo "Building for all platforms..."
	@for pkg in grove cx flow tend nb notify agentlogs nav hooks; do \
		echo "Cross-compiling $$pkg..."; \
		$(MAKE) -C $$pkg build-all; \
	done

# Register all binaries with grove dev
register: build
	@echo "Registering all Grove binaries with grove dev..."
	@for pkg in grove cx flow tend nb notify agentlogs nav hooks; do \
		echo "Registering $$pkg..."; \
		(cd $$pkg && ../grove/bin/grove dev link . --as $$(basename $$pkg)); \
	done
	@echo "Registration complete!"

# Set up development environment from source
setup:
	@./grove/scripts/setup-dev.sh

# Show help
help:
	@echo "Grove Ecosystem Makefile"
	@echo "======================="
	@echo ""
	@echo "Available targets:"
	@echo "  make build          - Build all packages"
	@echo "  make schema         - Generate JSON schemas for all packages"
	@echo "  make clean          - Clean all build artifacts"
	@echo "  make test           - Run tests for all packages"
	@echo "  make fmt            - Format all Go code"
	@echo "  make vet            - Run go vet on all packages"
	@echo "  make lint           - Run linter on all packages"
	@echo "  make check          - Run all checks (fmt, vet, lint, test)"
	@echo "  make test-e2e       - Run E2E tests for all packages that have them"
	@echo "  make test-ecosystem - Run the ecosystem integration test"
	@echo "  make dev            - Build with race detector"
	@echo "  make build-all      - Cross-compile for multiple platforms"
	@echo "  make register       - Register all binaries with grove dev"
	@echo "  make setup          - Set up dev environment from source"
	@echo "  make help           - Show this help"
	@echo ""
	@echo "Individual package targets:"
	@echo "  Run 'make help' in any subdirectory for package-specific targets"
