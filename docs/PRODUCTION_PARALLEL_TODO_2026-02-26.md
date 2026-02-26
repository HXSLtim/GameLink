# Production Readiness Todo (Agent Parallel Plan)

Last Updated: 2026-02-26
Owner: Agent Team

## Status Legend
- `TODO`: not started
- `IN_PROGRESS`: currently being worked on
- `DONE`: completed and verified
- `BLOCKED`: cannot proceed without external dependency/decision

## P0 (Must Fix Before Production)

| ID | Module | Task | Acceptance Criteria | Status | Owner |
|---|---|---|---|---|---|
| P0-01 | API Routing | Register public payment callback routes (`/api/v1/public/payments/wechat/notify`, `/api/v1/public/payments/alipay/notify`) | Routes are mounted and reachable; callback path test added/passed | DONE | Agent-A |
| P0-02 | Payment Safety | Remove/disable production fallback to mock providers when provider init fails | In production mode, provider init failure causes explicit startup/runtime failure; no silent mock fallback | DONE | Agent-B |
| P0-03 | Payment Flow | Replace mock pay-info generation path for production mode | In production mode, real provider order creation is used; mock path only in non-prod | DONE | Agent-G |
| P0-04 | SMS | Remove production mock SMS sending behavior | In production mode, SMS provider performs real signed request or fails explicitly (no fake success) | DONE | Agent-C |
| P0-05 | Identity/Callback Verification | Ensure third-party verification does not silently use mock in production | Production mode cannot use mock verifier; missing config returns clear error | DONE | Agent-C |

## P1 (Should Fix)

| ID | Module | Task | Acceptance Criteria | Status | Owner |
|---|---|---|---|---|---|
| P1-01 | Admin API | Implement sensitive word detail endpoint (`GetSensitiveWord`) | `GET /admin/sensitive-words/:id` returns detail with 200/404 | DONE | Agent-D |
| P1-02 | Admin API | Enable batch player verification/revocation routes or remove dead docs | Route behavior aligns with docs and permission matrix | DONE | Agent-H |
| P1-03 | Admin API | Add batch payment operations (or explicitly de-scope + doc update) | Batch payment action path is implemented or documented as intentionally unsupported | DONE | Agent-H |
| P1-04 | ServiceItem | Add order reference check before batch delete | Deleting referenced items is rejected with clear error | DONE | Agent-I |
| P1-05 | SLA Alerting | Implement 2h dispute escalation supervisor alert | SLA >2h creates actionable alert record/notification | DONE | Agent-J |

## P2 (Nice to Have / Product Hardening)

| ID | Module | Task | Acceptance Criteria | Status | Owner |
|---|---|---|---|---|---|
| P2-01 | Finance | Reconciliation workflow (service/handler/router) | End-to-end reconciliation APIs + minimal operator flow | DONE | Agent-F |
| P2-02 | Metrics | Align metrics path and permission sync exclusions | `/metrics` exposure and permission sync list are consistent | DONE | Agent-F |

## Execution Notes
- Enforce TDD where applicable: failing test first, then implementation, then refactor.
- Every completed item must include:
  - changed files
  - verification commands
  - result (`PASS`/`FAIL`)
- Blocked items must include exact external dependency needed (credentials, provider contract, policy decision).

## Progress Log
- 2026-02-26: Initial todo created; awaiting first parallel execution wave.
- 2026-02-26: Wave-1 parallel completed for `P0-01`, `P0-02`, `P0-04`, `P0-05`, `P1-01`; package verification passed:
  - `go test ./internal/router ./internal/service/payment ./internal/service/external ./internal/service/sms ./internal/handler/admin -count=1`
- 2026-02-26: `P0-03` remains in progress (production flow has fail-closed guard, but full real provider order-create wiring is still pending).
- 2026-02-26: Wave-2 started in parallel for `P0-03`, `P1-02`, `P1-03`, `P1-04`, `P1-05`.
- 2026-02-26: Wave-2 completed; verification passed:
  - `go test ./internal/service/payment ./internal/handler/admin ./internal/service/item ./internal/repository/serviceitem ./pkg/scheduler -count=1`
  - `go test ./internal/router ./internal/service/external ./internal/service/sms -count=1`
- 2026-02-26: `P2-02` completed: permission sync skip-path updated to `/metrics` to match runtime route; stale legacy metrics-path references updated in related docs/examples. Verification passed:
  - `docker run --rm -v "<repo>/api:/src" -w /src golang:1.25.5-alpine sh -lc "/usr/local/go/bin/go test ./internal/router -run TestSyncAPIPermissionsSkipList_UsesRootMetricsPath -count=1"`
  - `docker run --rm -v "<repo>/api:/src" -w /src golang:1.25.5-alpine sh -lc "/usr/local/go/bin/go test ./internal/router -count=1"`
- 2026-02-26: `P2-02` ownership closure confirmed (`DONE`): metrics-path alignment remains `/metrics` across permission sync and docs/examples. Verification passed:
  - `docker run --rm -v "<repo>/api:/src" -w /src golang:1.25.5-alpine sh -lc "/usr/local/go/bin/go test ./internal/router -run TestSyncAPIPermissionsSkipList_UsesRootMetricsPath -count=1"`
- 2026-02-26: Wave-3 completed for `P2-01` reconciliation workflow: added admin reconciliation list/detail/create/execute APIs with repository/service/handler/router wiring. Verification passed:
  - `go test ./internal/repository/reconciliation ./internal/service/reconciliation ./internal/handler/admin ./internal/router -count=1`
- 2026-02-26: Post-P2 engineering hardening: fixed `scripts/pre-deployment-check.sh` early-exit bug under `set -e`, added `.env.production`/`admin/.env.production` example fallback for repository precheck runs, and made crypto signature/key checks conditional on encryption enabled. Verification passed:
  - `bash scripts/pre-deployment-check.sh`
