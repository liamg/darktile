#!/usr/bin/env bash
#
# Runs jquast/ucs-detect against darktile to find which Unicode version the
# terminal's character widths agree with.
#
#   scripts/ucs-detect.sh          run and print a summary
#   scripts/ucs-detect.sh --full   no limits; slow, but the real answer
#
# darktile currently advances the cursor by exactly one column per rune, so
# ucs-detect stops at its entry probe: it writes U+231A WATCH, measures the
# cursor, gets 1 instead of 2, and skips every test. The JSON it saves still
# records that, and the run becomes useful the moment character widths are
# implemented. The cheap regression guard until then is
# TestDeviationEveryRuneOccupiesOneColumn in termutil.
#
# ucs-detect is not packaged in nixpkgs, so it goes in a venv under .cache.

set -euo pipefail

cd "$(dirname "$0")/.."

VENV="${UCS_VENV:-.cache/ucs-detect-venv}"
REPORT="${UCS_REPORT:-.cache/ucs-detect.json}"
COLS="${COLS:-80}"
ROWS="${ROWS:-25}"
TIMEOUT="${TIMEOUT:-900}"

limits=(--limit-codepoints 64 --limit-graphemes 8 --limit-category-time 10 --no-languages-test)
[[ "${1:-}" == "--full" ]] && limits=()

if [[ ! -x "$VENV/bin/ucs-detect" ]]; then
    echo "installing ucs-detect into $VENV"
    python3 -m venv "$VENV"
    "$VENV/bin/pip" install --quiet --upgrade pip
    "$VENV/bin/pip" install --quiet ucs-detect
fi

workdir=$(mktemp -d)
trap 'rm -rf "$workdir"' EXIT

go build -o "$workdir/darktile-headless" ./cmd/darktile-headless

mkdir -p "$(dirname "$REPORT")"
rm -f "$REPORT"

# ucs-detect drives the terminal on stderr and prints its report on stdout, so
# only stdout is redirected. It also puts the tty in cbreak mode, which a
# background process group may not do, so it cannot be wrapped in `timeout`
# here — the bound goes on the terminal holding the pty instead.
cat > "$workdir/run.sh" <<EOF
#!/bin/sh
"$(cd "$(dirname "$VENV")" && pwd)/$(basename "$VENV")/bin/ucs-detect" \\
    ${limits[*]} --save-json "$(pwd)/$REPORT" > "$workdir/stdout.txt"
EOF
chmod +x "$workdir/run.sh"

echo "running ucs-detect (${COLS}x${ROWS})"
if ! timeout "$TIMEOUT" "$workdir/darktile-headless" \
        --cols "$COLS" --rows "$ROWS" \
        --command "$workdir/run.sh; exit" >/dev/null 2>&1; then
    echo "ucs-detect did not finish within ${TIMEOUT}s" >&2
    exit 1
fi

if [[ ! -s "$REPORT" ]]; then
    echo "ucs-detect produced no report" >&2
    exit 1
fi

python3 - "$REPORT" <<'PY'
import json, sys

report = json.load(open(sys.argv[1]))
results = report["test_results"]

print()
print(f"report: {sys.argv[1]}")
print(f"ambiguous width: {report['ambiguous_width']} (-1 = not determined)")

tested = {name: value for name, value in results.items() if value}
if not tested:
    print()
    print("no category ran: the terminal failed the wide-character entry probe,")
    print("so ucs-detect had nothing to measure. Character widths are the gap.")
    sys.exit(0)

for name, value in sorted(tested.items()):
    if isinstance(value, dict):
        best = max(value, key=lambda version: value[version].get("pct_success", 0))
        print(f"{name}: best {best} at {value[best].get('pct_success')}%")
PY
