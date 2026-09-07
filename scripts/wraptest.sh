#!/usr/bin/env bash
#
# Runs mattiase/wraptest against darktile and diffs the report against a
# checked-in baseline.
#
#   scripts/wraptest.sh            check against the baseline
#   scripts/wraptest.sh --update   rewrite the baseline from this run
#
# wraptest probes the hidden Last Column Flag: the deferred-wrap state a VT
# enters after printing into the rightmost column, and which operations clear
# it. esctest barely covers this, and it is the part of the spec emulators most
# often invent for themselves, so it gets its own baseline.
#
# The program talks to the terminal on stderr and prints its report on stdout,
# so only stdout is redirected here.

set -euo pipefail

cd "$(dirname "$0")/.."

BASELINE="internal/app/darktile/termutil/testdata/wraptest-baseline.txt"
WRAPTEST_DIR="${WRAPTEST_DIR:-.cache/wraptest}"
WRAPTEST_REPO="https://github.com/mattiase/wraptest.git"
COLS="${COLS:-80}"
ROWS="${ROWS:-25}"
TIMEOUT="${TIMEOUT:-60}"

update=false
[[ "${1:-}" == "--update" ]] && update=true

if [[ ! -d "$WRAPTEST_DIR" ]]; then
    echo "cloning wraptest into $WRAPTEST_DIR"
    git clone --depth 1 "$WRAPTEST_REPO" "$WRAPTEST_DIR"
fi

workdir=$(mktemp -d)
trap 'rm -rf "$workdir"' EXIT

report="$workdir/report.txt"

cc -o "$workdir/wraptest" "$WRAPTEST_DIR/wraptest.c"
go build -o "$workdir/darktile-headless" ./cmd/darktile-headless

# wraptest puts the tty in raw mode, which a background process group is not
# allowed to do: it would be stopped by SIGTTOU and produce nothing. So the
# run cannot be wrapped in `timeout` here; the bound goes on the terminal
# itself, which is what holds the pty open.
cat > "$workdir/run.sh" <<EOF
#!/bin/sh
"$workdir/wraptest" > "$report"
EOF
chmod +x "$workdir/run.sh"

echo "running wraptest (${COLS}x${ROWS})"
if ! timeout "$TIMEOUT" "$workdir/darktile-headless" \
        --cols "$COLS" --rows "$ROWS" \
        --command "$workdir/run.sh; exit" >/dev/null 2>&1; then
    echo "wraptest did not finish within ${TIMEOUT}s" >&2
    exit 1
fi

if [[ ! -s "$report" ]]; then
    echo "wraptest produced no report" >&2
    exit 1
fi

if $update; then
    mkdir -p "$(dirname "$BASELINE")"
    cp "$report" "$BASELINE"
    echo "baseline updated"
    exit 0
fi

if [[ ! -f "$BASELINE" ]]; then
    echo "no baseline at $BASELINE — create one with: $0 --update" >&2
    exit 1
fi

if ! diff -u "$BASELINE" "$report"; then
    echo >&2
    echo "wrapping behaviour changed; re-baseline with: $0 --update" >&2
    exit 1
fi

echo
echo "wrapping behaviour matches the baseline"
