# ==============================================================================
#  CONFIGURATION
# ==============================================================================
APP_NAME    := my-app
DIST_DIR    := dist

# Paths for Go
# Point to the MAIN entry point folder
GO_ENTRY    := ./cmd/git-markdown-editor
GO_FILES    := .

# Paths for Tailwind (Moved into internal/ui)
TW_BINARY   := tailwindcss
INPUT_CSS   := internal/ui/assets/input.css
OUTPUT_CSS  := internal/ui/assets/output.css

# ==============================================================================
#  SYSTEM DETECTION (Linux/Mac/Windows)
# ==============================================================================
UNAME_S := $(shell uname -s)
OS      := linux
EXT     := 

ifeq ($(UNAME_S),Darwin)
	OS := macos
endif
ifneq (,$(findstring MINGW,$(UNAME_S)))
	OS := windows
	EXT := .exe
endif

ARCH_RAW := $(shell uname -m)
ARCH     := x64
ifeq ($(ARCH_RAW),aarch64)
	ARCH := arm64
endif
ifeq ($(ARCH_RAW),arm64)
	ARCH := arm64
endif

DOWNLOAD_URL := https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-$(OS)-$(ARCH)$(EXT)
BINARY_FILE  := $(TW_BINARY)$(EXT)

# ==============================================================================
#  TARGETS
# ==============================================================================
.PHONY: all install dev run watch prod-embed clean

all: prod-embed

install:
	 @echo "⬇️  Downloading Tailwind..."
	 @curl -sL -o $(BINARY_FILE) $(DOWNLOAD_URL)
	 @chmod +x $(BINARY_FILE)

# --- Development ---
dev:
	 @echo "🚀 Starting Dev Mode (cmd/server)..."
	 @$(MAKE) -j2 watch run

watch:
	 @./$(BINARY_FILE) -i $(INPUT_CSS) -o $(OUTPUT_CSS) --watch

run:
	 @go run $(GO_ENTRY)

# --- Production (Embedded) ---
prod-embed: clean build-css
	 @echo "📦 Building Single Embedded Binary..."
	 @mkdir -p $(DIST_DIR)
	
	# We build the package at ./cmd/server
	 @GOTOOLCHAIN=auto go build -ldflags="-w -s" -o $(DIST_DIR)/$(APP_NAME)$(EXT) $(GO_ENTRY)
	
	 @echo "✅ Build complete: ./$(DIST_DIR)/$(APP_NAME)$(EXT)"

build-css:
	 @echo "🎨 Minifying CSS..."
	 @./$(BINARY_FILE) -i $(INPUT_CSS) -o $(OUTPUT_CSS) --minify

clean:
	 @rm -rf $(DIST_DIR)
	 @rm -f $(BINARY_FILE)