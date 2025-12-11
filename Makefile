# SSG Project Makefile
# Configuration variables (override via command line: make PORT=9000)

PORT ?= 8080
DEV_DIR ?= dev
BIN_DIR ?= bin
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo 'dev')
LDFLAGS := -X 'main.Version=$(VERSION)'

.PHONY: all test test-v test-cover build install serve serve-only kill assets fmt vet lint clean dev help changelog

# Composite targets
all: fmt vet test build

dev: kill assets serve-only

# Core build and test targets
test:
	go test ./...

test-v:
	go test -v ./...

test-cover:
	go test -coverprofile=coverage.out ./...

build:
	@mkdir -p $(BIN_DIR)
	go build -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/ssg ./cmd/ssg

install:
	go install -ldflags "$(LDFLAGS)" ./cmd/ssg

# Development server targets
serve: kill assets build
	cd $(DEV_DIR) && ../$(BIN_DIR)/ssg serve --build --port $(PORT)

serve-only:
	cd $(DEV_DIR) && ../$(BIN_DIR)/ssg serve --build --port $(PORT)

kill:
	@lsof -ti :$(PORT) | xargs kill -9 2>/dev/null || true

assets:
	@mkdir -p $(DEV_DIR)
	@cp -r assets $(DEV_DIR)/ 2>/dev/null || true

# Quality and maintenance targets
fmt:
	go fmt ./...

vet:
	go vet ./...

lint:
	@if which staticcheck > /dev/null 2>&1; then staticcheck ./...; else echo "staticcheck not installed, skipping lint"; fi

clean:
	rm -rf $(BIN_DIR)
	rm -f coverage.out

# Help target
help:
	@echo "SSG Makefile targets:"
	@echo ""
	@echo "  Composite:"
	@echo "    all       - Run fmt, vet, test, build (CI pipeline)"
	@echo "    dev       - Quick restart: kill, assets, serve-only"
	@echo ""
	@echo "  Build & Test:"
	@echo "    test      - Run tests"
	@echo "    test-v    - Run tests (verbose)"
	@echo "    test-cover- Run tests with coverage"
	@echo "    build     - Build binary to $(BIN_DIR)/ssg"
	@echo "    install   - Install binary via go install"
	@echo ""
	@echo "  Development:"
	@echo "    serve     - Full workflow: kill, assets, build, serve"
	@echo "    serve-only- Serve without rebuild"
	@echo "    kill      - Kill process on port $(PORT)"
	@echo "    assets    - Copy assets to $(DEV_DIR)/"
	@echo ""
	@echo "  Quality:"
	@echo "    fmt       - Format code"
	@echo "    vet       - Run go vet"
	@echo "    lint      - Run staticcheck"
	@echo "    clean     - Remove build artifacts"
	@echo "    changelog - Generate CHANGELOG.md from bd issues"
	@echo ""
	@echo "  Variables (override with VAR=value):"
	@echo "    PORT      - Server port (default: 8080)"
	@echo "    DEV_DIR   - Development directory (default: dev)"
	@echo "    BIN_DIR   - Binary output directory (default: bin)"

# Changelog generation
changelog:
	@echo "# Changelog" > CHANGELOG.md
	@echo "" >> CHANGELOG.md
	@echo "## [Unreleased]" >> CHANGELOG.md
	@echo "" >> CHANGELOG.md
	@echo "### Added" >> CHANGELOG.md
	@bd export --status=closed --type=feature 2>/dev/null | jq -r '"- " + .title + " (" + .id + ")"' >> CHANGELOG.md || echo "- None" >> CHANGELOG.md
	@echo "" >> CHANGELOG.md
	@echo "### Fixed" >> CHANGELOG.md
	@bd export --status=closed --type=bug 2>/dev/null | jq -r '"- " + .title + " (" + .id + ")"' >> CHANGELOG.md || echo "- None" >> CHANGELOG.md
