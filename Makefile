.DEFAULT_GOAL := help

# general define
BIN_INSTALL_DIR := $(HOME)/.local/bin
UNAME := $(shell uname)
PATH := $(BIN_INSTALL_DIR):$(PATH)


# buf install information
# 1. buf
BUF_RELEASE := https://github.com/bufbuild/buf/releases
BUF_VERSION := 1.0.0-rc6
BUF_NAME := buf
BUF_FULL_URL := $(BUF_RELEASE)/download/v$(BUF_VERSION)/buf-$(shell uname -s)-$(shell uname -m) -o $(BIN_INSTALL_DIR)/$(BUF_NAME)

# 2. protoc breaking name
PROTOC_BREAKING_NAME := protoc-gen-buf-breaking
PROTOC_BREAKING_FULL_URL := $(BUF_RELEASE)/download/v$(BUF_VERSION)/buf-$(shell uname -s)-$(shell uname -m) -o $(BIN_INSTALL_DIR)/$(PROTOC_BREAKING_NAME)

# 3. protoc lint name
PROTOC_LINT_NAME := protoc-gen-buf-lint
PROTOC_LINT_FULL_URL := $(BUF_RELEASE)/download/v$(BUF_VERSION)/buf-$(shell uname -s)-$(shell uname -m) -o $(BIN_INSTALL_DIR)/$(PROTOC_LINT_NAME)

# install
.PHONY: install
install: 
	mkdir -p "$(BIN_INSTALL_DIR)"
	make install.buf

.PHONY: install.buf
install.buf: 
	curl -sSL $(BUF_FULL_URL) -o "$(BIN_INSTALL_DIR)/$(BUF_NAME)"
	chmod +x "$(BIN_INSTALL_DIR)/$(BUF_NAME)"
	curl -sSL $(PROTOC_BREAKING_FULL_URL) -o "$(BIN_INSTALL_DIR)/$(PROTOC_BREAKING_NAME)"
	chmod +x "$(BIN_INSTALL_DIR)/$(PROTOC_BREAKING_NAME)"
	curl -sSL $(PROTOC_LINT_FULL_URL) -o "$(BIN_INSTALL_DIR)/$(PROTOC_LINT_NAME)"
	chmod +x "$(BIN_INSTALL_DIR)/$(PROTOC_LINT_NAME)"

.PHONY: uninstall
uninstall: 
	make uninstall.buf

.PHONY: uninstall.buf
uninstall.buf:
	rm -rf $(BIN_INSTALL_DIR)/*

.PHONY: gen.proto
gen.proto:
	cd client && rm -rf wsmsg
	export PATH="$${PATH}:$(BIN_INSTALL_DIR)"
	buf generate
	cp -r server/wsmsg client/wsmsg
	
# reserved for the later reference ##################################
.PHONY: install.go
install.go: install.go.notidy ## install go with dependencies
	cd server && go mod tidy

.PHONY: install.go.notidy
install.go.notidy: ## install go with dependencies but no tidy
	cd server && go mod download && grep _ ./cmd/tools/tools.go | cut -d' ' -f2 | sed 's/\r//' | xargs go install
	cd dbctl && go mod download

.PHONY: gen.all
gen.all: 
	cd server && wire ./...
	buf generate

.PHONY: test
test: test.go test.dart ## test go and dart files

.PHONY: format
format: format.dart format.go ## format Go and Dart Files

.PHONY: test.go
test.go: lint.go ## test go files
	cd server && go test ./...
	cd dbctl && go test ./...

.PHONY: format.go
format.go: ## format go files
	cd server && gci -w . && gofumpt -w -s . && go mod tidy
	cd dbctl && gci -w . && gofumpt -w -s . && go mod tidy

.PHONY: lint.go
lint.go: ## lint go files
	cd server && golangci-lint run
	cd dbctl && golangci-lint run

.PHONY: clean
clean: ## clean up proto generated files and mocks
	rm -rf ./client/lib/protos
	rm -rf ./server/pkg/pr12er/*.pb.go
	find ./client/test/ -name *.mocks.dart -delete

.PHONY: help
help: ## show available commands
	@printf "make commands:\n\n"
	@grep -E '^[a-zA-Z_.]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\t\033[1m%-30s\033[0m %s\n", $$1, $$2}'
