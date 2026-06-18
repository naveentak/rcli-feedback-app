# Unified Feedback System

Central feedback & ticketing for **r:clip**, **BoKa**, **ThxBud**, **MamZo**, **GlassCourt**, and future apps. GitHub Issues is the source of truth; GitHub Projects is the Kanban board.

## Architecture

```
Apps / CLI / Web Page
        │
        ▼
  Feedback Service (Go + Gin)
        │
        ▼
  GitHub Issues + Projects
```

## Quick Start

### 1. GitHub setup

Create a repo named `feedback` and a Personal Access Token with `repo` scope.

```bash
cp .env.example .env
# Edit .env with your GITHUB_TOKEN, GITHUB_OWNER, API_KEYS

make setup-labels   # creates labels via gh CLI
```

Create a GitHub Project with columns: **Inbox → Triaged → In Progress → Done**.

### 2. Run the service

```bash
make build
make run
```

Service runs at `http://localhost:8080`. User page at `/`.

### 3. Use the CLI

```bash
export FEEDBACK_API_URL=http://localhost:8080
export FEEDBACK_API_KEY=secret-rclip-key
export FEEDBACK_APP=rclip

# Submit
./bin/rcli feedback submit \
  --app rclip --type bug \
  --title "Crash on export" \
  --description "Steps to reproduce..."

# List
./bin/rcli feedback list --app rclip --status open
./bin/rcli feedback list --app all

# View / comment / update
./bin/rcli feedback view 42
./bin/rcli feedback comment 42 "Looking into this"
./bin/rcli feedback update 42 --status in-progress

# Agent context dump
./bin/rcli feedback agents-check --app rclip
```

## API

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/api/v1/feedback` | API Key | Submit ticket |
| GET | `/api/v1/feedback` | — | List tickets |
| GET | `/api/v1/feedback/:number` | — | Get ticket |
| POST | `/api/v1/feedback/:number/comments` | API Key | Add comment |
| PATCH | `/api/v1/feedback/:number` | API Key | Update status |

Auth headers: `X-API-Key` + `X-App` (or `Authorization: Bearer <key>`).

## Giving keys to client apps (r:clip, BoKa, etc.)

Client apps have **no gcloud access**. Only whoever runs this backend manages keys in GCP Secret Manager.

To onboard an app:

1. Generate a key: `openssl rand -hex 24`
2. Add to Secret Manager `api-keys`: `rclip:NEWKEY,...` then redeploy Cloud Run
3. Send the **r:clip key only** to the app developer via password manager or secure channel
4. They paste it into their local `Secrets.xcconfig` (gitignored) before Release builds

Each app's key only allows **submitting** tickets for that app.

## Labels

- `source:rclip`, `source:boka`, `source:thxbud`, `source:mamzo`, `source:glasscourt`
- `type:bug`, `type:feature-request`
- `status:triaged`, `status:in-progress`, `status:done`

## Project Structure

```
cmd/
  server/     # Gin HTTP service
  rcli/       # CLI entrypoint
internal/
  config/     # Env-based config
  github/     # GitHub API client
  feedback/   # Domain logic + handlers
  auth/       # API key middleware
pkg/cli/      # Cobra commands
web/          # User-facing page
scripts/      # GitHub setup helpers
```