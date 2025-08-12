#!/bin/sh

script_dir="$(dirname "$0")"

cd "${script_dir}/../src"

go install golang.org/x/vuln/cmd/govulncheck@latest

if ! govulncheck ./...; then
  echo "❌ govulncheck failed" >&2
  exit 1
fi

echo "✅ govulncheck reported no vulnerabilities"
