---
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
---

## Why this stack

A solo developer building a meal-planning web app in 3 weeks after-hours chose Go standard library, the recommended default for the `(web, go)` cell. Go compiles to a single binary with no runtime dependency, which pairs well with Fly.io — a containerised platform where a lean binary means fast startup and low memory overhead. All four agent-friendly criteria pass: Go is typed by the language, Go projects follow strong layout conventions (`cmd/`, `internal/`, `pkg/`), and Go is well-represented in both training data and official documentation. Bootstrapper confidence is first-class — the scaffold command is registered and expected to work, though not yet battle-tested end-to-end. Auth is not included first-class and will need manual addition post-scaffold (session management or OAuth via `golang.org/x/oauth2`); the user acknowledged this gap before proceeding. CI runs on GitHub Actions with manual promotion to Fly.io — a deliberate gate over production deploys during an after-hours build.
