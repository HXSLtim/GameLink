# Agent Guide for GameLink Backend

Scope: This file applies to the entire directory tree rooted at `backend/`.

Use this guide to align code style, structure, and commands when making changes with an agent. When in doubt, defer to docs/go-coding-standards.md and existing code patterns.

## Communication

- 语言要求：所有与本仓库相关的代理回复一律使用中文（简体），除非用户在对话中明确要求使用其他语言。

## Project Structure

- Entry points: `cmd/<service>/main.go` (compose dependencies only; no business logic)
- Core modules under `internal/`:
  - `config` – configuration loading and defaults
  - `handler` – HTTP routes, params binding, uniform API responses
  - `service` – business rules, validation, caching, orchestration
  - `repository` – interfaces and pagination utils; `repository/gormrepo` holds GORM impls
  - `model` – entities, enums, DTOs
- Docs live in `docs/`. Style rules: `docs/go-coding-standards.md`.

If you add new top‑level folders, describe them in `docs/project-structure.md` as per org guidelines.

## Coding Standards

- Follow `docs/go-coding-standards.md` strictly.
- Formatting/imports: `gofmt` + `goimports` with local prefix `gamelink`.
- Naming: short lowercase package names; exported identifiers use UpperCamelCase.
- Errors: fail fast, wrap with `%w`, use `errors.Is`; use `service.ErrValidation` and `repository.ErrNotFound` consistently.
- Context: pass `ctx` to external I/O (`db.WithContext(ctx)`, cache, etc.).
- HTTP layer: RESTful, snake_case JSON, unified envelope `{success, code, message, data}`; admin routes must have auth + rate limit.
- Service: validation, cache invalidation; no HTTP concerns.
- Repository: normalize pagination; check `RowsAffected`; return `repository.ErrNotFound` when missing.

## Commands

- Dependencies: `make deps`
- Lint: `make lint` (uses `.golangci.yml`); ensure `golangci-lint` installed
- Tests: `make test`
- Run: `make run CMD=user-service`
- Build: `make build`

## Agent Workflow Expectations

- Prefer minimal, surgical patches using `apply_patch`.
- Keep changes focused on the request; do not refactor unrelated code.
- Maintain import grouping: stdlib, third‑party, then `gamelink`.
- Update docs when adding behaviors, configs, or directories.
- Do not introduce secrets or commit `.env` files.

## Review Checklist

- Code compiles locally (`go build ./...`) and tests pass (`go test ./...`).
- `golangci-lint run` is clean (or justified in PR).
- Public APIs preserve the response envelope and versioned paths.
- New/changed configs are documented in `docs/`.
