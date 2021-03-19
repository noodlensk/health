dep: ## Install deps
	go mod tidy
lint: ## Run linter
	golangci-lint run

lint-fix: ## Run linter with fix command
	golangci-lint run --fix

fmt: ## gofmt and goimports all go files
	find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do gofmt -w -s "$$file"; goimports -w "$$file"; done

build-server: ## Build server
	go build ./cmd/server

release:
	mkdir -p release && \
	cd ./cmd/server && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./../../release/server_linux && \
	cd ../../ && \
	zip --junk-paths  health ./release/server_linux &&\
	rm -rf release

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
