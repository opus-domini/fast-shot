GOCMD        ?= go
LINT          = golangci-lint
VULNCHECK     = govulncheck
LINT_GOCACHE ?= /tmp/go-cache
LINT_CACHE   ?= /tmp/golangci-lint-cache

.DEFAULT_GOAL := help

# ─── Quality ──────────────────────────────────────────────────

.PHONY: test
test: check-go ## Run tests
	$(GOCMD) test ./...

.PHONY: test-coverage
test-coverage: check-go ## Run tests with race detection and coverage
	$(GOCMD) test -race -covermode=atomic -coverprofile=coverage.txt ./...

.PHONY: benchmark
benchmark: check-go ## Run benchmarks
	$(GOCMD) test -run=^$$ -bench=. -benchmem ./...

.PHONY: fmt
fmt: check-go check-lint ## Format code
	GOCACHE=$(LINT_GOCACHE) GOLANGCI_LINT_CACHE=$(LINT_CACHE) $(LINT) fmt

.PHONY: lint
lint: check-go check-lint ## Lint code
	GOCACHE=$(LINT_GOCACHE) GOLANGCI_LINT_CACHE=$(LINT_CACHE) $(LINT) run

.PHONY: tidy
tidy: check-go ## Tidy go.mod and go.sum
	$(GOCMD) mod tidy

.PHONY: security
security: check-go check-vuln ## Run vulnerability scanner
	$(VULNCHECK) ./...

# ─── Build ────────────────────────────────────────────────────

.PHONY: build
build: check-go ## Build project examples
	$(GOCMD) build -o build/examples/ ./examples/...

# ─── CI ───────────────────────────────────────────────────────

.PHONY: ci-fast
ci-fast: fmt lint tidy test security ## Fast PR gate

.PHONY: ci-full
ci-full: ci-fast test-coverage benchmark ## Full mainline gate

.PHONY: ci
ci: ci-full ## Run full CI pipeline

# ─── Maintenance ──────────────────────────────────────────────

.PHONY: clean
clean: ## Remove build artifacts
	$(GOCMD) clean
	rm -rf build coverage.txt

.PHONY: help
help: ## Show available targets
	@awk '\
		/^# ─── / { printf "\n\033[1m%s\033[0m\n", substr($$0, 7) } \
		/^[a-zA-Z_-]+:.*## / { \
			target = $$0; \
			sub(/:.*/, "", target); \
			desc = $$0; \
			sub(/.*## /, "", desc); \
			printf "  \033[36m%-18s\033[0m %s\n", target, desc; \
		}' $(MAKEFILE_LIST)
	@echo

# Dependency checks (internal)

.PHONY: check-go check-lint check-vuln

check-go:
	@command -v $(GOCMD) >/dev/null 2>&1 || { echo "error: go is not installed — https://golang.org/doc/install"; exit 1; }

check-lint:
	@command -v $(LINT) >/dev/null 2>&1 || { echo "error: $(LINT) is not installed — go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest"; exit 1; }

check-vuln:
	@command -v $(VULNCHECK) >/dev/null 2>&1 || { echo "error: $(VULNCHECK) is not installed — go install golang.org/x/vuln/cmd/govulncheck@latest"; exit 1; }
