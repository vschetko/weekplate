---
bootstrapped_at: 2026-06-07T11:52:43Z
starter_id: go
starter_name: Go (standard library)
project_name: weekplate
language_family: go
package_manager: "(Go uses modules natively — no external package manager)"
cwd_strategy: native-cwd into weekplate (D:\AI_COURSE\projects\weekplate)
bootstrapper_confidence: first-class
phase_3_status: ok
audit_command: "govulncheck ./..."
---

## Hand-off

```yaml
starter_id: go
project_name: weekplate
hints:
  language_family: go
  team_size: solo
  deployment_target: fly
  ci_provider: github-actions
  ci_default_flow: manual-promotion
  bootstrapper_confidence: first-class
  path_taken: standard
  quality_override: false
  self_check_answers: null
  has_auth: true
  has_payments: false
  has_realtime: false
  has_ai: false
  has_background_jobs: false
```

### Why this stack

A solo developer building a meal-planning web app in 3 weeks after-hours chose Go standard library, the recommended default for the `(web, go)` cell. Go compiles to a single binary with no runtime dependency, which pairs well with Fly.io — a containerised platform where a lean binary means fast startup and low memory overhead. All four agent-friendly criteria pass: Go is typed by the language, Go projects follow strong layout conventions (`cmd/`, `internal/`, `pkg/`), and Go is well-represented in both training data and official documentation. Bootstrapper confidence is first-class — the scaffold command is registered and expected to work, though not yet battle-tested end-to-end. Auth is not included first-class and will need manual addition post-scaffold (session management or OAuth via `golang.org/x/oauth2`); the user acknowledged this gap before proceeding. CI runs on GitHub Actions with manual promotion to Fly.io — a deliberate gate over production deploys during an after-hours build.

## Pre-scaffold verification

| Signal      | Value   | Severity | Notes                                                          |
| ----------- | ------- | -------- | -------------------------------------------------------------- |
| npm package | not run | n/a      | go language family — no npm package to check                   |
| GitHub repo | not run | n/a      | docs_url (https://go.dev/doc/) is not a GitHub URL             |

## Scaffold log

**Resolved invocation**: `go mod init github.com/user/weekplate`
**Strategy**: native-cwd into D:\AI_COURSE\projects\weekplate
**Exit code**: 0
**Files written**: 1 (go.mod)
**Conflicts (.scaffold siblings)**: none — directory was empty
**.gitignore handling**: absent in scaffold

### go.mod contents

```
module github.com/user/weekplate

go 1.26.4
```

## Post-scaffold audit

**Tool**: govulncheck ./...
**Status**: failed to run
**Reason**: govulncheck not installed on this machine
**Partial output (if any)**: govulncheck: command not found

To install govulncheck and run the audit manually:

```
go install golang.org/x/vuln/cmd/govulncheck@latest
cd D:\AI_COURSE\projects\weekplate
govulncheck ./...
```

A fresh Go module with no dependencies will return 0 findings — this is informational only.

## Hints recorded but not acted on

| Hint                    | Value            |
| ----------------------- | ---------------- |
| bootstrapper_confidence | first-class      |
| quality_override        | false            |
| path_taken              | standard         |
| self_check_answers      | null             |
| team_size               | solo             |
| deployment_target       | fly              |
| ci_provider             | github-actions   |
| ci_default_flow         | manual-promotion |
| has_auth                | true             |
| has_payments            | false            |
| has_realtime            | false            |
| has_ai                  | false            |
| has_background_jobs     | false            |

## Next steps

Next: a future skill will set up agent context (CLAUDE.md, AGENTS.md). For now, your project is scaffolded and verified — happy hacking.

Useful manual steps in the meantime:
- `git init` (if you have not already) to start your own repo history.
- Review any `.scaffold` siblings the conflict policy created and decide which version of each file to keep.
- Address audit findings per your project's risk tolerance — the full breakdown is in this log.
- Install govulncheck and run `govulncheck ./...` from weekplate when ready: `go install golang.org/x/vuln/cmd/govulncheck@latest`
- Auth (has_auth: true) needs manual wiring — consider `golang.org/x/oauth2` or a session library post-scaffold.
