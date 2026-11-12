#!/bin/bash
# regen.sh - update shark-go and regenerate schemas
set -euo pipefail

cd "$(dirname "$0")"

echo "** go generate"
go generate

echo "** shark-go"
# Run tests so that we don't install a broken shark-go.
(cd shark-go && go generate && go test && go install)

echo "** schemas"
(cd std/shark; ../gen.sh compile)
(cd std/shark; shark compile --no-standard-import -I.. -o- schema.shk) | (cd internal/schema && shark-go -promises=0 -schemas=0 -structstrings=0)
(
	top="$PWD"
	cd flowcontrol/internal/test-tool
	shark compile --no-standard-import -I "$top/std" -ogo *.shk
)
go generate ./...
