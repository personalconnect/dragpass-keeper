VERSION := 1.0
EXTENSION_ID := cmgjlocmnppfpknaipdfodjhbplnhimk
PKG_IDENTIFIER := com.blindfold.keeper.pkg

MAC_BIN := mac-blindfold-keeper
MAC_PKG_DIR := output/macos
MAC_PKG := $(MAC_PKG_DIR)/Blindfold-Keeper.pkg
MAC_JSON_TPL := build_root/Library/Application\ Support/Google/Chrome/NativeMessagingHosts/com.blindfold.keeper.json

WIN_BIN := win-blindfold-keeper.exe
WIN_PKG_DIR := output/windows
WIN_PKG := $(WIN_PKG_DIR)/Blindfold-Keeper.exe
# Go 크로스 컴파일용 Windows 컴파일러
WIN_CC := x86_64-w64-mingw32-gcc

.PHONY: all build pkg clean build-macos build-windows pkg-macos pkg-windows

all: build pkg

build: build-macos build-windows
	@echo "All binaries built."

pkg: pkg-macos pkg-windows
	@echo "All installers packaged."

build-macos: $(MAC_BIN)
$(MAC_BIN): main.go go.mod
	@echo "Building macOS binary: $(MAC_BIN)..."
	@go build -o $(MAC_BIN) .

build-windows: $(WIN_BIN)
$(WIN_BIN): main.go go.mod
	@echo "Building Windows binary: $(WIN_BIN)..."
	@CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=$(WIN_CC) go build -o $(WIN_BIN) .

pkg-macos: $(MAC_PKG)
$(MAC_PKG): $(MAC_BIN)
	@echo "Creating macOS package structure in ./build_root..."
	@rm -rf build_root
	@mkdir -p build_root/Library/Application\ Support/Blindfold
	@mkdir -p build_root/Library/Application\ Support/Google/Chrome/NativeMessagingHosts	
	@cp $(MAC_BIN) build_root/Library/Application\ Support/Blindfold/
	@echo "{\n  \"name\": \"com.blindfold.keeper\",\n  \"description\": \"Blindfold Device Key Storage\",\n  \"path\": \"/Library/Application Support/Blindfold/$(MAC_BIN)\",\n  \"type\": \"stdio\",\n  \"allowed_origins\": [\n    \"chrome-extension://$(EXTENSION_ID)/\"\n  ]\n}" > $(MAC_JSON_TPL)
	
	@echo "Creating output directory: $(MAC_PKG_DIR)..."
	@mkdir -p $(MAC_PKG_DIR)

	@echo "Building macOS package: $(MAC_PKG)..."
	@pkgbuild --root ./build_root \
            --identifier $(PKG_IDENTIFIER) \
            --version $(VERSION) \
            $(MAC_PKG)
	@echo "Successfully built $(MAC_PKG)"

pkg-windows: $(WIN_PKG)
$(WIN_PKG): $(WIN_BIN) setup.iss
	@echo "Building Windows installer via Docker: $(WIN_PKG)..."
	@docker run --rm -v "$$PWD:/work" amake/innosetup setup.iss
	@echo "Successfully built Windows installer."

clean:
	@echo "Cleaning up build artifacts..."
	@rm -f $(MAC_BIN) $(WIN_BIN)
	@rm -rf build_root output