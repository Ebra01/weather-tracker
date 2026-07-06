SQLC ?= sqlc
GO ?= go

.PHONY: help generate-db test verify check-sqlc

help:
	@echo "Available commands:"
	@echo "  make generate-db  Generate sqlc database code"
	@echo "  make test         Run Go tests"
	@echo "  make verify       Generate database code, then run tests"

check-sqlc:
	@command -v $(SQLC) >/dev/null 2>&1 || { echo "sqlc not found. Install sqlc or run with SQLC=/path/to/sqlc."; exit 1; }

generate-db: check-sqlc
	$(SQLC) generate

test:
	$(GO) test ./...

verify: generate-db test
