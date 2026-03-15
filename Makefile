.DEFAULT_GOAL := help

# ─── Help ────────────────────────────────────────────────────────────────────

.PHONY: help
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

# ─── Dependencies ────────────────────────────────────────────────────────────

.PHONY: deps
deps:
	go mod tidy
	go mod download

# ─── Test ────────────────────────────────────────────────────────────────────

HEADLESS ?= true
BROWSER  ?= chrome
TEST_DIR  = ./tests/...
GO_TEST   = go test $(TEST_DIR) -v -count=1

.PHONY: test
test: ## Run all tests (Chrome, headless)
	BROWSER=$(BROWSER) HEADLESS=$(HEADLESS) $(GO_TEST)

.PHONY: test-chrome
test-chrome: ## Run tests on Chrome
	BROWSER=chrome HEADLESS=true $(GO_TEST)

.PHONY: test-firefox
test-firefox: ## Run tests on Firefox
	BROWSER=firefox HEADLESS=true $(GO_TEST)

.PHONY: test-webkit
test-webkit: ## Run tests on WebKit
	BROWSER=webkit HEADLESS=true $(GO_TEST)

.PHONY: test-headed
test-headed: ## Run tests with browser visible
	BROWSER=chrome HEADLESS=false $(GO_TEST)

.PHONY: test-one
test-one: ## Run a specific test: make test-one TEST=TestName
	@test -n "$(TEST)" || (echo "Usage: make test-one TEST=TestName" && exit 1)
	BROWSER=chrome HEADLESS=true $(GO_TEST) -run $(TEST)# Generate and open Allure report

# ─── Allure ──────────────────────────────────────────────────────────────────

.PHONY: allure-serve
allure-serve: ## Serve Allure report
	allure serve ./allure-results

.PHONY: allure-generate
allure-generate: ## Generate Allure report
	allure generate ./allure-results -o ./allure-report --clean

# ─── Clean ───────────────────────────────────────────────────────────────────

.PHONY: clean
clean: ## Remove generated artifacts
	rm -rf ./allure-results ./allure-report ./screenshots
	mkdir -p ./allure-results ./screenshots

# ─── Format ──────────────────────────────────────────────────────────────────

.PHONY: fmt
fmt: fmt-go fmt-nix fmt-yaml ## Run all formatters

.PHONY: fmt-go
fmt-go: ## Format Go files
	goimports -w .

.PHONY: fmt-nix
fmt-nix: ## Format Nix files
	nixpkgs-fmt flake.nix

.PHONY: fmt-yaml
fmt-yaml: ## Format YAML files
	yamlfmt .lint: lint-go lint-yaml lint-nix

# ─── Lint ────────────────────────────────────────────────────────────────────

.PHONY: lint
lint: lint-go lint-nix lint-yaml ## Run all linters

.PHONY: lint-go
lint-go: ## Lint Go files
	golangci-lint run --config .golangci.yml ./...

.PHONY: lint-nix
lint-nix: ## Lint Nix files
	nixpkgs-fmt --check flake.nix
	statix check .
	deadnix --fail .

.PHONY: lint-yaml
lint-yaml: ## Lint YAML files
	yamllint -c .yamllint.yml .

# ─── CI ──────────────────────────────────────────────────────────────────────

.PHONY: ci
ci: lint test ## Run lint + test (for CI/CD)
