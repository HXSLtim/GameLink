# GitHub Actions Workflows

This directory contains CI/CD workflows for the GameLink project.

## Workflow Overview

### Core CI/CD Workflows

| Workflow | Purpose | Trigger | Status |
|----------|---------|---------|--------|
| **ci.yml** | Continuous Integration (tests, linting, builds) | Push/PR to main/dev | ✅ Active |
| **deploy.yml** | Deployment to Staging/Production | Tag push / Manual | ✅ Active |
| **security.yml** | Security vulnerability scanning | Push/PR/Schedule | ✅ Active |
| **pre-commit-check.yml** | Pre-commit code quality gates | Pre-commit hook | ✅ Active |

### New Optimization Workflows

| Workflow | Purpose | Trigger | Priority |
|----------|---------|---------|----------|
| **performance.yml** | Performance benchmarking and regression testing | PR/Schedule/Manual | P0 |
| **e2e.yml** | End-to-end integration testing | Schedule/Manual | P0 |
| **flow-guard-regression.yml** | Order/Payment guard regression (unpaid-complete block) | Schedule/Manual | P0 |
| **withdraw-flow-regression.yml** | Withdraw flow regression (request/approve/complete) | Schedule/Manual | P0 |
| **test-report.yml** | Test coverage and results aggregation | PR/Manual | P1 |
| **dependabot-merge.yml** | Automated dependency updates | Dependabot PR | P2 |

## Workflow Details

### CI Workflow (`ci.yml`)

**Purpose:** Main continuous integration pipeline

**Features:**
- Change detection to skip unnecessary builds
- Parallel testing for backend, admin, and client
- Code coverage enforcement (70% threshold)
- Linter checks (golangci-lint, ESLint)
- Type checking (TypeScript)
- Docker image building with caching

**Quality Gates:**
- Backend: 70% test coverage required
- Admin: Type check + linter + tests
- Client: Type check + linter + tests

### Performance Workflow (`performance.yml`)

**Purpose:** Detect performance regressions

**Features:**
- Go benchmark tests with `-bench` and `-benchmem`
- Baseline comparison using `benchstat`
- 10% performance degradation threshold
- API response time load testing with `hey`
- Automated PR commenting with results

**Schedule:** Daily at midnight UTC + on PRs

**Configuration:**
- Baseline stored in GitHub Actions artifacts
- Updated automatically on main branch
- 90-day retention for baselines

### E2E Workflow (`e2e.yml`)

**Purpose:** End-to-end testing of critical user flows

**Features:**
- Backend E2E with real database
- User flow testing with Playwright
- Integration scenario testing
- Full-stack testing with Docker Compose

**Test Scenarios:**
- User registration and login
- Order creation flow
- Payment integration
- Real-time chat (WebSocket)

**Schedule:** Daily at 1 AM UTC

### Flow Guard Regression Workflow (`flow-guard-regression.yml`)

**Purpose:** Verify critical business guards for order completion and reviews

**Features:**
- Starts backend with PostgreSQL + Redis service containers
- Enables seed data and runs full flow-guard regression script
- Verifies unpaid order completion is blocked (`400`)
- Verifies review submission before completion is blocked (`400`)
- Verifies paid flow can complete and review successfully

**Schedule:** Daily at 2:30 AM UTC + manual trigger

### Withdraw Flow Regression Workflow (`withdraw-flow-regression.yml`)

**Purpose:** Verify withdraw end-to-end flow with automatic balance precheck

**Features:**
- Starts backend with PostgreSQL + Redis service containers
- Enables seed data and runs withdraw flow regression script
- Auto-clears pending/approved withdraws for target player (precheck)
- Auto-topups available balance when needed (paid + completed order path)
- Verifies request (`/player/earnings/withdraw`) to admin approve/complete
- Verifies player withdraw-history and summary delta assertions
- Uploads regression script output (`withdraw-flow-regression.log` + summary markdown artifact)

**Schedule:** Daily at 3:00 AM UTC + manual trigger

### Test Report Workflow (`test-report.yml`)

**Purpose:** Aggregate and display test results

**Features:**
- Merges coverage reports from all components
- Generates security scan summary
- Creates performance test reports
- Posts summary as PR comment

**Report Sections:**
- Backend (Go) coverage
- Admin Panel (TypeScript) coverage
- Client (TypeScript) coverage
- Security scan results
- Performance benchmarks

### Deploy Workflow (`deploy.yml`)

**Purpose:** Automated deployment to Staging/Production

**Features:**
- Pre-deployment validation
- Docker image building and pushing to GHCR
- Health checks with retry logic
- Deployment notifications (DingTalk)
- Automated rollback on failure
- GitHub deployment status updates

**Environments:**
- **Staging:** Automatic on `-rc` tags or manual trigger
- **Production:** Manual trigger only (requires approval)

**Health Checks:**
- Backend: `/api/v1/healthz` endpoint
- Frontend: Homepage accessibility
- 5 attempts with 10-second intervals

### Security Workflow (`security.yml`)

**Purpose:** Scan for security vulnerabilities

**Scanners:**
- **govulncheck:** Go dependency vulnerabilities
- **npm audit:** Node.js dependency vulnerabilities
- **Trivy:** Docker image vulnerabilities
- **Gitleaks:** Secret/key leakage detection

**Schedule:** Runs on every push, PR, and daily

### Pre-commit Check (`pre-commit-check.yml`)

**Purpose:** Fast pre-commit validation

**Checks:**
- File size limits (max 2MB)
- Credential scanning
- Basic syntax validation

**Benefits:** Catches issues before commit

### Dependabot Auto-Merge (`dependabot-merge.yml`)

**Purpose:** Automatically merge non-breaking dependency updates

**Conditions:**
- PR created by Dependabot
- Labeled with `dependencies`
- All CI checks passed
- No merge conflicts

**Merge Strategy:** Squash and merge with auto-deletion

## Environment Variables

### Required Secrets

| Secret | Description | Used By |
|--------|-------------|---------|
| `GITHUB_TOKEN` | GitHub token (auto-provided) | All workflows |
| `DINGTALK_WEBHOOK` | DingTalk notification URL | deploy.yml |
| `CODECOV_TOKEN` | Codecov upload token | ci.yml (optional) |

### Environment Configuration

Workflows use the following environment variables:

```yaml
NODE_VERSION: '20'
GO_VERSION: '1.25.5'
REGISTRY: ghcr.io
```

## Usage

### Triggering Workflows

**Automatic:**
- Push to main/dev branches → CI + Security
- Create PR → CI + Security + Test Report
- Create tag `v*` → Deploy
- Daily schedule → Performance + E2E + Flow Guard + Withdraw Flow + Security

**Manual:**
```bash
# Deploy to staging
gh workflow run deploy.yml -f environment=staging

# Deploy to production
gh workflow run deploy.yml -f environment=production

# Run performance tests
gh workflow run performance.yml

# Run E2E tests
gh workflow run e2e.yml -f environment=staging

# Run flow-guard regression
gh workflow run flow-guard-regression.yml

# Run withdraw-flow regression
gh workflow run withdraw-flow-regression.yml
```

### Checking Workflow Status

```bash
# List recent workflow runs
gh run list --workflow=ci.yml

# View specific run details
gh run view <run-id>

# Watch logs in real-time
gh run watch

# Download artifacts
gh run download <run-id>
```

## Troubleshooting

### Common Issues

**1. CI timeouts**
- Increase timeout in workflow if needed
- Check for slow tests or external dependencies
- Use `--short` flag for unit tests in CI

**2. Docker build cache miss**
- Cache is invalidated if Dockerfile changes
- First build after cache miss will be slower
- Subsequent builds will use cache

**3. Health check failures**
- Services may take longer to start than expected
- Increase `max_attempts` or `wait_seconds` in health check
- Check application logs for startup errors

**4. Permission errors**
- Ensure `GITHUB_TOKEN` has required permissions
- Add `permissions:` block to workflow if needed

### Debugging Failed Runs

```bash
# Re-run failed workflow
gh run rerun <run-id>

# Re-run with debug logging
gh run rerun <run-id> --debug

# Cancel running workflow
gh run cancel <run-id>
```

## Performance Optimization

### Current Optimization

- **Docker layer caching:** Uses GitHub Actions cache
- **Go module caching:** Speeds up dependency download
- **npm caching:** Caches node_modules
- **Change detection:** Skips builds for unchanged components

### Future Optimizations

- [ ] Self-hosted runners for faster builds
- [ ] Parallel test execution within packages
- [ ] Incremental build caching
- [ ] Test result caching for flaky tests

## Monitoring

### Key Metrics

- **Build duration:** Target < 10 minutes
- **Test pass rate:** Target > 95%
- **Code coverage:** Target > 70%
- **Deployment success rate:** Target > 95%

### Viewing Metrics

- GitHub Actions dashboard: https://github.com/your-org/GameLink/actions
- Workflow run history and trends
- Individual job timing and resource usage

## Contributing

When adding new workflows:

1. Follow existing naming conventions (`kebab-case.yml`)
2. Add documentation to this README
3. Use appropriate triggers (automatic vs manual)
4. Include error handling and notifications
5. Test in fork before merging to main

## Related Documentation

- **CI/CD Optimization Plan:** `docs/CICD_OPTIMIZATION_PLAN.md`
- **Deployment Checklist:** `docs/DEPLOYMENT_CHECKLIST.md`
- **Security Hardening:** `docs/SECURITY_HARDENING.md`
- **Monitoring Guide:** `docs/MONITORING_ALERT_GUIDE.md`

## Support

For issues or questions about CI/CD:

1. Check this README first
2. Review workflow logs for error details
3. Consult optimization plan documentation
4. Contact DevOps-Engineer team member
