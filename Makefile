.PHONY: build build-all clean test help

VERSION ?= dev
BINARY_NAME = easytime
BUILD_DIR = build

# 定义编译目标
TARGETS = \
	linux/amd64 \
	linux/arm64 \
	linux/386 \
	linux/arm \
	windows/amd64 \
	windows/386 \
	windows/arm64 \
	darwin/amd64 \
	darwin/arm64

help:
	@echo "Available targets:"
	@echo "  make build          - Build for current OS/ARCH"
	@echo "  make build-all      - Build for all platforms"
	@echo "  make test           - Run tests"
	@echo "  make clean          - Clean build directory"

build:
	@mkdir -p $(BUILD_DIR)
	@echo "Building $(BINARY_NAME) for $(GOOS)/$(GOARCH)..."
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) -ldflags="-s -w" .

build-all: clean
	@mkdir -p $(BUILD_DIR)
	@echo "Building for all platforms..."
	$(foreach target,$(TARGETS),$(call build-target,$(target)))

define build-target
	@os=$(word 1,$(subst /, ,$(1))); \
	arch=$(word 2,$(subst /, ,$(1))); \
	ext=""; \
	[ "$$os" = "windows" ] && ext=".exe"; \
	echo "Building for $$os/$$arch..."; \
	GOOS=$$os GOARCH=$$arch go build -o $(BUILD_DIR)/$(BINARY_NAME)-$$os-$$arch$$ext -ldflags="-s -w" .
endef

test:
	@go test -v -race -coverprofile=coverage.out ./...

clean:
	@echo "Cleaning build directory..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out
