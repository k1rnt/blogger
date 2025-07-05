#!/usr/bin/env bash
git diff --name-only "$1" "$2" | grep '^posts/.*\.md$' || true