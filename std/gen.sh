#!/bin/bash

std_dir="$(dirname "$0")"
repo_root="$(cd "$std_dir/.." && pwd)"

infer_package_name() {
	# Convert the filename $1 to a package name:
	# strip the .shk extension and dirname, turn dashes into underscores, convert '+' to 'x'.
	# cpp.shk keeps the Go package name cxx.
	name="$(printf '%s' "$(basename "$1")" | sed 's/\.shk$// ; s/-/_/g ; s/+/x/g')"
	if [ "$name" = "cpp" ]; then
		name="cxx"
	fi
	printf '%s' "$name"
}

gen_annotated_schema() {
	# Copy the schema from file "$2" into std/shark, and add Go annotations.
	# "$1" is the root of the input directory, used to trim the path prefix.
	indir="$1"
	infile="$2"

	base=$(echo -n "$infile" | sed -e "s@$indir@@" -e 's/\.capnp$/.shk/')
	outdir="$std_dir/shark"
	outfile="$outdir/$base"
	package_name="$(infer_package_name "$outfile")"
	[ -d "$(dirname "$outfile")" ] || mkdir -p "$(dirname "$outfile")"
	cat "$infile" - > "$outfile" << EOF
using Go = import "/go.shk";
\$Go.package("$package_name");
\$Go.import("github.com/kumo-rpc/sharkgo/std/shark/$package_name");
EOF
}

find_shk_files() {
	find "$1" -type f -name '*.shk'
}

# shark compile treats a bare language name as shark-<lang>. Pass the
# shark-go executable path so the plugin is this repo's generator.
shark_go_plugin() {
	if [ -n "${SHARK_GO:-}" ]; then
		printf '%s' "$SHARK_GO"
		return
	fi
	plugin="$(mktemp)"
	(cd "$repo_root" && go build -o "$plugin" ./shark-go) || return 1
	printf '%s' "$plugin"
}

gen_go_src() {
	# Generate go source code from the schema file $1. Create the package
	# directory if necessary. Output is <name>.sk.go.
	file="$1"
	filedir="$(dirname "$file")"
	package_name="$(infer_package_name "$file")"
	mkdir -p "$filedir/$package_name" && \
		shark compile --no-standard-import -I"$std_dir" -o"$shark_go":"$filedir/$package_name" --src-prefix="$filedir" "$file"
}

usage() {
	echo "Usage:"
	echo ""
	echo "    $0 import <path/to/shark/c++/src/shark>"
	echo "    $0 patch      # Apply necessary patches to schema"
	echo "    $0 compile    # Generate go source files"
	echo "    $0 clean-go   # Remove go source files"
	echo "    $0 clean-all  # Remove go source files and imported schemas"
}

# do_* implements the corresponding subcommand described in usage's output.
do_import() {
	input_dir="$1"
	for file in $(find "$input_dir" -type f -name '*.shk'); do
		gen_annotated_schema "$input_dir" "$file" || return 1
	done
}

do_patch() {
	cd "$std_dir" && patch -p1 < fixups.patch
}

do_compile() {
	shark_go="$(shark_go_plugin)" || return 1
	for file in $(find_shk_files "$std_dir"); do
		gen_go_src "$file" || return 1
	done
}

do_clean_go() {
	find "$std_dir" -name '*.sk.go' -delete
	find "$std_dir" -type d -empty -delete
}

do_clean_all() {
	do_clean_go
	find "$std_dir/shark" -name '*.shk' -delete
}

eq_or_usage() {
	# If "$1" is not equal to "$2", call usage and exit.
	if [ ! $1 = $2 ] ; then
		usage
		exit 1
	fi
}

case "$1" in
	import)    eq_or_usage $# 2; do_import "$2" ;;
	patch) eq_or_usage $# 1; do_patch ;;
	compile)   eq_or_usage $# 1; do_compile ;;
	clean-go)  eq_or_usage $# 1; do_clean_go ;;
	clean-all) eq_or_usage $# 1; do_clean_all ;;
	*) usage; exit 1 ;;
esac
