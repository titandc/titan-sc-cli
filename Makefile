BINARY_NAME := titan-sc

## —— Build ———————————————————————————————————————————————————————————————————
.PHONY: build
build: ## Build the CLI binary
	@go build -o $(BINARY_NAME) .

.PHONY: install
install: build ## Install to /usr/local/bin
	@sudo cp $(BINARY_NAME) /usr/local/bin/

.PHONY: build-all
build-all: ## Build for all platforms
	@mkdir -p dist
	GOOS=linux GOARCH=amd64 go build -o dist/$(BINARY_NAME)-linux-amd64 .
	GOOS=linux GOARCH=386 go build -o dist/$(BINARY_NAME)-linux-i386 .
	GOOS=linux GOARCH=arm64 go build -o dist/$(BINARY_NAME)-linux-arm64 .
	GOOS=darwin GOARCH=amd64 go build -o dist/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -o dist/$(BINARY_NAME)-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build -o dist/$(BINARY_NAME)-windows-amd64.exe .
	GOOS=windows GOARCH=386 go build -o dist/$(BINARY_NAME)-windows-i386.exe .

.PHONY: release
release: ## Build and package releases for all platforms
	@mkdir -p release
	@# Linux amd64
	GOOS=linux GOARCH=amd64 go build -o release/$(BINARY_NAME) .
	@tar -czf release/linux-amd64.tgz -C release $(BINARY_NAME)
	@rm release/$(BINARY_NAME)
	@# Linux i386
	GOOS=linux GOARCH=386 go build -o release/$(BINARY_NAME) .
	@tar -czf release/linux-i386.tgz -C release $(BINARY_NAME)
	@rm release/$(BINARY_NAME)
	@# Linux arm64
	GOOS=linux GOARCH=arm64 go build -o release/$(BINARY_NAME) .
	@tar -czf release/linux-arm64.tgz -C release $(BINARY_NAME)
	@rm release/$(BINARY_NAME)
	@# Darwin amd64
	GOOS=darwin GOARCH=amd64 go build -o release/$(BINARY_NAME) .
	@tar -czf release/darwin-amd64.tgz -C release $(BINARY_NAME)
	@rm release/$(BINARY_NAME)
	@# Darwin arm64
	GOOS=darwin GOARCH=arm64 go build -o release/$(BINARY_NAME) .
	@tar -czf release/darwin-arm64.tgz -C release $(BINARY_NAME)
	@rm release/$(BINARY_NAME)
	@# Windows amd64
	GOOS=windows GOARCH=amd64 go build -o release/$(BINARY_NAME).exe .
	@zip -jq release/windows-amd64.zip release/$(BINARY_NAME).exe
	@rm release/$(BINARY_NAME).exe
	@# Windows i386
	GOOS=windows GOARCH=386 go build -o release/$(BINARY_NAME).exe .
	@zip -jq release/windows-i386.zip release/$(BINARY_NAME).exe
	@rm release/$(BINARY_NAME).exe
	@echo "Release archives created in release/"
	@ls -la release/

.PHONY: clean
clean: ## Remove build artifacts
	@rm -f $(BINARY_NAME)
	@rm -rf dist/
	@rm -rf release/

## —— Dependencies ————————————————————————————————————————————————————————————
.PHONY: deps
deps: ## Download and tidy dependencies
	@go mod download
	@go mod tidy

## —— Help ————————————————————————————————————————————————————————————————————
.PHONY: help
help: ## Show available targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-12s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help