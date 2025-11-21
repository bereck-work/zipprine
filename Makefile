# zipprine Makefile
# Build for all major architectures

.PHONY: all build clean install test help deps build-all release

# Application name
BINARY_NAME=zipprine
VERSION?=1.0.0
BUILD_DIR=build
RELEASE_DIR=releases

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOINSTALL=$(GOCMD) install

# Build flags
LDFLAGS=-ldflags "-s -w -X main.Version=$(VERSION)"
BUILD_FLAGS=-trimpath

# Source
MAIN_PATH=./cmd/zipprine

# Color output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[0;33m
BLUE=\033[0;34m
MAGENTA=\033[0;35m
CYAN=\033[0;36m
NC=\033[0m # No Color

##@ General

help: ## Display this help screen
	@echo "$(CYAN)╔═══════════════════════════════════════════════════╗$(NC)"
	@echo "$(CYAN)║         🗜️  zipprine Build System 🚀              ║$(NC)"
	@echo "$(CYAN)╚═══════════════════════════════════════════════════╝$(NC)"
	@echo ""
	@awk 'BEGIN {FS = ":.*##"; printf "Usage:\n  make $(CYAN)<target>$(NC)\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  $(CYAN)%-15s$(NC) %s\n", $1, $2 } /^##@/ { printf "\n$(MAGENTA)%s$(NC)\n", substr($0, 5) } ' $(MAKEFILE_LIST)

##@ Development

deps: ## Download dependencies
	@echo "$(BLUE)📦 Downloading dependencies...$(NC)"
	@$(GOMOD) download
	@$(GOMOD) tidy
	@echo "$(GREEN)✅ Dependencies installed$(NC)"

test: ## Run tests
	@echo "$(BLUE)🧪 Running tests...$(NC)"
	@$(GOTEST) -race -timeout 30s ./...
	@echo "$(GREEN)✅ Tests passed$(NC)"

test-verbose: ## Run tests with verbose output
	@echo "$(BLUE)🧪 Running tests (verbose)...$(NC)"
	@$(GOTEST) -v -race -timeout 30s ./...
	@echo "$(GREEN)✅ Tests passed$(NC)"

test-coverage: ## Run tests with coverage report
	@echo "$(BLUE)📊 Running tests with coverage...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@$(GOTEST) -race -coverprofile=$(BUILD_DIR)/coverage.out -covermode=atomic ./...
	@$(GOCMD) tool cover -html=$(BUILD_DIR)/coverage.out -o $(BUILD_DIR)/coverage.html
	@echo "$(GREEN)✅ Coverage report generated: $(BUILD_DIR)/coverage.html$(NC)"
	@$(GOCMD) tool cover -func=$(BUILD_DIR)/coverage.out | tail -n 1

bench: ## Run benchmarks
	@echo "$(BLUE)📊 Running benchmarks...$(NC)"
	@$(GOTEST) -bench=. -benchmem -run=^$$ ./...
	@echo "$(GREEN)✅ Benchmarks complete$(NC)"

test-all: test-coverage bench ## Run all tests with coverage and benchmarks
	@echo "$(GREEN)✨ All tests and benchmarks complete$(NC)"

clean: ## Clean build artifacts
	@echo "$(YELLOW)🧹 Cleaning build artifacts...$(NC)"
	@$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@rm -rf $(RELEASE_DIR)
	@echo "$(GREEN)✅ Clean complete$(NC)"

fmt: ## Format code
	@echo "$(BLUE)📝 Formatting code...$(NC)"
	@$(GOCMD) fmt ./...
	@echo "$(GREEN)✅ Code formatted$(NC)"

vet: ## Run go vet
	@echo "$(BLUE)🔍 Running go vet...$(NC)"
	@$(GOCMD) vet ./...
	@echo "$(GREEN)✅ Vet complete$(NC)"

lint: fmt vet ## Run linters
	@echo "$(GREEN)✅ Linting complete$(NC)"

##@ Building

build: deps ## Build for current platform
	@echo "$(BLUE)🔨 Building $(BINARY_NAME) for current platform...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@$(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "$(GREEN)✅ Build complete: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

install: deps ## Install to $GOPATH/bin
	@echo "$(BLUE)📥 Installing $(BINARY_NAME)...$(NC)"
	@$(GOINSTALL) $(LDFLAGS) $(MAIN_PATH)
	@echo "$(GREEN)✅ Installed to $(shell go env GOPATH)/bin/$(BINARY_NAME)$(NC)"

run: build ## Build and run
	@echo "$(CYAN)▶️  Running $(BINARY_NAME)...$(NC)"
	@./$(BUILD_DIR)/$(BINARY_NAME)

##@ Cross-Platform Builds

build-linux-amd64: ## Build for Linux AMD64
	@echo "$(BLUE)🐧 Building for Linux AMD64...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	@echo "$(GREEN)✅ Built: $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64$(NC)"

build-linux-arm64: ## Build for Linux ARM64
	@echo "$(BLUE)🐧 Building for Linux ARM64...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=arm64 $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_PATH)
	@echo "$(GREEN)✅ Built: $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64$(NC)"

build-linux-arm: ## Build for Linux ARM
	@echo "$(BLUE)🐧 Building for Linux ARM...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=arm $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm $(MAIN_PATH)
	@echo "$(GREEN)✅ Built: $(BUILD_DIR)/$(BINARY_NAME)-linux-arm$(NC)"

build-darwin-amd64: ## Build for macOS AMD64 (Intel)
	@echo "$(BLUE)🍎 Building for macOS AMD64...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@GOOS=darwin GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	@echo "$(GREEN)✅ Built: $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64$(NC)"

build-darwin-arm64: ## Build for macOS ARM64 (Apple Silicon)
	@echo "$(BLUE)🍎 Building for macOS ARM64...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@GOOS=darwin GOARCH=arm64 $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	@echo "$(GREEN)✅ Built: $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64$(NC)"

build-windows-amd64: ## Build for Windows AMD64
	@echo "$(BLUE)🪟 Building for Windows AMD64...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@GOOS=windows GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	@echo "$(GREEN)✅ Built: $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe$(NC)"

build-windows-arm64: ## Build for Windows ARM64
	@echo "$(BLUE)🪟 Building for Windows ARM64...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@GOOS=windows GOARCH=arm64 $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-arm64.exe $(MAIN_PATH)
	@echo "$(GREEN)✅ Built: $(BUILD_DIR)/$(BINARY_NAME)-windows-arm64.exe$(NC)"

build-freebsd-amd64: ## Build for FreeBSD AMD64
	@echo "$(BLUE)👹 Building for FreeBSD AMD64...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@GOOS=freebsd GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-freebsd-amd64 $(MAIN_PATH)
	@echo "$(GREEN)✅ Built: $(BUILD_DIR)/$(BINARY_NAME)-freebsd-amd64$(NC)"

build-openbsd-amd64: ## Build for OpenBSD AMD64
	@echo "$(BLUE)🐡 Building for OpenBSD AMD64...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@GOOS=openbsd GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-openbsd-amd64 $(MAIN_PATH)
	@echo "$(GREEN)✅ Built: $(BUILD_DIR)/$(BINARY_NAME)-openbsd-amd64$(NC)"

build-all: deps ## Build for all platforms
	@echo "$(MAGENTA)╔═══════════════════════════════════════════════════╗$(NC)"
	@echo "$(MAGENTA)║     🌍 Building for ALL architectures 🚀         ║$(NC)"
	@echo "$(MAGENTA)╚═══════════════════════════════════════════════════╝$(NC)"
	@echo ""
	@$(MAKE) build-linux-amd64
	@$(MAKE) build-linux-arm64
	@$(MAKE) build-linux-arm
	@$(MAKE) build-darwin-amd64
	@$(MAKE) build-darwin-arm64
	@$(MAKE) build-windows-amd64
	@$(MAKE) build-windows-arm64
	@$(MAKE) build-freebsd-amd64
	@$(MAKE) build-openbsd-amd64
	@echo ""
	@echo "$(GREEN)╔═══════════════════════════════════════════════════╗$(NC)"
	@echo "$(GREEN)║          ✨ All builds complete! ✨              ║$(NC)"
	@echo "$(GREEN)╚═══════════════════════════════════════════════════╝$(NC)"
	@ls -lh $(BUILD_DIR)/

##@ Release

release: clean build-all ## Create release packages
	@echo "$(MAGENTA)📦 Creating release packages...$(NC)"
	@mkdir -p $(RELEASE_DIR)
	
	@echo "$(BLUE)  → Packaging Linux AMD64...$(NC)"
	@tar -czf $(RELEASE_DIR)/$(BINARY_NAME)-$(VERSION)-linux-amd64.tar.gz -C $(BUILD_DIR) $(BINARY_NAME)-linux-amd64
	
	@echo "$(BLUE)  → Packaging Linux ARM64...$(NC)"
	@tar -czf $(RELEASE_DIR)/$(BINARY_NAME)-$(VERSION)-linux-arm64.tar.gz -C $(BUILD_DIR) $(BINARY_NAME)-linux-arm64
	
	@echo "$(BLUE)  → Packaging Linux ARM...$(NC)"
	@tar -czf $(RELEASE_DIR)/$(BINARY_NAME)-$(VERSION)-linux-arm.tar.gz -C $(BUILD_DIR) $(BINARY_NAME)-linux-arm
	
	@echo "$(BLUE)  → Packaging macOS AMD64...$(NC)"
	@tar -czf $(RELEASE_DIR)/$(BINARY_NAME)-$(VERSION)-darwin-amd64.tar.gz -C $(BUILD_DIR) $(BINARY_NAME)-darwin-amd64
	
	@echo "$(BLUE)  → Packaging macOS ARM64...$(NC)"
	@tar -czf $(RELEASE_DIR)/$(BINARY_NAME)-$(VERSION)-darwin-arm64.tar.gz -C $(BUILD_DIR) $(BINARY_NAME)-darwin-arm64
	
	@echo "$(BLUE)  → Packaging Windows AMD64...$(NC)"
	@cd $(BUILD_DIR) && zip -q ../$(RELEASE_DIR)/$(BINARY_NAME)-$(VERSION)-windows-amd64.zip $(BINARY_NAME)-windows-amd64.exe
	
	@echo "$(BLUE)  → Packaging Windows ARM64...$(NC)"
	@cd $(BUILD_DIR) && zip -q ../$(RELEASE_DIR)/$(BINARY_NAME)-$(VERSION)-windows-arm64.zip $(BINARY_NAME)-windows-arm64.exe
	
	@echo "$(BLUE)  → Packaging FreeBSD AMD64...$(NC)"
	@tar -czf $(RELEASE_DIR)/$(BINARY_NAME)-$(VERSION)-freebsd-amd64.tar.gz -C $(BUILD_DIR) $(BINARY_NAME)-freebsd-amd64
	
	@echo "$(BLUE)  → Packaging OpenBSD AMD64...$(NC)"
	@tar -czf $(RELEASE_DIR)/$(BINARY_NAME)-$(VERSION)-openbsd-amd64.tar.gz -C $(BUILD_DIR) $(BINARY_NAME)-openbsd-amd64
	
	@echo ""
	@echo "$(GREEN)✅ Release packages created:$(NC)"
	@ls -lh $(RELEASE_DIR)/
	@echo ""
	@echo "$(CYAN)📊 Package sizes:$(NC)"
	@du -sh $(RELEASE_DIR)/*

checksums: ## Generate SHA256 checksums for releases
	@echo "$(BLUE)🔐 Generating checksums...$(NC)"
	@cd $(RELEASE_DIR) && shasum -a 256 * > SHA256SUMS
	@echo "$(GREEN)✅ Checksums generated: $(RELEASE_DIR)/SHA256SUMS$(NC)"
	@cat $(RELEASE_DIR)/SHA256SUMS

##@ Docker (Bonus)

docker-build: ## Build Docker image
	@echo "$(BLUE)🐳 Building Docker image...$(NC)"
	@docker build -t $(BINARY_NAME):$(VERSION) -t $(BINARY_NAME):latest .
	@echo "$(GREEN)✅ Docker image built$(NC)"

docker-run: ## Run in Docker
	@echo "$(CYAN)🐳 Running in Docker...$(NC)"
	@docker run -it --rm $(BINARY_NAME):latest

##@ Info

version: ## Show version
	@echo "$(CYAN)zipprine version: $(VERSION)$(NC)"

platforms: ## Show supported platforms
	@echo "$(CYAN)Supported platforms:$(NC)"
	@echo "  🐧 Linux:   AMD64, ARM64, ARM"
	@echo "  🍎 macOS:   AMD64 (Intel), ARM64 (Apple Silicon)"
	@echo "  🪟 Windows: AMD64, ARM64"
	@echo "  👹 FreeBSD: AMD64"
	@echo "  🐡 OpenBSD: AMD64"

size: ## Show binary sizes
	@echo "$(CYAN)Binary sizes:$(NC)"
	@if [ -d "$(BUILD_DIR)" ]; then \
		du -sh $(BUILD_DIR)/* | sort -h; \
	else \
		echo "$(RED)No builds found. Run 'make build-all' first.$(NC)"; \
	fi

.DEFAULT_GOAL := help