# These shell flags are REQUIRED for an early exit in case any program called by make errors!
.SHELLFLAGS=-euo pipefail -c
SHELL := /bin/bash

.PHONY: all fmt clean check build tidy goimports golangci-lint

# Set the GOBIN environment variable so that dependencies will be installed
# always in the same place, regardless of the value of GOPATH
CACHE := $(PWD)/.cache
export GOBIN := $(PWD)/bin
export PATH := $(GOBIN):$(PATH)

IMAGE := vaultie-talkie`
DEV_TAG := dev
all: build

clean: ## Clean this directory
	@rm -fr $(CACHE) $(GOBIN) bin/* dist/ || true


build: tidy ## Build binaries
	@CGO_ENABLED=0 go build -a -o $(GOBIN)/vaultie-talkie *.go

run: build
	@bin/vaultie-talkie

tidy:
	@go mod tidy

verify:
	@go mod verify

fmt:
	@go fmt ./...

check: golangci-lint goimports

vault-dev:
	@docker run -d --cap-add=IPC_LOCK -e 'VAULT_DEV_ROOT_TOKEN_ID=myroot' -e 'VAULT_DEV_LISTEN_ADDRESS=0.0.0.0:8200' -p 8200:8200 vault

GOIMPORTS := $(GOBIN)/goimports
goimports:
	@$(call go-get-tool,$(GOIMPORTS),golang.org/x/tools/cmd/goimports)
	@$(GOIMPORTS) -w -l $(shell find . -type f -name "*.go" -not -path "./vendor/*")

GOLANGCI_LINT := $(GOBIN)/golangci-lint
golangci-lint:
	@$(call go-get-tool,$(GOLANGCI_LINT),github.com/golangci/golangci-lint/cmd/golangci-lint@v1.43.0)
	@echo "Running golangci-lint..."
	@$(GOLANGCI_LINT) run --timeout=10m -E unused,gosimple,staticcheck --skip-dirs-use-default --verbose

install-pre-commit:
	@pip install pre-commit

pre-commit: install-pre-commit
	@pre-commit install --hook-type commit-msg

test:
	@go clean -testcache ./...
	@go test ./... -v

# go-get-tool will 'go get' any package $2 and install it to $1.
define go-get-tool
@[ -f $(1) ] || { \
set -e ;\
TMP_DIR=$$(mktemp -d) ;\
cd $$TMP_DIR ;\
go mod init tmp ;\
echo "Downloading $(2)" ;\
go get $(2) ;\
rm -rf $$TMP_DIR ;\
}
endef

help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
