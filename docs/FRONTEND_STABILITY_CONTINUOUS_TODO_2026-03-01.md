# Frontend Stability Continuous Todo (2026-03-01)

Owner: Frontend Engineering  
Scope: `admin/` + `app/`  
Rule: only move deprecated files to archive directory, no hard delete.

## Status Legend
- `TODO`: not started
- `IN_PROGRESS`: currently being worked on
- `DONE`: completed and locally verified
- `BLOCKED`: waiting for external dependency/decision

## P0 (Immediate Stability)

| ID | Task | Acceptance Criteria | Status |
|---|---|---|---|
| P0-01 | Disable new PWA registration by default (`admin`) | PWA plugin is off unless `VITE_ENABLE_PWA=true` | DONE |
| P0-02 | One-time cleanup for legacy SW/cache (`admin`) | Legacy workbox/service worker cache is cleaned without wiping all browser caches | DONE |
| P0-03 | Frontend env validation gate (`admin` + `app`) | `predev/prebuild/prepreview` fail fast on invalid env config | DONE |

## P1 (Theme and UX Consistency)

| ID | Task | Acceptance Criteria | Status |
|---|---|---|---|
| P1-01 | Single source of theme tokens | `ThemeContext` and theme constants share one token source | DONE |
| P1-02 | Remove hard-coded conflicting primary colors in shared theme layer | `admin/src/theme/colors.ts` follows Ant CSS vars (no green/blue conflict) | DONE |
| P1-03 | Eliminate page-level hard-coded colors (`#1890ff` etc.) | All high-traffic pages use `theme.useToken()` / semantic tokens | DONE |

## P1 (User-End Hardening)

| ID | Task | Acceptance Criteria | Status |
|---|---|---|---|
| P1-04 | Add user app env template | `app/.env.example` exists and is documented | DONE |
| P1-05 | Add user app smoke tests for login/order/payment route flow | Core route render tests pass in CI | IN_PROGRESS |

## P2 (Quality & Monitoring)

| ID | Task | Acceptance Criteria | Status |
|---|---|---|---|
| P2-01 | Add Sentry (or equivalent) runtime error reporting | Production frontend errors are reportable with release/version tags | IN_PROGRESS |
| P2-02 | Add Playwright E2E for login/order/payment | Core E2E flow added and runs in CI | TODO |
| P2-03 | Deployment verification for frontend images | Docker deploy validation catches blank-page regressions before release | TODO |

## Progress Log
- 2026-03-01: Completed P0-01/P0-02/P0-03 and P1-01/P1-02.
- 2026-03-01: Added `scripts/validate-frontend-env.mjs`, wired npm pre-scripts in `admin` and `app`.
- 2026-03-01: Added `app/.env.example`, introduced `VITE_ENABLE_PWA=false` in `admin/.env.example`.
- 2026-03-01: Started P1-05 by adding `app/src/router.smoke.test.tsx` for login/order/payment-result route smoke coverage (execution still blocked by local Rollup optional dependency issue).
- 2026-03-01: Unblocked local Rollup optional dependency, app route smoke tests now executable and passing locally (`3 passed`).
- 2026-03-01: Started P2-01 with runtime monitoring scaffold in app (`setupGlobalErrorMonitoring`, optional Sentry global SDK bridge, env validation hooks).
- 2026-03-01: Completed P1-03 page-level primary-color cleanup in admin high-traffic pages/components; business code no longer contains `#1890ff/#7ACC35` hard-coded literals.
