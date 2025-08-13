#!/bin/sh

script_dir="$(dirname "$0")"

cd "${script_dir}/../src/internal"

if ! go test -v ./...; then
  exit_code=$?
  echo "❌ Tests failed" >&2
  exit $exit_code
fi

echo "✅ All tests passed"