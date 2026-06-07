# Repository Guidelines

`weekplate` is a meal-planning web app for couples and small households: a user builds a weekly meal plan and gets a deduplicated grocery list in one flow. Stack: Go standard library, server-side HTML templates, single binary deployed to Railway. MVP deadline: 2026-06-29.

## Hard Rules

- **Grocery list must always be complete.** Every ingredient for every meal in the current week's plan must appear in the grocery list, deduplicated and quantities summed. A partial list is a trust-breaking one-strike failure — do not ship any feature that can produce an incomplete list.
- **Plans must persist server-side.** Calorie target, food exclusions, and the current weekly plan must survive a page refresh and a device switch. Do not store authoritative state only in browser memory or localStorage.
- **All plan and grocery routes require authentication.** No anonymous access to plan data. Auth model is flat: one shared household account, no roles or per-member permissions.
- **Food exclusions are hard constraints; calorie target is a soft preference.** Never include a meal that contains an excluded ingredient regardless of calorie fit. Calorie target guides selection, not hard filtering.
- **No custom recipe creation in v1.** Meals come from a pre-built library only. Do not implement endpoints that accept user-defined recipes.
- **Mobile-first layout.** All HTML templates must be fully usable on a phone-sized viewport. Desktop is secondary.

## Project Structure

Layout not established yet. `cmd/server/` — entry point. `internal/` — domain logic, handlers, templates, data access. Module path for all imports: `github.com/vschetko/weekplate` — see @go.mod.

See @context/foundation/prd.md for full requirements and @context/foundation/tech-stack.md for stack rationale and deployment notes.

## Build & Development

- `go build ./...` — compile all packages
- `go test ./...` — run all tests
- `go vet ./...` — static analysis
- `go fmt ./...` — format (run before committing)

No Makefile yet; prefer adding `make build` / `make test` targets rather than inventing ad hoc scripts.

## Coding Style

Go 1.26.4. Standard library-first — introduce chi or gorilla/mux only when `net/http.ServeMux` is insufficient. No third-party ORM. `gofmt` formatting required.

## Testing

Use the standard `testing` package. Mock dependencies with `gomock` (`go.uber.org/mock/gomock`); generate mocks via `mockgen`. No other external test framework.

## Commit & Pull Request Guidelines

No commit convention established in history yet. Use Conventional Commits prefixes (`feat:`, `fix:`, `chore:`, `test:`) going forward. Remote: `github.com/vschetko/weekplate`.

## Deployment

Railway (Nixpacks auto-build; `railway.toml` provides explicit fallback). Build target: `CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/server`. Railway injects `PORT` and `DATABASE_URL` automatically. See @context/foundation/infrastructure.md for platform decision and risk register.

**CLI install (Windows):** `npm install -g @railway/cli` — the `winget` package `Railway.RailwayCLI` does not exist. Node.js 16+ required.
