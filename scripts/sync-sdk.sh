#!/usr/bin/env bash
# sync-sdk.sh — copy this repo's sdk/ into a downstream contract repo, verbatim.
#
# go-contract-template is the PRIMARY copy of the Magi contract SDK. Every
# contract repo (magi_nft-contract, magi_token-contract, ...) carries a
# byte-identical sdk/ directory, so a host function added here reaches all of
# them by running this script — never by hand-editing a downstream copy.
#
# The sdk package deliberately imports nothing from its own module (the
# custom-GC runtime is blank-imported by the contract's main package instead),
# which is what makes a verbatim copy possible.
#
# Usage:
#   scripts/sync-sdk.sh /path/to/magi_nft-contract        # copy, then verify
#   scripts/sync-sdk.sh --check /path/to/magi_nft-contract # verify only (CI)
set -euo pipefail

CHECK_ONLY=0
if [[ "${1:-}" == "--check" ]]; then
	CHECK_ONLY=1
	shift
fi

TARGET="${1:-}"
if [[ -z "$TARGET" ]]; then
	echo "usage: $0 [--check] <path-to-contract-repo>" >&2
	exit 2
fi

SRC="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/sdk"
DST="$TARGET/sdk"

if [[ ! -d "$DST" ]]; then
	echo "error: $DST does not exist (is $TARGET a contract repo?)" >&2
	exit 2
fi

if [[ "$CHECK_ONLY" -eq 0 ]]; then
	rm -f "$DST"/*.go
	cp "$SRC"/*.go "$DST"/
fi

if diff -r "$SRC" "$DST" >/dev/null; then
	echo "OK: $DST is identical to the template sdk/"
else
	echo "DRIFT: $DST differs from the template sdk/" >&2
	diff -r "$SRC" "$DST" >&2 || true
	exit 1
fi

if [[ "$CHECK_ONLY" -eq 0 ]]; then
	cat <<'NOTE'

Next steps in the target repo:
  1. rebuild the wasm (see its readme; tinygo needs a go <= 1.25 toolchain,
     e.g. GOTOOLCHAIN=go1.25.3)
  2. go test ./test/...

If the target's contract calls a host function its pinned vsc-node does not
provide (see the `replace vsc-node` line in its go.mod), the module will fail
to instantiate in that repo's tests — bump the pin in that case.
NOTE
fi
