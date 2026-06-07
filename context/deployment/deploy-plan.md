---
project: weekplate
deployed_at: 2026-06-07
platform: Railway
environment: production
---

## First Deployment — Skeleton Server

### What was deployed

Minimal Go HTTP server (no feature code yet) to validate the deployment pipeline end-to-end before feature development begins.

- **Commit:** `e7e35f2` — "chore: Railway deployment scaffold and minimal HTTP server skeleton"
- **Branch:** `master`
- **Build:** Nixpacks v1.41.0 detected Go, ran `go build -o app ./cmd/server` from `railway.toml`
- **Health check:** `/health` passed on first deploy

### Railway project

| Field | Value |
|---|---|
| Project name | weekplate |
| Project ID | e27fcd42-51e4-4136-8d45-a2d75171ee31 |
| Environment | production |
| Region | sfo (San Francisco) |
| App service name | web |
| App service ID | 06fb9661-a3f5-4911-a1e8-e86e9124de8f |
| Public URL | https://web-production-cd9e9.up.railway.app |
| Database service | Postgres (PostgreSQL 18, SSL) |
| Database service ID | 7291d546-047c-4378-b8cd-61b60abd855c |

### Secrets wired

| Variable | Set by | Notes |
|---|---|---|
| `DATABASE_URL` | Railway (auto-injected from Postgres service reference) | Internal private network URL |
| `DATABASE_PUBLIC_URL` | Railway (auto-injected) | External proxy URL for local DB access |
| `APP_ENV` | Manual (`railway variable set`) | Value: `production` |
| `SESSION_SECRET` | Manual (`railway variable set`) | 64-char hex, generated via PowerShell CSPRNG |
| `PORT` | Railway (auto-injected) | Read by server at startup |

### Verification

```
GET https://web-production-cd9e9.up.railway.app/health  → 200 ok
GET https://web-production-cd9e9.up.railway.app/         → 200 weekplate — coming soon
```

Database connection confirmed live: `/health` pings Postgres on every request and returned 200.

### Operational notes

- **Auto-deploy:** Railway watches the `master` branch on `github.com/vschetko/weekplate`. Every push to `master` triggers a new Nixpacks build.
- **Rollback:** `railway service redeploy --service web` to re-run the previous deployment, or use the Railway dashboard deployment history.
- **Logs:** `railway logs` (runtime) or `railway logs --build <deployment-id>` (build output).
- **CLI install:** `npm install -g @railway/cli` (winget package does not exist; Node.js is the correct path on Windows).

### Known issues resolved during first deploy

1. **winget package doesn't exist** — `Railway.RailwayCLI` is not in the winget registry. Correct install: `npm install -g @railway/cli`. `infrastructure.md` Getting Started section updated.
2. **Branch mismatch** — Railway was initially connected to `main` (empty branch); redirected to `master` via `railway service source connect --branch master`.
3. **First build triggered before code landed** — Railway triggered a build immediately on repo connection (before our push). Resolved by `railway service redeploy` after source branch was corrected.

### Next step

Implement feature code in `internal/` following the PRD at `context/foundation/prd.md`. Every push to `master` deploys automatically.
