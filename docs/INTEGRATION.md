# App Integration Guide

**Reference implementation:** r:clip (`rclip-macos-app` + `rclip-webapp`)

## Choose an integration pattern

| Pattern | Best for | API key in app? |
|---------|----------|-----------------|
| **A — Signed web form** | macOS menu bar / utility apps with branded web UX | No |
| **B — API key direct** | CLI tools, servers, mobile apps with backend | Yes (build-time) |

r:clip uses **Pattern A**. This guide documents that flow; Pattern B uses `POST /api/v1/feedback` — see [API.md](API.md).

---

## Pattern A — Signed web form (r:clip)

### Overview

```
macOS Preferences → signed URL → product.com/feedback → POST /feedback/signed → GitHub Issue
```

### Step 1 — macOS app: `FeedbackTokenManager`

Location: `{App}/Core/FeedbackTokenManager.swift`

Responsibilities:

- Persist `deviceId` in `UserDefaults`
- Sign `"{deviceId}|{timestamp}"` with HMAC-SHA256
- Open browser: `https://{product}.refactory.co.za/feedback?did&ts&sig&app&type&v&os`

Query parameters:

| Param | Value |
|-------|-------|
| `did` | Device UUID |
| `ts` | Unix timestamp (seconds) |
| `sig` | HMAC hex digest |
| `app` | App ID (`rclip`) |
| `type` | `bug` or `feature` |
| `v` | `CFBundleShortVersionString` |
| `os` | `ProcessInfo.operatingSystemVersionString` |

Preferences buttons call:

```swift
FeedbackTokenManager.shared.openFeedback(type: "bug")      // or "feature"
```

**Sandbox:** requires `com.apple.security.network.client` entitlement (outbound HTTPS only).

### Step 2 — Web app: `FeedbackForm.tsx`

Location: `{product}-webapp/src/components/FeedbackForm.tsx`

Responsibilities:

- Parse URL query params; show "use the app" message if invalid/missing
- Render branded form (type toggle, title, description, email, version badges)
- POST to `{VITE_FEEDBACK_API_URL}/api/v1/feedback/signed`
- Success state: show ticket number + **View on GitHub** link

Env:

```bash
VITE_FEEDBACK_API_URL=https://feedback-service-689297474325.europe-west1.run.app
```

Route: `/feedback` (see `App.tsx`)

### Step 3 — Backend onboarding (ops team only)

1. Add app ID to `internal/feedback/types.go` if new
2. Add `source:{app}` label: `make setup-labels`
3. Add API key to Secret Manager `api-keys` (for CLI/ops, not the web form)
4. Ensure `FEEDBACK_HMAC_SECRET` is set on Cloud Run (shared across r:labs macOS apps)
5. Add product origin to CORS in `cmd/server/main.go` if new domain
6. Deploy: `make deploy` + `firebase deploy` (webapp)

### Step 4 — Deploy order

```
1. rcli-feedback-app  →  Cloud Run
2. {product}-webapp   →  Firebase / hosting
3. {product}-macos-app →  App Store build (only if FeedbackTokenManager changed)
```

---

## Onboarding checklist (new app)

```
[ ] App ID chosen (lowercase, e.g. boka)
[ ] source:{app} label created
[ ] API key generated and added to Secret Manager
[ ] FeedbackTokenManager.swift (or equivalent) in macOS app
[ ] FeedbackForm.tsx on product domain /feedback route
[ ] VITE_FEEDBACK_API_URL set in webapp
[ ] CORS origin added for product domain
[ ] End-to-end test: app → form → GitHub issue
[ ] GitHub Project board receives new issues
```

---

## Copy-paste template: web form submit handler

```typescript
const res = await fetch(`${API_URL}/api/v1/feedback/signed`, {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    deviceId: token.did,
    timestamp: token.ts,
    signature: token.sig,
    appId: token.app,
    type: feedbackType,       // 'bug' | 'feature'
    title,
    description,
    email: email || undefined,
    appVersion: token.v,
    osVersion: token.os,
  }),
})
const ticket = await res.json()
// ticket.number, ticket.url
```

---

## Per-repo pointers

| Repo | Doc |
|------|-----|
| `rclip-macos-app` | `FEEDBACK_SETUP.md` |
| `rclip-webapp` | `.env.example` |
| `rcli-feedback-app` | This guide + [API.md](API.md) |