# rcli Feedback System — Documentation

Formal specification for the unified feedback & ticketing platform across r:labs products.

## Documents

| Doc | Audience | Purpose |
|-----|----------|---------|
| [Architecture](ARCHITECTURE.md) | Everyone | System design, data flow, components |
| [API Reference](API.md) | App & backend devs | Endpoints, auth, payloads |
| [App Integration](INTEGRATION.md) | App devs | How to wire a product (r:clip is the reference) |
| [Deployment](DEPLOYMENT.md) | Backend / DevOps | GCP Cloud Run, secrets, redeploy |
| [Operations](OPERATIONS.md) | You + agents | Triage, CLI, GitHub Projects |
| [ADR-001](ADR-001-github-source-of-truth.md) | — | Why GitHub Issues is the database |
| [ADR-002](ADR-002-signed-web-feedback.md) | — | Why macOS apps use HMAC + web forms |

## Live services

| Service | URL |
|---------|-----|
| Feedback API | https://feedback-service-689297474325.europe-west1.run.app |
| GitHub repo | https://github.com/naveentak/rcli-feedback-app |
| r:clip form | https://rclip.refactory.co.za/feedback |

## Supported apps

| App ID | Product | Integration status |
|--------|---------|-------------------|
| `rclip` | r:clip | Live (signed web form) |
| `boka` | BoKa | Labels ready |
| `thxbud` | ThxBud | Labels ready |
| `mamzo` | MamZo | Labels ready |
| `glasscourt` | GlassCourt | Labels ready |

## Quick commands

```bash
make run              # local server
make deploy           # Cloud Run
make e2e              # smoke test
bash scripts/setup-github-labels.sh
bash scripts/setup-github-project.sh   # requires gh project scope
```