# ADR-002: Signed Web Forms for macOS Utility Apps

**Status:** Accepted  
**Date:** 2026-06-18

## Context

macOS utility apps (r:clip) should offer polished, branded feedback UX — not native AppKit forms or raw API keys embedded in the binary. The existing `rclip.refactory.co.za/feedback` form was built for this purpose with HMAC URL signing.

## Decision

**Pattern A:** macOS app opens a **HMAC-signed browser URL** → product web form → `POST /api/v1/feedback/signed` → GitHub Issue.

- Shared `FEEDBACK_HMAC_SECRET` between macOS app and Cloud Run
- Signature over `{deviceId}|{timestamp}`, 1-hour expiry
- No API key in the macOS app or browser
- Web form POSTs to Cloud Run with CORS for product domain

## Consequences

**Positive**

- Keeps the designed web UX (animations, copy, branding)
- macOS app stays thin — one `FeedbackTokenManager`, no secrets in binary
- Only app users can reach the form (unsigned URLs see "use the app" message)
- Same pattern works for Harbor and future r:labs macOS apps

**Negative**

- Requires browser context switch (acceptable for infrequent feedback)
- HMAC secret rotation requires coordinated macOS + server deploy
- Supabase edge function deprecated for r:clip

## Alternatives considered

| Option | Rejected because |
|--------|------------------|
| Native AppKit form | Poor UX vs existing web form; user rejected |
| API key in macOS app | Extractable from binary; wrong for utility apps |
| Public unauthenticated form | Spam risk without HMAC gate |