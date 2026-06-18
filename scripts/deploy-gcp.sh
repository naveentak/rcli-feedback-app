#!/usr/bin/env bash
# Deploy feedback service to Cloud Run on rlabs-app.
set -euo pipefail

PROJECT_ID="${GCP_PROJECT:-rlabs-app}"
REGION="${GCP_REGION:-europe-west1}"
SERVICE="feedback-service"
IMAGE="${REGION}-docker.pkg.dev/${PROJECT_ID}/rcli/feedback:latest"
SA_NAME="feedback-run"

echo "→ Project: $PROJECT_ID | Region: $REGION | Account: $(gcloud config get-value account)"

gcloud config set project "$PROJECT_ID"

echo "→ Enabling APIs..."
gcloud services enable \
  run.googleapis.com \
  cloudbuild.googleapis.com \
  artifactregistry.googleapis.com \
  secretmanager.googleapis.com \
  --quiet

echo "→ Artifact Registry..."
if ! gcloud artifacts repositories describe rcli --location="$REGION" &>/dev/null; then
  gcloud artifacts repositories create rcli \
    --repository-format=docker \
    --location="$REGION" \
    --description="rcli services"
fi

echo "→ Secrets..."
token_no_newline() { gh auth token | tr -d '\n'; }
if ! gcloud secrets describe github-token &>/dev/null; then
  token_no_newline | gcloud secrets create github-token --data-file=-
else
  token_no_newline | gcloud secrets versions add github-token --data-file=-
fi

if ! gcloud secrets describe api-keys &>/dev/null; then
  API_KEYS_VAL="${API_KEYS:-rclip:prod-rclip-key,boka:prod-boka-key,thxbud:prod-thxbud-key,mamzo:prod-mamzo-key,glasscourt:prod-glasscourt-key}"
  echo -n "$API_KEYS_VAL" | gcloud secrets create api-keys --data-file=-
else
  API_KEYS_VAL="${API_KEYS:-}"
  if [ -n "$API_KEYS_VAL" ]; then
    echo -n "$API_KEYS_VAL" | gcloud secrets versions add api-keys --data-file=-
  fi
fi

echo "→ Service account..."
if ! gcloud iam service-accounts describe "${SA_NAME}@${PROJECT_ID}.iam.gserviceaccount.com" &>/dev/null; then
  gcloud iam service-accounts create "$SA_NAME" --display-name="Feedback Cloud Run"
fi
SA_EMAIL="${SA_NAME}@${PROJECT_ID}.iam.gserviceaccount.com"

bind_secret() {
  local secret=$1 attempt
  for attempt in 1 2 3; do
    if gcloud secrets add-iam-policy-binding "$secret" \
      --member="serviceAccount:${SA_EMAIL}" \
      --role="roles/secretmanager.secretAccessor" --quiet 2>/dev/null; then
      return 0
    fi
    sleep 5
  done
  gcloud secrets add-iam-policy-binding "$secret" \
    --member="serviceAccount:${SA_EMAIL}" \
    --role="roles/secretmanager.secretAccessor" --quiet
}
for SECRET in github-token api-keys; do bind_secret "$SECRET"; done

echo "→ Building & pushing image (Cloud Build)..."
gcloud builds submit --tag "$IMAGE" --project="$PROJECT_ID" --region="$REGION" --quiet

echo "→ Deploying to Cloud Run..."
gcloud run deploy "$SERVICE" \
  --image="$IMAGE" \
  --region="$REGION" \
  --platform=managed \
  --allow-unauthenticated \
  --service-account="$SA_EMAIL" \
  --port=8080 \
  --memory=256Mi \
  --cpu=1 \
  --min-instances=0 \
  --max-instances=3 \
  --set-secrets="GITHUB_TOKEN=github-token:latest,API_KEYS=api-keys:latest,FEEDBACK_HMAC_SECRET=feedback-hmac-secret:latest" \
  --set-env-vars="GITHUB_OWNER=naveentak,GITHUB_REPO=rcli-feedback-app,DEV_MODE=false,PUBLIC_SUBMIT=true,GIN_MODE=release" \
  --quiet

URL=$(gcloud run services describe "$SERVICE" --region="$REGION" --format='value(status.url)')
echo ""
echo "✓ Deployed: $URL"
echo "  Health:  $URL/health"
echo "  Form:    $URL/"