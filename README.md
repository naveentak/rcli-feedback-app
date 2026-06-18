# rcli — Unified Feedback System

Central feedback & ticketing for **r:clip**, **BoKa**, **ThxBud**, **MamZo**, **GlassCourt**, and future r:labs apps.

**GitHub Issues** = source of truth · **GitHub Projects** = Kanban · **`rcli`** = operator CLI

## Live

| | |
|---|---|
| API | https://feedback-service-689297474325.europe-west1.run.app |
| Issues | https://github.com/naveentak/rcli-feedback-app/issues |
| r:clip form | https://rclip.refactory.co.za/feedback |

## Documentation

Full specification: **[docs/README.md](docs/README.md)**

| Doc | What |
|-----|------|
| [Architecture](docs/ARCHITECTURE.md) | System design & data flow |
| [API](docs/API.md) | Endpoints & auth |
| [Integration](docs/INTEGRATION.md) | Wire a new app (r:clip = reference) |
| [Deployment](docs/DEPLOYMENT.md) | Cloud Run & secrets |
| [Operations](docs/OPERATIONS.md) | Triage with `rcli` |

## Quick start (local)

```bash
cp .env.example .env    # GITHUB_TOKEN, GITHUB_OWNER
make setup-labels
make build && make run  # http://localhost:8080
```

## Quick start (operator)

```bash
export FEEDBACK_API_URL=https://feedback-service-689297474325.europe-west1.run.app
export FEEDBACK_API_KEY=<from-secret-manager>
export FEEDBACK_APP=rclip

rcli feedback list --app all
rcli feedback view 42
rcli feedback update 42 --status triaged
```

## Deploy

```bash
gcloud config set account refactoryza@gmail.com
gcloud config set project rlabs-app
make deploy
```

## Project layout

```
cmd/server/       Go API (Cloud Run)
cmd/rcli/         Operator CLI
internal/         GitHub client, auth, domain logic
pkg/cli/          Cobra commands
web/              Generic feedback page
docs/             Formal specification
scripts/          Setup & deploy helpers
```