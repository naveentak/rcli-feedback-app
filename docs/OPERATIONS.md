# Operations

Day-to-day ticket management for the r:labs feedback system.

## GitHub Project board

**One-time setup:**

```bash
gh auth refresh -h github.com -s project,read:project
bash scripts/setup-github-project.sh
```

Set columns in GitHub UI: **Inbox → Triaged → In Progress → Done**

New issues land in **Inbox** (no `status:` label). Move via labels using `rcli`.

## Triage workflow

```
Inbox (open, no status label)
  → rcli feedback update N --status triaged
In Progress
  → rcli feedback update N --status in-progress
Done
  → rcli feedback update N --status done   (closes issue)
```

## CLI setup

```bash
cd rcli-feedback-app && make build

export FEEDBACK_API_URL=https://feedback-service-689297474325.europe-west1.run.app
export FEEDBACK_API_KEY=<ops-key-for-any-app>   # from Secret Manager api-keys
export FEEDBACK_APP=rclip                         # or boka, thxbud, etc.
```

## Common commands

```bash
# All open tickets
rcli feedback list --app all

# Per product
rcli feedback list --app rclip --status open

# Inspect
rcli feedback view 42

# Reply (posts GitHub comment)
rcli feedback comment 42 "Fixed in 1.1.1 — please update"

# Triage
rcli feedback update 42 --status in-progress

# AI agent context dump
rcli feedback agents-check --app rclip
```

## Filtering in GitHub

| View | GitHub filter |
|------|---------------|
| r:clip bugs | `label:source:rclip label:type:bug` |
| Open inbox | `is:open -label:status:triaged -label:status:in-progress -label:status:done` |
| In progress | `label:status:in-progress` |

## Smoke test after deploy

```bash
make run          # terminal 1 (or use production URL)
make e2e          # terminal 2
```

## Incident checklist

| Symptom | Check |
|---------|-------|
| Form submit fails | Cloud Run logs; HMAC secret match; CORS origin |
| 401 on signed submit | Clock skew; URL older than 1h; wrong HMAC secret |
| No GitHub issue | `GITHUB_TOKEN` secret; repo permissions |
| Old form still showing | Redeploy `rclip-webapp`; hard-refresh browser |