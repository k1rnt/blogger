#!/usr/bin/env bash
before="$1"
after="$2"
set -euo pipefail

if [ -z "$before" ] || ! git rev-parse -q --verify "$before^{commit}" >/dev/null; then
  before=$(git rev-parse "$after^" 2>/dev/null || echo "$after")
fi

git diff --name-only "$before" "$after" -- 'posts/**/*.md' || true
