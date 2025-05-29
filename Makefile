APP_NAME := goeditor
SRC_FILE := main.go
BUILD_DIR := output

GOOS_LIST := darwin linux windows
GOARCH_LIST := amd64 arm64

.PHONY: all deps build clean

# Default target
all: deps build

# Install OS-specific dependencies
deps:
	@echo "Installing dependencies for $(shell uname)..."
	@if [ "$$(uname)" = "Linux" ]; then \
		sudo apt-get update && \
		sudo apt-get install -y golang gcc libgl1-mesa-dev xorg-dev libxkbcommon-dev; \
	elif [ "$$(uname)" = "Darwin" ]; then \
		xcode-select --install || true; \
	else \
		echo "Dependency installation not supported on this OS"; \
	fi

# Build for all platforms
build:
	@mkdir -p $(BUILD_DIR)
	@for GOOS in $(GOOS_LIST); do \
		for GOARCH in $(GOARCH_LIST); do \
			if [ "$$GOOS" = "windows" ] && [ "$$GOARCH" = "arm64" ]; then continue; fi; \
			OUT_FILE="$(BUILD_DIR)/$(APP_NAME)-$$GOOS-$$GOARCH"; \
			[ "$$GOOS" = "windows" ] && OUT_FILE="$$OUT_FILE.exe"; \
			echo "Building for $$GOOS/$$GOARCH -> $$OUT_FILE"; \
			GOOS=$$GOOS GOARCH=$$GOARCH CGO_ENABLED=1 go build -o "$$OUT_FILE" $(SRC_FILE); \
		done; \
	done

# Clean build output
clean:
	@echo "Cleaning build directory..."
	@rm -rf $(BUILD_DIR)
