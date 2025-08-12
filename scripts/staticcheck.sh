#!/bin/sh

script_dir="$(dirname "$0")"

cd "${script_dir}/../src"

go install honnef.co/go/tools/cmd/staticcheck@latest

if ! staticcheck ./...; then
  echo "❌ staticcheck failed" >&2
  exit 1
fi

echo "✅ staticcheck reported no bugs or performance issues"
