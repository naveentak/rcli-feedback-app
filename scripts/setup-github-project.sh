#!/usr/bin/env bash
# Creates Feedback Kanban project and links it to rcli-feedback-app.
# Requires: gh auth refresh -h github.com -s project,read:project
set -euo pipefail

OWNER="${GITHUB_OWNER:-naveentak}"
REPO="${GITHUB_REPO:-rcli-feedback-app}"

echo "Checking gh project scope..."
if ! gh auth status 2>&1 | grep -q 'project'; then
  echo "Run this first (interactive):"
  echo "  gh auth refresh -h github.com -s project,read:project"
  exit 1
fi

PROJECT=$(gh project create --owner "$OWNER" --title "Feedback Kanban" --format json)
PROJECT_NUM=$(echo "$PROJECT" | python3 -c "import sys,json; print(json.load(sys.stdin)['number'])")
PROJECT_ID=$(echo "$PROJECT" | python3 -c "import sys,json; print(json.load(sys.stdin)['id'])")

echo "→ Linking to $OWNER/$REPO..."
gh project link "$PROJECT_NUM" --owner "$OWNER" --repo "$OWNER/$REPO"

echo "→ Adding Status field options..."
# GitHub Projects v2 uses built-in Status field — configure columns in the UI:
# Inbox (default) → Triaged → In Progress → Done
echo ""
echo "✓ Project created: https://github.com/users/$OWNER/projects/$PROJECT_NUM"
echo ""
echo "Manual step: open the project and set Status columns to:"
echo "  Inbox → Triaged → In Progress → Done"
echo ""
echo "Add existing issues:"
gh issue list --repo "$OWNER/$REPO" --limit 50 --json number --jq '.[].number' | while read -r num; do
  gh project item-add "$PROJECT_NUM" --owner "$OWNER" --url "https://github.com/$OWNER/$REPO/issues/$num" 2>/dev/null || true
done
echo "Done."