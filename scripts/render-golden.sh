#!/usr/bin/env bash
#
# Renders a set of scenes in a real darktile window under Xvfb and compares the
# screenshots to checked-in goldens.
#
#   scripts/render-golden.sh            check against the goldens
#   scripts/render-golden.sh --update   rewrite the goldens from this run
#   scripts/render-golden.sh text-attributes 256-colour   only these scenes
#
# The parser has unit tests, esctest and fuzzing behind it; the renderer has
# nothing. Everything from the buffer to the pixels — glyph rasterisation,
# colour resolution, sixel blitting, ligatures — is only covered here.
#
# Differences are reported as a pixel count with a diff image written next to
# the golden, because a change here is usually meant to be looked at rather
# than read.

set -euo pipefail

cd "$(dirname "$0")/.."

GOLDEN_DIR="internal/app/darktile/gui/testdata/render"
OUT_DIR="${OUT_DIR:-.cache/render}"
SCREENSHOT_MS="${SCREENSHOT_MS:-2500}"
# Software rendering is deterministic for a given mesa, which the flake pins.
# A non-zero budget only exists so a mesa bump does not have to be a re-baseline.
TOLERANCE="${TOLERANCE:-0}"

update=false
if [[ "${1:-}" == "--update" ]]; then
    update=true
    shift
fi

# A scene runs as the terminal's shell rather than as a command typed into one,
# so no prompt, echoed command line or shell version string ends up in the
# image. Each one hides the cursor first, since whether it is drawn depends on
# the window having focus.
scene_body() {
    case "$1" in
    text-attributes)
        cat <<'EOF'
printf '\033[1mbold\033[0m \033[2mdim\033[0m \033[3mitalic\033[0m \033[4munderline\033[0m\n'
printf '\033[7mreverse\033[0m \033[9mstrikethrough\033[0m \033[5mblink\033[0m\n\n'
for i in 0 1 2 3 4 5 6 7; do printf "\033[3${i}m██ normal ${i}\033[0m\n"; done
for i in 0 1 2 3 4 5 6 7; do printf "\033[9${i}m██ bright ${i}\033[0m\n"; done
EOF
        ;;
    256-colour)
        cat <<'EOF'
i=0
while [ $i -lt 256 ]; do
  printf "\033[48;5;${i}m  \033[0m"
  i=$((i + 1))
  [ $((i % 32)) -eq 0 ] && printf '\n'
done
printf '\n'
i=16
while [ $i -lt 256 ]; do
  printf "\033[38;5;${i}m▓\033[0m"
  i=$((i + 1))
done
printf '\n'
EOF
        ;;
    true-colour)
        cat <<'EOF'
r=0
while [ $r -lt 8 ]; do
  c=0
  while [ $c -lt 64 ]; do
    printf "\033[48;2;$((c * 4));$((r * 32));$((255 - c * 4))m "
    c=$((c + 1))
  done
  printf '\033[0m\n'
  r=$((r + 1))
done
EOF
        ;;
    box-drawing)
        cat <<'EOF'
printf '\033(0lqqqwqqqk\033(B\n'
printf '\033(0x   x   x\033(B\n'
printf '\033(0tqqqnqqqu\033(B\n'
printf '\033(0x   x   x\033(B\n'
printf '\033(0mqqqvqqqj\033(B\n'
printf '\033(0`afgijklmnopqrstuvwxyz{|}~\033(B\n'
EOF
        ;;
    sixel)
        cat <<'EOF'
printf '\033Pq'
printf '#0;2;90;10;10#1;2;10;90;20#2;2;20;30;95'
printf '#0~~~~~~~~$#1??????~~~~~~~~$#2????????????~~~~~~~~-'
printf '#2~~~~~~~~$#0??????~~~~~~~~$#1????????????~~~~~~~~'
printf '\033\\\n'
printf 'sixel above\n'
EOF
        ;;
    scroll-region)
        cat <<'EOF'
i=1
while [ $i -le 12 ]; do printf "line $i\n"; i=$((i + 1)); done
printf '\033[3;8r\033[8;1Hscrolled\n\n\n'
printf '\033[r\033[12;1Hafter reset\n'
EOF
        ;;
    *)
        echo "unknown scene: $1" >&2
        return 1
        ;;
    esac
}

ALL_SCENES=(text-attributes 256-colour true-colour box-drawing sixel scroll-region)
scenes=("${@:-}")
[[ -z "${scenes[0]:-}" ]] && scenes=("${ALL_SCENES[@]}")

workdir=$(mktemp -d)
trap 'rm -rf "$workdir"' EXIT

mkdir -p "$OUT_DIR" "$GOLDEN_DIR"

# Two runs at once write the same screenshot paths and delete each other's
# output, which shows up as "no screenshot produced" rather than as a diff.
exec 9>"$OUT_DIR/.lock"
if ! flock -n 9; then
    echo "waiting for another render run to finish..." >&2
    flock 9
fi

binary="$workdir/darktile"
go build -o "$binary" ./cmd/darktile

# A config or theme in the real home would change what is drawn. HOME moves
# too, so LD_LIBRARY_PATH has to survive it: ebiten dlopens libGL, which the
# devShell puts there.
export XDG_CONFIG_HOME="$workdir/config"
export HOME="$workdir/home"
export ENV=/dev/null
mkdir -p "$XDG_CONFIG_HOME" "$HOME"

# The goldens are compared pixel for pixel, so the renderer has to be the same
# everywhere: llvmpipe is, a GPU driver is not. MESA_PREFIX comes from the
# devShell and is what makes this work off NixOS, where there is no
# /run/opengl-driver for libglvnd to load drivers from.
#
# This goes on darktile's own command line rather than in the environment
# because Xvfb must not inherit it -- it links the system mesa, and mixing the
# two segfaults the X server before darktile ever starts.
gl_env=(env LIBGL_ALWAYS_SOFTWARE=1)
if [[ -n "${MESA_PREFIX:-}" ]]; then
    gl_env+=(
        "LIBGL_DRIVERS_PATH=$MESA_PREFIX/lib/dri"
        "__EGL_VENDOR_LIBRARY_DIRS=$MESA_PREFIX/share/glvnd/egl_vendor.d"
    )
fi

failed=0
changed=()

for scene in "${scenes[@]}"; do
    script="$workdir/$scene.sh"
    {
        echo '#!/bin/sh'
        echo "printf '\033[?25l'"
        # The window is created at a default size and resized once the real one
        # is known, which reflows the buffer. A scene that prints into that gap
        # races the resize and lands on a different scroll offset, so wait it
        # out before emitting anything.
        echo "sleep 1"
        scene_body "$scene"
        # Outlives the screenshot, then lets the terminal close on its own.
        echo "sleep $((SCREENSHOT_MS / 1000 + 3))"
    } > "$script"
    chmod +x "$script"

    shot="$OUT_DIR/$scene.png"
    rm -f "$shot"

    xvfb-run -a -s "-screen 0 1024x768x24" \
        "${gl_env[@]}" "$binary" --shell "$script" \
        --screenshot-after-ms "$SCREENSHOT_MS" \
        --screenshot-filename "$shot" </dev/null >"$workdir/$scene.log" 2>&1 || true

    if [[ ! -s "$shot" ]]; then
        echo "FAIL $scene: no screenshot produced" >&2
        tail -5 "$workdir/$scene.log" >&2
        failed=1
        continue
    fi

    golden="$GOLDEN_DIR/$scene.png"

    if $update; then
        cp "$shot" "$golden"
        echo "updated $scene"
        continue
    fi

    if [[ ! -f "$golden" ]]; then
        echo "FAIL $scene: no golden — create one with: $0 --update" >&2
        failed=1
        continue
    fi

    diff_image="$OUT_DIR/$scene.diff.png"
    # Not `compare -metric AE`: under an HDRI build that returns a fractional
    # error, so a handful of changed pixels reads as 0.93 and any truncation to
    # an integer silently becomes "identical". This counts pixels where some
    # channel differs at all, as an integer.
    pixels=$(magick "$golden" "$shot" -compose difference -composite \
        -separate -evaluate-sequence max -threshold 0 \
        -format '%[fx:int(mean*w*h+0.5)]' info: 2>/dev/null || true)
    compare -metric AE "$golden" "$shot" "$diff_image" >/dev/null 2>&1 || true

    if [[ ! "$pixels" =~ ^[0-9]+$ ]]; then
        echo "FAIL $scene: could not compare (image sizes differ?)" >&2
        failed=1
        continue
    fi

    if (( pixels > TOLERANCE )); then
        echo "FAIL $scene: $pixels pixels differ (tolerance $TOLERANCE)"
        echo "     golden $golden"
        echo "     actual $shot"
        echo "     diff   $diff_image"
        changed+=("$scene")
        failed=1
    else
        echo "ok   $scene"
        rm -f "$diff_image"
    fi
done

if $update; then
    echo
    echo "goldens written to $GOLDEN_DIR"
    exit 0
fi

if (( failed )); then
    echo
    if (( ${#changed[@]} )); then
        echo "render changed in: ${changed[*]}" >&2
        echo "look at the diff images, then re-baseline with: $0 --update" >&2
    fi
    exit 1
fi

echo
echo "render matches the goldens"
