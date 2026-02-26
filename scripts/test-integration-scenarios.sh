#!/usr/bin/env bash
set -euo pipefail

API_BASE_URL="${API_BASE_URL:-http://localhost:8080}"
WEB_BASE_URL="${WEB_BASE_URL:-http://localhost:5173}"
MAX_ATTEMPTS="${MAX_ATTEMPTS:-5}"
RETRY_DELAY_SECONDS="${RETRY_DELAY_SECONDS:-2}"

if ! command -v curl >/dev/null 2>&1; then
  echo "ERROR: curl is required but was not found in PATH"
  exit 1
fi

check_http() {
  local name="$1"
  local url="$2"
  local expected_status="$3"
  local body_must_contain="${4:-}"

  local tmp status
  tmp="$(mktemp)"
  status="$(curl -sS -L --max-time 10 -o "$tmp" -w "%{http_code}" "$url" || true)"

  if [[ "$status" != "$expected_status" ]]; then
    echo "FAIL: $name - expected HTTP $expected_status, got ${status:-curl_error} ($url)"
    rm -f "$tmp"
    return 1
  fi

  if [[ -n "$body_must_contain" ]] && ! grep -qi "$body_must_contain" "$tmp"; then
    echo "FAIL: $name - response body does not contain '$body_must_contain' ($url)"
    rm -f "$tmp"
    return 1
  fi

  echo "PASS: $name ($url)"
  rm -f "$tmp"
  return 0
}

retry_check() {
  local name="$1"
  local url="$2"
  local expected_status="$3"
  local body_must_contain="${4:-}"
  local attempt=1

  while [[ "$attempt" -le "$MAX_ATTEMPTS" ]]; do
    if check_http "$name" "$url" "$expected_status" "$body_must_contain"; then
      return 0
    fi

    if [[ "$attempt" -lt "$MAX_ATTEMPTS" ]]; then
      echo "Retrying '$name' in ${RETRY_DELAY_SECONDS}s (attempt ${attempt}/${MAX_ATTEMPTS})..."
      sleep "$RETRY_DELAY_SECONDS"
    fi
    attempt=$((attempt + 1))
  done

  return 1
}

failures=0

run_scenario() {
  local name="$1"
  local url="$2"
  local expected_status="$3"
  local body_must_contain="${4:-}"

  if ! retry_check "$name" "$url" "$expected_status" "$body_must_contain"; then
    failures=$((failures + 1))
  fi
}

echo "Running integration scenario checks..."
echo "API_BASE_URL=$API_BASE_URL"
echo "WEB_BASE_URL=$WEB_BASE_URL"

run_scenario "Backend health" "${API_BASE_URL}/api/v1/healthz" "200" "status"
run_scenario "Backend root health" "${API_BASE_URL}/health" "200" "status"
run_scenario "Admin homepage" "${WEB_BASE_URL}" "200" "<html"

if [[ "$failures" -gt 0 ]]; then
  echo "Integration scenarios failed: $failures"
  exit 1
fi

echo "All integration scenarios passed."
