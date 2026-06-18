# API Reference

Base URL (production): `https://feedback-service-689297474325.europe-west1.run.app`

## Endpoints

### `GET /health`

```json
{ "status": "ok" }
```

---

### `POST /api/v1/feedback/signed` — macOS web forms

**Auth:** HMAC signature (no API key in browser).

**Request:**

```json
{
  "deviceId": "uuid",
  "timestamp": "1718700000",
  "signature": "hex-hmac-sha256",
  "appId": "rclip",
  "type": "bug",
  "title": "Short summary",
  "description": "Details…",
  "email": "optional@example.com",
  "appVersion": "1.1.0",
  "osVersion": "macOS 15.0"
}
```

| Field | Notes |
|-------|-------|
| `type` | `bug` or `feature` (`feature` → `feature-request` label) |
| `appId` | Must be a registered app: `rclip`, `boka`, `thxbud`, `mamzo`, `glasscourt` |
| `signature` | `HMAC-SHA256("{deviceId}\|{timestamp}", FEEDBACK_HMAC_SECRET)` as hex |

**Validation:** timestamp within 1 hour; signature must match.

**Response `201`:**

```json
{
  "number": 42,
  "title": "…",
  "url": "https://github.com/naveentak/rcli-feedback-app/issues/42",
  "app": "rclip",
  "type": "bug",
  "state": "open",
  "labels": ["source:rclip", "type:bug"]
}
```

---

### `POST /api/v1/feedback` — API key submit

**Headers:** `X-API-Key`, `X-App` (or `Authorization: Bearer <key>`)

**Request:**

```json
{
  "app": "rclip",
  "type": "bug",
  "title": "…",
  "description": "…",
  "reporter": "optional"
}
```

| `type` | Values |
|--------|--------|
| | `bug`, `feature-request` |

---

### `GET /api/v1/feedback`

**Query:** `app`, `status`, `type`, `state` (`open`|`closed`|`all`)

Public read. No auth.

---

### `GET /api/v1/feedback/:number`

Public read. Returns single ticket.

---

### `POST /api/v1/feedback/:number/comments`

**Auth:** API key required.

```json
{ "body": "Comment text" }
```

---

### `PATCH /api/v1/feedback/:number`

**Auth:** API key required.

```json
{ "status": "triaged" }
```

| Status | Effect |
|--------|--------|
| `triaged` | Adds `status:triaged` label |
| `in-progress` | Adds `status:in-progress` label |
| `done` | Adds `status:done` label, closes issue |

---

## CORS

Allowed origins for browser POST:

- `https://rclip.refactory.co.za`
- `http://localhost:5173`
- `http://localhost:3000`

Add more via `ALLOWED_ORIGINS` env var (comma-separated).

## Errors

```json
{ "error": "human-readable message" }
```

| Code | Typical cause |
|------|---------------|
| 400 | Invalid app, type, or missing fields |
| 401 | Bad API key or invalid/expired HMAC |
| 404 | Issue not found |
| 503 | HMAC secret not configured |