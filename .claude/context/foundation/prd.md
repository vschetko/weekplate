---
project: weekplate
version: 1
status: draft
created: 2026-06-07
context_type: greenfield
product_type: web-app
target_scale:
  users: medium
  qps: "# TODO: not specified — see Open Questions"
  data_volume: "# TODO: not specified — see Open Questions"
timeline_budget:
  mvp_weeks: 3
  hard_deadline: "2026-06-29"
  after_hours_only: true
---

## Vision & Problem Statement

Tired, busy people — working adults in a couple or small household (2–3 people) — hit the weekend without a meal plan. Faced with the question "what do we eat next week?", they default to takeout and delivery because planning feels like work on top of work. The real cost is financial (delivery fees) and health-related (poor food choices).

The insight driving this product: no existing tool connects meal planning directly to a grocery list in a single, low-friction flow. Users who try to plan end up juggling two separate tools — a recipe browser and a manual list — and the habit breaks down at that gap.

This app closes that gap: plan meals for the week, get a ready-to-use grocery list, cook with confidence.

## User & Persona

**Primary persona:** A couple or small household (2–3 people), working adults, low energy on weekends. They share one household account — flat access model, no per-member permissions. They know roughly what they like to eat but struggle to translate that into an actionable weekly plan and a single shopping trip. They are not food enthusiasts or cooking hobbyists — they want to eat well without thinking hard about it.

## Success Criteria

### Primary
A couple signs up, sets their preferred daily calorie target (from predefined levels or custom), sets food exclusions, requests a weekly meal plan, receives a linked grocery list, adjusts the plan via removal with auto-replacement offer, views the completed grocery list, and sees today's meal — in a single session. If this flow works end-to-end, the product is proven.

### Secondary
- Cooking instructions inline: tapping a meal shows a basic recipe so the app serves as the cooking companion during the week.

### Guardrails
- Grocery list must be complete — no missing ingredients when the user shops. Partial lists destroy trust on first use.
- The weekly plan must persist across browser refreshes and app restarts. Losing a plan after setup is a one-strike failure.

## User Stories

### US-01: Generate a weekly meal plan

- **Given** a logged-in user with a calorie target set and at least one food exclusion applied (or none),
- **When** they request a weekly meal plan,
- **Then** the app displays a list of meals for each day of the upcoming week that fits their daily calorie target and excludes flagged ingredients, along with a linked grocery list containing all required ingredients (deduplicated and summed).

#### Acceptance Criteria
# TODO: acceptance criteria for US-01 — see Open Questions

## Functional Requirements

### Authentication & Profile
- FR-001: User can create an account (email/password or OAuth). Priority: must-have
  > Socrates: Counter-argument considered: "account creation adds drop-off before any value is shown." Resolution: acknowledged; a guest/trial mode is a valid v2 consideration but auth is required in v1 for plan persistence across devices and household sharing.

- FR-002: User can log in to an existing account. Priority: must-have
  > Socrates: Handled together with FR-001. Auth is foundational; stands as written.

- FR-003: User can set a preferred daily calorie target (custom input or pick from predefined levels, e.g. 1500 / 2000 / 2500 kcal). Priority: must-have
  > Socrates: Counter-argument considered: "asking for weight is privacy-invasive and BMR calculation adds accuracy risk." Resolution: FR updated — app asks for a preferred daily calorie target directly, skipping sex/age/weight collection and BMR computation. Simpler, less invasive, and eliminates formula accuracy risk.

### Meal Planning
- FR-004: User can request a weekly meal plan. Priority: must-have
  > Socrates: Counter-argument considered: "7-day plans go stale mid-week as real life disrupts them." Resolution: accepted risk for v1; mid-week regeneration and day-level flexibility are v2 features. Recorded in Open Questions.

- FR-005: User can set food exclusions (allergens and dislikes) that the app filters out of all meal suggestions. Priority: must-have
  > Socrates: Counter-argument considered: "too many filter options cause decision paralysis." Resolution: scope clarified — filters are exclusion-only (what to remove, not what to select). Target: a short list of common allergens plus a free-text "don't include" field. No style/cuisine selectors in v1.

- FR-006: User can view the generated weekly meal list. Priority: must-have
  > Socrates: No counter-argument raised; stands as written.

- FR-007: User can remove a meal from the plan; the app removes linked ingredients from the grocery list and offers an automatic replacement for the empty day slot. Priority: must-have
  > Socrates: Counter-argument considered: "removing a meal leaves a day without a plan." Resolution: FR updated — on removal the app offers to auto-fill the gap with a suggested replacement from the available recipe pool, rather than leaving the slot empty.

- FR-008: User can view today's assigned meal. Priority: must-have
  > Socrates: Counter-argument considered: "first-session users see an empty state with no guidance." Resolution: the empty state for this view must include a clear call-to-action to generate a plan. UX requirement, not a new FR.

### Grocery List
- FR-009: User can view the grocery list linked to the current week's meal plan, with ingredients deduplicated and quantities summed across all meals. Priority: must-have
  > Socrates: No counter-argument raised; stands as written. Deduplication and quantity summing is a data quality requirement built into the implementation, not a separate FR.

### Recipes
- FR-011: User can view cooking instructions for a meal. Priority: nice-to-have
  > Socrates: Counter-argument considered: "every meal needing instructions doubles content work and slows library growth." Resolution: kept as nice-to-have with a content quality gate — recipes without complete instructions should not be published to the library. This limits initial library size but protects quality.

## Non-Functional Requirements

- **Grocery list completeness**: the list must contain every ingredient needed for the week's meals, deduplicated and with quantities summed. The target is a single shopping trip with no mid-week store returns.
- **Mobile-first responsive web**: the app must be fully usable on a phone screen without a native install. The primary use context is a phone on the weekend.
- **Cross-device plan persistence**: the user's calorie target, food exclusions, and current week's plan must be accessible from any device after logging in — logging out and back in on a different device must show the same plan.

## Business Logic

The app recommends a 7-day meal plan that avoids flagged ingredients and approximates the user's calorie target, then automatically derives the complete consolidated grocery list from that plan.

Rule shape: **recommendation + derivation**. The app recommends meals filtered by food exclusions (hard constraint) and guided by the calorie target (soft preference); it then calculates the consolidated grocery list (ingredients summed and deduplicated across all 7 selected meals). No additional rules in v1: variety enforcement, macro balancing, and ranking by preference are explicitly deferred.

Inputs the rule consumes: user's daily calorie target (soft preference); set of excluded ingredients/allergens (hard constraint); available recipe library with per-recipe calorie counts and ingredient lists.
Output: a 7-day meal assignment + a deduplicated, summed ingredient list ready for shopping.

## Access Control

Authentication: email/password or a supported third-party sign-in provider. Users create an account; the household shares one account (one login, one plan). No role separation — the model is flat: anyone logged in can view and edit the shared weekly plan and grocery list.

The smallest access model that makes the MVP useful: a single shared household account with no per-member permissions. Multi-user role management is explicitly out of scope for v1.

## Non-Goals

- **No social features**: the app serves one household. No plan sharing with other households, no public profiles, no community features.
- **No native mobile app**: responsive web only. No iOS/Android native build, no app store submission in v1.
- **No grocery delivery integration**: the grocery list is for the user to shop however they choose. No API connection to any delivery service.
- **No macro tracking beyond calories**: nutritional logic stops at daily calorie target. Protein/carb/fat breakdown is explicitly a v2 concern.
- **No custom meal creation**: meals come from a pre-built recipe library only. Users cannot define their own recipes in v1.
- **No saved favorites**: meal favorites are out of scope for v1. Favorites only add value after several plan cycles; deferred to v2 backlog.

## Open Questions

1. **Will a 7-day plan go stale mid-week?** A 7-day plan may become inaccurate as real life disrupts it. Mid-week plan regeneration and day-level flexibility are deferred to v2. Risk: users may abandon plans by Wednesday. Owner: user. Action: monitor in v1 user feedback.

2. **What is the expected requests-per-second (qps) load?** Not specified during shaping. Required to complete `target_scale.qps` in frontmatter. Owner: user.

3. **What is the expected data volume?** Not specified during shaping. Required to complete `target_scale.data_volume` in frontmatter. Owner: user.

4. **What are the acceptance criteria for US-01?** The Given/When/Then block exists but specific pass/fail acceptance criteria are not captured in shape-notes. Owner: user.
