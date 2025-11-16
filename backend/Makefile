GOCMD ?= go
LINT ?= golangci-lint

.PHONY: deps lint test run build swagger

deps:
	$(GOCMD) mod tidy

lint:
	$(LINT) run --timeout=5m

test:
	$(GOCMD) test ./...

run:
ifndef CMD
	$(error usage: make run CMD=<service-name>)
endif
	$(GOCMD) run ./cmd/$(CMD)

build:
	$(GOCMD) build ./cmd/...

swagger:
	@which swag >/dev/null 2>&1 || (echo "Installing swag with generics support" && GOFLAGS= go install github.com/swaggo/swag/cmd/swag@latest)
	@echo "Generating Swagger docs from annotations (with generics support)..."
	$(GOCMD) env GOPATH >/dev/null 2>&1
	@# Ensure we have a recent version of swag that supports generics
	swag init -g cmd/main.go -o docs/swagger -q --parseDependency --parseInternal
	@# Keep embedded swagger json in sync for /swagger.json endpoint
	cp docs/swagger/swagger.json internal/handler/swagger/openapi.json
	@echo "Swagger docs generated successfully"

swagger-check:
	@echo "Checking swag version for generics support..."
	@which swag >/dev/null 2>&1 || (echo "swag not installed" && exit 1)
	@swag --version | grep -E "v[2-9]\." >/dev/null || (echo "Please update swag to v2.0+ for generics support" && exit 1)

fix-swagger:
	@echo "Fixing Swagger annotations..."
	@which python3 >/dev/null 2>&1 && python3 fix_swagger_annotations.py internal/handler || echo "Python3 not available, manual fixing required"
