#!/bin/sh
# GitHub Action entrypoint. Prints a proof receipt and, when a token and PR
# number are present, posts it as a comment.
set -eu

cd "${GITHUB_WORKSPACE:-.}"

args="verify --repo . --base ${INPUT_BASE:-main} --head ${INPUT_HEAD:-HEAD} --reruns ${INPUT_RERUNS:-3} --max-mutants ${INPUT_MAX_MUTANTS:-12}"

body_file=""
if [ -n "${INPUT_PR_BODY:-}" ]; then
  body_file="$(mktemp)"
  printf '%s\n' "$INPUT_PR_BODY" > "$body_file"
  args="$args --pr-body $body_file"
fi

if [ "${INPUT_LLM_SUMMARY:-false}" = "true" ]; then
  args="$args --llm-summary"
fi

receipt="$(mktemp)"
set +e
# shellcheck disable=SC2086
proofline $args > "$receipt"
code=$?
set -e

cat "$receipt"

verdict="$(awk -F': ' '/^\*\*Verdict:/{print $2; exit}' "$receipt" | tr -d '*')"
if [ -n "${GITHUB_OUTPUT:-}" ]; then
  echo "verdict=${verdict}" >> "$GITHUB_OUTPUT"
fi

if [ "${INPUT_COMMENT:-true}" != "true" ]; then
  exit "$code"
fi

token="${INPUT_TOKEN:-${GITHUB_TOKEN:-}}"
pr="${INPUT_PR:-}"
if [ -z "$pr" ] && [ -n "${GITHUB_EVENT_PATH:-}" ] && [ -f "$GITHUB_EVENT_PATH" ]; then
  pr="$(python3 -c 'import json,os; d=json.load(open(os.environ["GITHUB_EVENT_PATH"])); print((d.get("pull_request") or {}).get("number",""))')"
fi

if [ -z "$token" ] || [ -z "$pr" ] || [ -z "${GITHUB_REPOSITORY:-}" ]; then
  echo "proofline: skipping comment (need token, pull request number, and GITHUB_REPOSITORY)" >&2
  exit "$code"
fi

payload="$(python3 -c 'import json,sys; print(json.dumps({"body": sys.stdin.read()}))' < "$receipt")"
curl -sS -o /tmp/proofline-comment.json -w "%{http_code}" \
  -X POST \
  -H "Authorization: Bearer ${token}" \
  -H "Accept: application/vnd.github+json" \
  -H "Content-Type: application/json" \
  "https://api.github.com/repos/${GITHUB_REPOSITORY}/issues/${pr}/comments" \
  -d "$payload" > /tmp/proofline-comment.code || true

http="$(cat /tmp/proofline-comment.code 2>/dev/null || echo 0)"
case "$http" in
  2*) ;;
  *)
    echo "proofline: comment failed with HTTP ${http}" >&2
    cat /tmp/proofline-comment.json >&2 || true
    if [ "$code" -eq 0 ]; then
      exit 1
    fi
    ;;
esac

exit "$code"
