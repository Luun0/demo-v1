# Strict targets: every step must exist and pass. Nothing is skipped.
SHELL := /bin/bash
.SHELLFLAGS := -eu -o pipefail -c
GOBIN ?= $(shell go env GOPATH)/bin
export PATH := $(GOBIN):$(PATH)
NPM := npm --prefix frontend

.PHONY: help tools install check lint typecheck test build \
        backend-lint backend-typecheck backend-test backend-build \
        frontend-lint frontend-typecheck frontend-test frontend-build contract-lint \
        secrets size-check run dev seed e2e vuln

help:
	@echo "make tools     install golangci-lint, goimports, gitleaks"
	@echo "make install   npm ci for the frontend"
	@echo "make check     lint + typecheck + test + build (backend and frontend)"
	@echo "make run       start the frontend dev server"

tools:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.15.0
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/zricethezav/gitleaks/v8@v8.30.1

install:
	$(NPM) ci

check: lint typecheck test build
	@echo "make check: OK"

lint: backend-lint frontend-lint contract-lint size-check
typecheck: backend-typecheck frontend-typecheck
test: backend-test frontend-test
build: backend-build frontend-build

backend-lint:
	@out=$$(cd backend && gofmt -l .); if [ -n "$$out" ]; then echo "gofmt needed:"; echo "$$out"; exit 1; fi
	cd backend && golangci-lint run ./...

backend-typecheck:
	cd backend && go vet ./...

backend-test:
	cd backend && ./scripts/covergate.sh

backend-build:
	cd backend && go build ./...

frontend-lint:
	$(NPM) run lint

contract-lint:
	$(NPM) run lint:contract
	$(NPM) run check:api

frontend-typecheck:
	$(NPM) run typecheck

frontend-test:
	$(NPM) test

frontend-build:
	$(NPM) run build
	$(NPM) run size:bundle

size-check:
	./scripts/check-file-size.sh

secrets:
	gitleaks git --no-banner --redact .

run:
	$(NPM) run dev

dev: ## full stack in Docker (db, migrate, api, web)
	docker compose up --build -d
	@echo "web: http://localhost:$${WEB_PORT:-8081}  api: http://localhost:$${APP_PORT:-8080}"

seed: ## synthetic demo data (PRV-08) into the compose db
	docker compose run --rm api seed

e2e: ## Playwright + axe against a running API on :8080 (see docs/verification-plan.md)
	$(NPM) run test:e2e

vuln: ## dependency vulnerability scans
	cd backend && go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	$(NPM) audit --audit-level=high
