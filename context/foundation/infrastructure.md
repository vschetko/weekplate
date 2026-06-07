---
project: weekplate
researched_at: 2026-06-07T00:00:00Z
recommended_platform: Railway
runner_up: Fly.io
context_type: mvp
tech_stack:
  language: go
  framework: standard library
  runtime: compiled binary (single Go binary via Docker)
---

## Recommendation

**Deploy on Railway.**

Railway is the strongest fit for weekplate's Go binary + PostgreSQL stack at MVP scale: it provides co-located Postgres with automatic `DATABASE_URL` injection, Nixpacks auto-builds Go without a Dockerfile, and has the most mature Claude Code / MCP integration of any candidate platform (remote hosted MCP at `mcp.railway.com`, dedicated Claude Code agent page). The Hobby plan covers both app and DB in a single pay-as-you-go account with a $5 monthly credit. The scoring edge over Fly.io comes from Railway's first-class agent integration and lower ops burden for a solo after-hours developer — the cross-check on Fly.io surfaced unmanaged Postgres maintenance as a real risk given the timeline.

**Note:** Fly.io was the original `deployment_target` in `tech-stack.md`. After research and anti-bias cross-check, Railway was selected as the stronger fit. If Railway's Hobby-tier reliability becomes a concern at launch, Fly.io remains the validated runner-up and migration is straightforward (both use Docker containers).

## Platform Comparison

| Platform | CLI-first | Managed/Serverless | Agent-readable docs | Stable deploy API | MCP / Integration | Total |
|---|---|---|---|---|---|---|
| **Railway** | Pass | Pass | Pass | Pass | Pass | **5 / 5** |
| **Fly.io** | Pass | Pass | Pass (llms.txt + MDX) | Pass | Partial (experimental) | **4.5 / 5** |
| **Render** | Partial | Pass | Partial | Partial | Fail | **2 / 5** |
| Cloudflare Workers | Pass | Pass | Pass | Pass | Pass | Dropped (Go WASM only; no Postgres) |
| Vercel | Pass | Pass | Pass | Pass | Partial | Dropped (no co-located DB) |
| Netlify | Pass | Pass | Partial | Pass | Pass | Dropped (no DB; serverless-only Go) |

**Hard filters applied:**
- Cloudflare Workers: Go requires WASM compilation — incompatible with Go standard library HTTP server binary. D1 is SQLite-only, no Postgres.
- Vercel: No co-located database; primarily frontend/edge. Q5 answer (co-location preferred) eliminates it.
- Netlify: Go is serverless-functions-only (60s timeout); no database co-location.

**Interview weights applied:**
- Q2 (no strong cost preference) → neutral between Railway Hobby (~$10-18/month real cost) and Fly.io (~$7-15/month)
- Q3 (familiar with Railway/Render/Fly.io) → ties broken within this group
- Q4 (single region) → no edge/CDN advantage needed; both Railway EU-West and Fly.io single-region work
- Q5 (co-location preferred) → Railway and Fly.io both pass; Render's Starter Postgres (256MB RAM) fails for production use

### Shortlisted Platforms

#### 1. Railway (Recommended)

Railway scores the highest on all five agent-friendly criteria. The CLI (`railway`) covers deploy, logs, variables, and domain management. Nixpacks eliminates the Dockerfile requirement for Go projects — `git push` triggers an auto-detected Go build. PostgreSQL is a first-class citizen: one-click provision, auto-injected `DATABASE_URL`, co-located in the same region as the app. The MCP story is the strongest in the pool: a local CLI-based MCP server, a hosted remote MCP at `mcp.railway.com` (OAuth-backed), and a dedicated Claude Code integration page. Hobby plan includes $5 monthly credit; real all-in cost for app + DB running 24/7 is $10–18/month. No SLA on Hobby is the primary trade-off.

#### 2. Fly.io

Fly.io was the original tech-stack choice and remains a strong second. `flyctl` is comprehensive (deploy, rollback, logs, SSH, Postgres management). Go binary deployment via Docker is well-established: lean single-binary images with fast startup and low memory overhead map directly to Fly.io's Machine model. Postgres is co-located (same region, private network). Documentation is agent-readable (llms.txt + MDX on GitHub). The MCP story is weaker — `fly mcp server` is experimental, and Fly.io's own team advised against MCP in 2026. The primary risk surfaced by cross-check: unmanaged Postgres (the default) requires manual disk monitoring, version upgrades, and backup configuration — ops burden for a solo after-hours developer.

#### 3. Render

Render offers native Go support, managed Postgres, and a genuinely simple UI. The free tier for static assets and the $7/month Starter web service are the cheapest continuous-deployment options in the shortlist. However, Render's CLI story is the weakest of the three (deploy webhooks rather than a first-class CLI tool), there is no MCP server, and the April 2026 workspace plan changes introduced per-user fees on the Professional plan. The Starter Postgres ($7/month, 256MB RAM, 1GB storage) is adequate for development but underpowered for production. Render is a reasonable choice for a developer who wants the simplest possible UI-driven setup and doesn't prioritize agent operability.

## Anti-Bias Cross-Check: Railway

### Devil's Advocate — Weaknesses

1. **Hobby plan $5 credit is not a free tier.** A Go app + Postgres running continuously will exhaust the credit in 5–10 days. Real all-in cost: $10–18/month. Railway's marketing implies "starts free"; the reality is "starts with $5 credit."
2. **No SLA on Hobby — reliability is best-effort.** Railway has experienced platform incidents in 2024–2025. For a meal-planning app peaking on weekend mornings, an unannounced Saturday incident means missing the only high-value usage window. No recourse without upgrading to Pro ($20/month).
3. **Nixpacks can produce non-obvious build failures.** Non-standard `go.mod` layouts, CGO dependencies, or custom build tags can cause silent failures or binaries that compile but can't connect to `DATABASE_URL` correctly. The fallback is adding a Dockerfile — but this defeats the main DX advantage.
4. **Only 4 regions.** US-West, US-East, EU-West (Amsterdam), Asia-Pacific. If target users are in Central Europe or other underserved regions, latency may be higher than expected.
5. **Railway CLI requires Node.js on Windows.** The `railway` CLI is an npm package. On the user's Windows environment, a standalone binary installer exists but is less prominently documented than `npm install -g @railway/cli`.

### Pre-Mortem — How This Could Fail

Six months after deploying weekplate on Railway, the project ran into a quiet resource ceiling. The solo dev had enjoyed frictionless onboarding: Nixpacks detected Go, Postgres provisioned in one click, `DATABASE_URL` was auto-injected. The Hobby plan's $5 credit covered the first week, then monthly bills of $12–18 started arriving — not painful, but not zero as implied.

The critical failure came at month 4: a Railway platform incident took the service down for 3 hours on a Saturday morning — peak meal-planning time. With no SLA on Hobby, there was no recourse. A handful of early adopters who had tried the app that weekend never returned. Upgrading to Pro ($20/month) would have provided a stronger reliability posture but was added too late.

A second issue emerged when the dev attempted to add recipe image uploads: Railway's volume (persistent disk) support was in limited beta on EU-West and behaved inconsistently. The workaround — external R2 object storage — worked, but it added an unexpected integration step mid-MVP.

The root assumption was that Railway's "it just works" positioning would hold under production conditions. At Hobby tier, the managed feel is genuine for the happy path; but platform reliability and feature availability are not at the same level as Pro or a more established platform.

### Unknown Unknowns

- **`DATABASE_URL` format requires correct parsing in Go.** Railway injects Postgres credentials as a single `postgresql://user:pass@host:port/db` URL. Go's `database/sql` DSN format is different. Libraries like `pgx` or `lib/pq` accept the URL format, but the app must be written to parse it — not zero work if connection setup assumes separate env vars.
- **Nixpacks cache invalidates on `go.mod` changes.** Every dependency update triggers a full module download (3–5 min first build). Predictable but surprising if you expect Docker-layer-style caching.
- **Remote MCP requires browser OAuth for initial auth.** The hosted `mcp.railway.com` endpoint needs a browser-based OAuth flow. In a CI or headless environment, the local CLI-based MCP (`railway mcp`) must be used instead. The docs don't prominently surface this distinction.
- **Railway volumes (persistent disk) have limited region availability.** Any local disk writes should be treated as ephemeral. All persistent data must go to Postgres or an external object store; Railway volumes are not yet a safe assumption at Hobby tier.
- **Railway CLI on Windows has two install paths.** `npm install -g @railway/cli` (requires Node) and a standalone binary installer. The binary installer path is not the first option shown in the docs. This can cause a 15-minute Node install detour on first setup on Windows.

## Operational Story

- **Preview deploys**: Railway automatically creates a preview environment for each GitHub PR branch when connected to a repo. Preview environments share the project's environment variables unless overridden. Preview URLs are public by default — add Railway's environment-level auth (or Cloudflare Access in front) if preview URLs should not be publicly accessible. Fork PRs from external contributors do not trigger preview deploys by default.
- **Secrets**: Environment variables (including `DATABASE_URL`) live in Railway's project environment settings. The Railway CLI can read and set them via `railway variables set KEY=VALUE`. Tokens for the Railway CLI itself are stored in `~/.railway/config.json` on the local machine — not committed to the repo. For CI/CD, a Railway API token scoped to the project is set as a GitHub Actions secret.
- **Rollback**: Railway maintains a deployment history in the dashboard. To revert: `railway rollback` (CLI) or click "Rollback" on a prior deployment in the dashboard — typical time-to-revert is under 2 minutes. Database migrations do not roll back automatically; schema changes should use reversible migration scripts (e.g., `golang-migrate` with `migrate down`).
- **Approval**: The following actions require a human: deleting a Railway project or service, rotating the primary Railway API token, dropping a Postgres database. An agent may perform unattended: deploy, rollback, read logs, set/get environment variables, scale replica count (Pro only).
- **Logs**: `railway logs` streams runtime logs to the terminal. `railway logs --deployment <id>` fetches logs for a specific deployment. For build logs: available in the Railway dashboard under the deployment detail view. MCP-connected agents can query logs via the Railway MCP server's structured tools (`deploymentLogs`, `serviceLogs`).

## Risk Register

| Risk | Source | Likelihood | Impact | Mitigation |
|---|---|---|---|---|
| Hobby $5 credit exhausted — unexpected monthly charge | Devil's advocate | H | L | Set Railway spending limit; monitor usage dashboard weekly; budget $15-20/month for app + DB |
| Platform incident during peak weekend usage (no SLA) | Pre-mortem | M | M | Upgrade to Pro ($20/month) at first retained user cohort; add status page (statuspage.io or Railway's built-in) |
| Nixpacks build fails silently on non-standard Go layout | Devil's advocate | L | M | Add a `railway.toml` with explicit build command (`go build -o app ./cmd/server`); test build locally with `nixpacks build` |
| `DATABASE_URL` format mismatch in Go DSN parsing | Unknown unknowns | M | H | Use `pgx/v5` with `pgx.ParseConfig(os.Getenv("DATABASE_URL"))` from day 1; test connection string locally before first deploy |
| Volume / local disk assumed persistent — data loss | Unknown unknowns | M | H | Treat Railway filesystem as ephemeral; all recipe images and uploads go to Cloudflare R2 or AWS S3 |
| Remote MCP OAuth fails in headless CI environment | Unknown unknowns | L | L | Use Railway CLI local MCP (`railway mcp`) in CI; reserve remote MCP for interactive Claude Code sessions |
| CLI install friction on Windows (Node.js dependency) | Devil's advocate | M | L | Use Railway's standalone binary installer (`winget install Railway.RailwayCLI`); document in AGENTS.md |
| Migration rollback not automatic — schema drift under incident | Pre-mortem | M | H | Use `golang-migrate` with numbered migrations; test `migrate down N` before each deploy to staging |

## Getting Started

1. **Install Railway CLI (Windows):**
   Use the standalone binary installer to avoid the Node.js dependency:
   ```powershell
   winget install Railway.RailwayCLI
   ```
   Verify: `railway --version`

2. **Authenticate and link the project:**
   ```bash
   railway login
   railway init          # creates a new Railway project for weekplate
   railway link          # if project already exists: link current directory
   ```

3. **Provision PostgreSQL:**
   In the Railway dashboard → New Service → Database → PostgreSQL. Railway auto-injects `DATABASE_URL` into the linked service. Confirm with:
   ```bash
   railway variables | grep DATABASE_URL
   ```

4. **Connect the GitHub repo for automatic deploys:**
   In Railway dashboard → Service Settings → Source → Connect to GitHub repo (`weekplate`). Railway will use Nixpacks to detect Go and build automatically on every push to `main`. For the first deploy, verify the build detects Go correctly — if Nixpacks fails, add a `railway.toml`:
   ```toml
   [build]
   builder = "nixpacks"
   buildCommand = "go build -o app ./cmd/server"

   [deploy]
   startCommand = "./app"
   ```

5. **Set the Go HTTP port and any app environment variables:**
   Railway injects `PORT` automatically. Ensure the Go server binds to `os.Getenv("PORT")` (not a hardcoded `:8080`):
   ```bash
   railway variables set APP_ENV=production SESSION_SECRET=<value>
   ```

6. **Verify first deploy and tail logs:**
   ```bash
   railway up            # manual deploy from local, or push to GitHub
   railway logs          # stream runtime logs
   railway status        # check deployment health
   ```

## Out of Scope

The following were not evaluated in this research:
- Docker image configuration (not needed for Nixpacks build; add if Nixpacks fails)
- CI/CD pipeline setup (GitHub Actions with `railway up` or Railway's native GitHub integration covers MVP)
- Production-scale architecture (multi-region, HA, DR)
- Cost optimization beyond MVP scale
