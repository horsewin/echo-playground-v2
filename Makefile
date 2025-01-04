SHELL := /bin/sh

# Goコマンド
GOCMD         = go
# ビルド後の出力先ディレクトリ
BUILD_DIR     = bin
# ビルドするバイナリの名前
BINARY_NAME   = main
# メインのエントリーポイント
MAIN_PACKAGE  = ./main.go

.PHONY: all validate build test run clean update-deps build-linux

# allターゲットでは「validate → build → run」を一括実行
all: validate build run

##
# 検証系: fmt, vet, test
##
validate:
	@echo "==> Running go fmt"
	@$(GOCMD) fmt ./...

	@echo "==> Running go vet"
	@$(GOCMD) vet ./...

	@echo "==> Running golangci-lint"
	@golangci-lint run

#   必要に応じてテストを実行する
# 	@echo "==> Running tests"
# 	@$(GOCMD) test -v ./...

##
# ビルド & 実行
##
build:
	@echo "==> Building..."
	@mkdir -p $(BUILD_DIR)
	@$(GOCMD) build -ldflags "-s -w" -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)

run:
	@echo "==> Running..."
	@$(BUILD_DIR)/$(BINARY_NAME)

##
# テスト (validate でも実行しているが、個別でも呼び出せるように)
##
test:
	@echo "==> Running tests"
	@$(GOCMD) test -v ./...

##
# 不要ファイルの削除
##
clean:
	@echo "==> Cleaning build outputs"
	@rm -rf $(BUILD_DIR)

##
# 依存関係の更新: go mod tidy
##
update-deps:
	@echo "==> Updating dependencies"
	@$(GOCMD) mod tidy

##
# Linux向けクロスコンパイル
##
build-linux:
	@echo "==> Cross compiling for Linux (amd64)"
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(MAKE) build
