#!/usr/bin/env bash
# Creates labels in the feedback repo. Requires: gh CLI authenticated.
set -euo pipefail

OWNER="${GITHUB_OWNER:?Set GITHUB_OWNER}"
REPO="${GITHUB_REPO:-rcli-feedback-app}"

labels=(
  "source:rclip|0E8A16|Feedback from r:clip"
  "source:boka|1D76DB|Feedback from BoKa"
  "source:thxbud|FBCA04|Feedback from ThxBud"
  "source:mamzo|E99695|Feedback from MamZo"
  "source:glasscourt|5319E7|Feedback from GlassCourt"
  "type:bug|D73A4A|Bug report"
  "type:feature-request|A2EEEF|Feature request"
  "status:triaged|C5DEF5|Triaged"
  "status:in-progress|FEF2C0|In progress"
  "status:done|0E8A16|Done"
)

for entry in "${labels[@]}"; do
  IFS='|' read -r name color description <<< "$entry"
  gh label create "$name" --repo "$OWNER/$REPO" --color "$color" --description "$description" --force
  echo "✓ $name"
done

echo ""
echo "Labels created. Create a GitHub Project board manually:"
echo "  Columns: Inbox → Triaged → In Progress → Done"