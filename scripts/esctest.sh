#!/usr/bin/env bash
#
# Runs the esctest2 conformance suite against darktile and compares the set of
# failing tests to a checked-in baseline. New failures fail the script; tests
# that started passing are reported so the baseline can be tightened.
#
#   scripts/esctest.sh            check against the baseline
#   scripts/esctest.sh --update   rewrite the baseline from this run
#
# esctest drives a terminal from the inside: it writes escape sequences to
# stdout and reads the replies from stdin. darktile-headless gives it a real
# pty attached to the terminal core, with no window or GPU involved.

set -euo pipefail

# The baseline is a sorted file that comm reads back, so the collation has to be
# the same wherever this runs. Under en_US the test names order differently than
# under C, and comm then reports the whole file as new failures.
export LC_ALL=C

cd "$(dirname "$0")/.."

BASELINE="internal/app/darktile/termutil/testdata/esctest-baseline.txt"
ESCTEST_DIR="${ESCTEST_DIR:-.cache/esctest2}"
ESCTEST_REPO="https://github.com/ThomasDickey/esctest2.git"
VT_LEVEL="${VT_LEVEL:-4}"
COLS="${COLS:-80}"
ROWS="${ROWS:-25}"

update=false
[[ "${1:-}" == "--update" ]] && update=true

if [[ ! -d "$ESCTEST_DIR" ]]; then
    echo "cloning esctest2 into $ESCTEST_DIR"
    git clone --depth 1 "$ESCTEST_REPO" "$ESCTEST_DIR"
fi

workdir=$(mktemp -d)
trap 'rm -rf "$workdir"' EXIT

binary="$workdir/darktile-headless"
log="$workdir/esctest.log"

go build -o "$binary" ./cmd/darktile-headless

echo "running esctest (vt level $VT_LEVEL, ${COLS}x${ROWS})"
"$binary" --cols "$COLS" --rows "$ROWS" --command \
    "cd '$(cd "$ESCTEST_DIR/esctest" && pwd)' && python3 esctest.py \
        --force --no-print-logs --max-vt-level $VT_LEVEL \
        --timeout 0.3 --logfile '$log'; exit" >/dev/null 2>"$workdir/stderr"

# A panic in the terminal kills the run partway through, which would otherwise
# look like a large batch of newly fixed tests.
if grep -q "^panic:" "$workdir/stderr"; then
    echo "darktile-headless panicked during the run:" >&2
    head -20 "$workdir/stderr" >&2
    exit 1
fi

grep -oP '^\*\*\* TEST \K\S+(?= FAILED)' "$log" | sort -u > "$workdir/failures.txt"

echo
grep -E '^\*\*\* [0-9]+ tests passed' "$log" | tail -1

if $update; then
    mkdir -p "$(dirname "$BASELINE")"
    cp "$workdir/failures.txt" "$BASELINE"
    echo "baseline updated: $(wc -l < "$BASELINE") failing tests"
    exit 0
fi

if [[ ! -f "$BASELINE" ]]; then
    echo "no baseline at $BASELINE — create one with: $0 --update" >&2
    exit 1
fi

new=$(comm -13 "$BASELINE" "$workdir/failures.txt")
fixed=$(comm -23 "$BASELINE" "$workdir/failures.txt")

if [[ -n "$fixed" ]]; then
    echo
    echo "tests that now pass (tighten the baseline with --update):"
    echo "$fixed" | sed 's/^/  + /'
fi

if [[ -n "$new" ]]; then
    echo
    echo "NEW FAILURES:" >&2
    echo "$new" | sed 's/^/  - /' >&2
    exit 1
fi

echo
echo "no new failures"
