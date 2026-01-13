# Repository Guidelines

## Project Structure & Module Organization

This repo is a monorepo:

- `api/`: Go backend (Gin/GORM). Entry point: `api/cmd/main.go`. Main code lives in `api/internal/` (layered `handler/ → service/ → repository/ → model/`). Tests live alongside code (`*_test.go`) and under `api/tests/` for integration/load tooling.
- `admin/`: Admin web UI (React + TypeScript + Ant Design + Vite). Source in `admin/src/`, unit tests in `admin/src/test/`, and Playwright E2E tests under `admin/tests/`.
- `client/`: **Desktop PWA** for User/Player (React + TypeScript + shadcn/ui + Vite). Targets desktop browsers with PWA support for offline capability. Source in `client/src/`.
- `app/`: **Mobile** Taro mini‑program app for User/Player. Source in `app/src/` (Taro conventions).
- `docs/`: Architecture, standards, and feature docs. `scripts/`: deployment and Docker helper scripts (mostly PowerShell).

## Build, Test, and Development Commands

- Docker: `docker-compose up -d` (uses `.env`; start from `.env.example`).
- Backend: `cd api && make deps && go run cmd/main.go` (run API locally).
- Backend tests: `cd api && make test` or `make test-coverage` (writes `coverage.out`).
- Backend quality: `cd api && make check` (runs `fmt`, `vet`, `lint`, `test`).
- Admin: `cd admin && npm install && npm run dev` (Vite dev server).
- Admin tests: `cd admin && npm run test` (Vitest) and `npm run test:e2e` (Playwright).
- Client: `cd client && npm install && npm run dev`.
- Taro app: `cd app && pnpm install && pnpm run dev:weapp` (or other `dev:*` targets).

## Coding Style & Naming Conventions

- Go: `gofmt` + `goimports` (local prefix `gamelink`), and `golangci-lint` configured in `api/.golangci.yml`.
- TypeScript/React: ESLint is enabled in each frontend package; `admin/` also uses Prettier (`admin/.prettierrc`). Prefer existing patterns/components in the same module.
- Naming: use `camelCase` for JS/TS identifiers, `PascalCase` for React components, and standard Go exported naming. Keep filenames descriptive (e.g., `order_service.go`, `OrderDetailPage.tsx`).

## Testing Guidelines

- Go: add/keep tests close to code (`*_test.go`). Prefer table-driven tests and `testify` assertions. Use build tag `integration` when a test requires external services.
- Frontend: unit tests via Vitest; E2E tests via Playwright. Name unit tests `*.test.ts(x)` and keep E2E specs in `admin/tests/`.

## Commit & Pull Request Guidelines

- Commits follow Conventional Commits: `type(scope): subject` (examples: `feat(admin): ...`, `fix(ci): ...`, `test(integration): ...`).
- PRs should include: a clear description, linked issue/PRD when applicable, screenshots for UI changes, and the commands/results you ran (e.g., `make test-coverage`, `npm run test`).

## Communication

- 与用户交流时请使用中文（简体）。

## Security & Configuration Tips

- Never commit secrets. Use `.env.example` as the template and document new env vars in `README.md` or `docs/`.
