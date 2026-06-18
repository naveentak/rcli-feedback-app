# Deployment

## Environments

| Env | API URL | GitHub repo |
|-----|---------|-------------|
| Production | `https://feedback-service-689297474325.europe-west1.run.app` | `naveentak/rcli-feedback-app` |
| Local | `http://localhost:8080` | same (via `.env`) |

## GCP (production)

| Setting | Value |
|---------|-------|
| Account | `refactoryza@gmail.com` |
| Project | `rlabs-app` |
| Region | `europe-west1` |
| Service | `feedback-service` |
| Image | `europe-west1-docker.pkg.dev/rlabs-app/rcli/feedback:latest` |

### Secrets (Secret Manager)

| Secret | Purpose |
|--------|---------|
| `github-token` | GitHub API (repo scope) |
| `api-keys` | Per-app keys: `rclip:KEY,boka:KEY,…` |
| `feedback-hmac-secret` | HMAC for signed web forms |

### Environment variables

| Var | Production value |
|-----|------------------|
| `GITHUB_OWNER` | `naveentak` |
| `GITHUB_REPO` | `rcli-feedback-app` |
| `DEV_MODE` | `false` |
| `PUBLIC_SUBMIT` | `true` |
| `GIN_MODE` | `release` |

### Deploy

```bash
gcloud config set account refactoryza@gmail.com
gcloud config set project rlabs-app
make deploy
```

`make deploy` runs `scripts/deploy-gcp.sh` which:

1. Enables APIs
2. Builds via Cloud Build (no local Docker required)
3. Deploys to Cloud Run with secrets mounted

### Manual redeploy

```bash
gcloud builds submit --tag europe-west1-docker.pkg.dev/rlabs-app/rcli/feedback:latest --region=europe-west1
gcloud run deploy feedback-service \
  --image=europe-west1-docker.pkg.dev/rlabs-app/rcli/feedback:latest \
  --region=europe-west1 \
  --set-secrets="GITHUB_TOKEN=github-token:latest,API_KEYS=api-keys:latest,FEEDBACK_HMAC_SECRET=feedback-hmac-secret:latest"
```

### Rotating secrets

**GitHub token** (strip newlines):

```bash
gh auth token | tr -d '\n' | gcloud secrets versions add github-token --data-file=-
gcloud run services update feedback-service --region=europe-west1 --update-secrets=GITHUB_TOKEN=github-token:latest
```

**API keys:**

```bash
echo -n 'rclip:NEW,boka:NEW,...' | gcloud secrets versions add api-keys --data-file=-
gcloud run services update feedback-service --region=europe-west1 --update-secrets=API_KEYS=api-keys:latest
```

**HMAC secret** (must match macOS apps + web forms):

```bash
echo -n 'your-hex-secret' | gcloud secrets versions add feedback-hmac-secret --data-file=-
gcloud run services update feedback-service --region=europe-west1 --update-secrets=FEEDBACK_HMAC_SECRET=feedback-hmac-secret:latest
```

## Product web apps (Firebase)

r:clip example:

```bash
cd rclip-webapp
# ensure VITE_FEEDBACK_API_URL is in .env
npm run build
firebase deploy --only hosting --project rlabs-app
```

## Local development

```bash
cp .env.example .env
# GITHUB_TOKEN, GITHUB_OWNER, API_KEYS, DEV_MODE=true
make run
```

Test signed endpoint locally — set `FEEDBACK_HMAC_SECRET` in `.env` to match macOS app.

## Cost

Cloud Run: scale-to-zero, 256Mi, max 3 instances. Expect ~$0–2/month at current volume. See architecture notes in project README.