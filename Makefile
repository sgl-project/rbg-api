# rbg-api Makefile

# Module information
MODULE ?= sigs.k8s.io/rbgs/api

# Go environment
GO ?= go
GOFLAGS ?=
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

# Tools
CONTROLLER_GEN ?= $(shell which controller-gen 2>/dev/null || echo "$(GOPATH)/bin/controller-gen")
GOLANGCI_LINT  ?= $(shell which golangci-lint 2>/dev/null || echo "$(GOPATH)/bin/golangci-lint")

# CRD output path
CRD_OUTPUT_DIR ?= config/crd/bases

##@ General

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: generate
generate: ## Run all code generation (deepcopy + manifests).
	$(MAKE) generate-deepcopy
	$(MAKE) generate-manifests

.PHONY: generate-deepcopy
generate-deepcopy: ## Generate DeepCopy methods for API types.
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" \
		paths="./apis/..."

.PHONY: generate-manifests
generate-manifests: ## Generate CRD manifests.
	$(CONTROLLER_GEN) crd \
		paths="./apis/..." \
		output:crd:artifacts:config=$(CRD_OUTPUT_DIR)

.PHONY: generate-clients
generate-clients: ## Regenerate client-go, informers, and listers using k8s code-generator.
	bash hack/update-codegen.sh

.PHONY: fmt
fmt: ## Run go fmt against code.
	$(GO) fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	$(GO) vet ./...

.PHONY: lint
lint: ## Run golangci-lint.
	$(GOLANGCI_LINT) run ./...

.PHONY: tidy
tidy: ## Run go mod tidy.
	$(GO) mod tidy

.PHONY: vendor
vendor: ## Update vendor directory.
	$(GO) mod vendor

##@ Build

.PHONY: build
build: ## Build all packages (compile check).
	$(GO) build ./...

##@ Test

.PHONY: test
test: ## Run unit tests.
	$(GO) test ./... -v -count=1

##@ Tools

.PHONY: install-tools
install-tools: ## Install required code-generation tools.
	$(GO) install sigs.k8s.io/controller-tools/cmd/controller-gen@latest
	$(GO) install k8s.io/code-generator/cmd/client-gen@latest
	$(GO) install k8s.io/code-generator/cmd/lister-gen@latest
	$(GO) install k8s.io/code-generator/cmd/informer-gen@latest
	$(GO) install k8s.io/code-generator/cmd/deepcopy-gen@latest
