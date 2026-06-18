#!/usr/bin/env bash
# End-to-end smoke test. Requires server running at FEEDBACK_API_URL (default localhost:8080).
set -euo pipefail

BASE="${FEEDBACK_API_URL:-http://localhost:8080}"
RCLI="${RCLI_BIN:-./bin/rcli}"
PASS=0
FAIL=0

pass() { echo "✓ $1"; PASS=$((PASS+1)); }
fail() { echo "✗ $1 — $2"; FAIL=$((FAIL+1)); }

echo "Testing against $BASE"
echo ""

# --- Web ---
curl -sf "$BASE/health" | grep -q '"status":"ok"' && pass "GET /health" || fail "GET /health" "down?"
curl -sf "$BASE/" | grep -q 'Submit Feedback' && pass "GET / (form)" || fail "GET /" "missing form"

# --- API submit ---
BUG=$(curl -sf -X POST "$BASE/api/v1/feedback" -H "Content-Type: application/json" \
  -d '{"app":"mamzo","type":"bug","title":"E2E smoke bug","description":"smoke test"}')
BUG_NUM=$(echo "$BUG" | python3 -c "import sys,json; print(json.load(sys.stdin)['number'])")
echo "$BUG" | grep -q '"app":"mamzo"' && pass "POST submit #$BUG_NUM" || fail "POST submit" "$BUG"

# --- API get (by number, no index lag) ---
GET=$(curl -sf "$BASE/api/v1/feedback/$BUG_NUM")
echo "$GET" | grep -q 'E2E smoke bug' && pass "GET issue #$BUG_NUM" || fail "GET issue" "$GET"

# --- API list (retry — GitHub label index can lag a few seconds) ---
LIST=""
for i in 1 2 3 4 5; do
  LIST=$(curl -sf "$BASE/api/v1/feedback?app=mamzo&state=open" || echo "null")
  echo "$LIST" | grep -q 'E2E smoke bug' && break
  sleep 2
done
echo "$LIST" | grep -q 'E2E smoke bug' && pass "GET list app=mamzo" || fail "GET list" "$LIST"

# --- CLI (before status change) ---
if [ -x "$RCLI" ]; then
  $RCLI feedback view "$BUG_NUM" | grep -q 'E2E smoke bug' && pass "CLI view" || fail "CLI view" ""
  $RCLI feedback list --app mamzo | grep -q 'E2E smoke bug' && pass "CLI list" || fail "CLI list" ""
else
  fail "CLI" "binary not found at $RCLI"
fi

# --- API key endpoints ---
curl -sf -X POST "$BASE/api/v1/feedback/$BUG_NUM/comments" \
  -H "Content-Type: application/json" -H "X-API-Key: ${FEEDBACK_API_KEY:-local-dev-key}" -H "X-App: mamzo" \
  -d '{"body":"smoke comment"}' | grep -q 'comment added' && pass "POST comment" || fail "POST comment" ""

UPD=$(curl -sf -X PATCH "$BASE/api/v1/feedback/$BUG_NUM" \
  -H "Content-Type: application/json" -H "X-API-Key: ${FEEDBACK_API_KEY:-local-dev-key}" -H "X-App: mamzo" \
  -d '{"status":"triaged"}')
echo "$UPD" | grep -q 'triaged' && pass "PATCH status" || fail "PATCH status" "$UPD"

# --- Auth ---
HTTP=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE/api/v1/feedback/$BUG_NUM/comments" \
  -H "Content-Type: application/json" -d '{"body":"no key"}')
[ "$HTTP" = "401" ] && pass "Reject missing API key" || fail "Auth" "HTTP $HTTP"

echo ""
echo "=== $PASS passed, $FAIL failed ==="
[ "$FAIL" -eq 0 ]