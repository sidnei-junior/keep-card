#!/bin/sh
# Validates that the commit message header follows Conventional Commits.
# See CONTRIBUTING.md for the full convention.

max_length=72
types="feat|fix|docs|refactor|test|perf|build|ci|chore"
pattern="^($types)(\([a-z0-9-]+\))?!?: .+$"

header=$(grep -v '^#' "$1" | head -n 1 | tr -d '\r')
types_list=$(printf '%s' "$types" | sed 's/|/, /g')

case "$header" in
  "Merge "* | "Revert "* | "fixup! "* | "squash! "*)
    exit 0
    ;;
esac

if ! printf '%s' "$header" | grep -Eq "$pattern"; then
  cat >&2 <<EOF

Invalid commit message:

  $header

Expected Conventional Commits format:

  <type>(<optional scope>)<optional !>: <description>

Allowed types: $types_list

Examples:

  feat(card-service): add purchase registration endpoint
  fix: handle empty routing key
  docs: update local setup guide

EOF
  exit 1
fi

length=$(printf '%s' "$header" | wc -m | tr -d ' ')
if [ "$length" -gt "$max_length" ]; then
  echo "" >&2
  echo "Commit message header is $length characters long; the limit is $max_length." >&2
  echo "" >&2
  exit 1
fi
