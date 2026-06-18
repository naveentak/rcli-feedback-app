# ADR-001: GitHub Issues as Source of Truth

**Status:** Accepted  
**Date:** 2026-06-18

## Context

r:labs has multiple products (r:clip, BoKa, ThxBud, MamZo, GlassCourt) that need a unified feedback channel. Previous r:clip flow used Supabase `tickets` table — a second database to maintain.

## Decision

Use **GitHub Issues** in `naveentak/rcli-feedback-app` as the only persistent store. GitHub Projects provides Kanban. A thin Go service translates API calls into GitHub API operations.

## Consequences

**Positive**

- Zero database ops; issues are visible, searchable, linkable
- Native integration with dev workflow, PRs, and agents
- `rcli` CLI can manage tickets without custom admin UI
- Labels replace schema columns (`source`, `type`, `status`)

**Negative**

- GitHub API rate limits (negligible at current scale)
- Label index lag (~2–5s after create)
- No custom fields beyond labels + issue body

## Alternatives considered

| Option | Rejected because |
|--------|------------------|
| Supabase tickets | Duplicate source of truth; already migrating away |
| Linear / Jira | Cost + context switch away from code |
| Local SQLite | Not accessible to agents or web |