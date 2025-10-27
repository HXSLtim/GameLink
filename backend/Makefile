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
	@which swag >/dev/null 2>&1 || (echo "Installing swag" && GOFLAGS= go install github.com/swaggo/swag/cmd/swag@latest)
	@echo "Generating Swagger docs from annotations..."
	$(GOCMD) env GOPATH >/dev/null 2>&1
	swag init -g cmd/user-service/main.go -o docs/swagger -q
