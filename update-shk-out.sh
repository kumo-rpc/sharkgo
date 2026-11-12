#!/bin/bash
# Manually refresh CodeGeneratorRequest fixtures (*.shk.out).
# These files are read by tests and are not produced by go generate.
set -euo pipefail

cd "$(dirname "$0")"
std="$PWD/std"

compile() {
	out="$1"
	shift
	shark compile --no-standard-import -I "$std" -o- "$@" > "$out"
	echo "wrote $out"
}

(
	cd shark-go/testdata
	compile go.shk.out go.shk
	compile const.shk.out const.shk
	compile group.shk.out group.shk
	compile scopes.shk.out scopes.shk
	compile util.shk.out util.shk
	compile persistent-simple-and-samepkg.shk.out persistent-simple.shk persistent-samepkg.shk
)

compile shark-go/testdata/aircraft.shk.out internal/aircraftlib/aircraft.shk
compile shark-go/testdata/rpc.shk.out std/shark/rpc.shk
compile encoding/text/testdata/txt.shk.out encoding/text/testdata/txt.shk
