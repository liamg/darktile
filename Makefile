

default: build

build:
	./scripts/build.sh

test:
	go test ./...

# A flake only ever sees a commit, never a tag, so VERSION is the source of
# truth and the tag is derived from it rather than the other way round. Bump
# VERSION, commit, then run this.
tag:
	@v=$$(cat VERSION); \
	git diff --quiet HEAD -- VERSION || { echo "VERSION differs from HEAD; commit it first" >&2; exit 1; }; \
	git tag -a "v$$v" -m "v$$v" && echo "tagged v$$v"

# Known CVEs in the dependency tree, filtered down to the ones actually
# reachable from this code.
vuln:
	govulncheck ./...

# Escape sequence conformance against esctest2, ratcheted by a checked-in
# baseline of known failures.
esctest:
	./scripts/esctest.sh

esctest-update:
	./scripts/esctest.sh --update

# Deferred-wrap (Last Column Flag) semantics, which esctest barely covers.
wraptest:
	./scripts/wraptest.sh

wraptest-update:
	./scripts/wraptest.sh --update

# Which Unicode version the character widths agree with. Stops at its entry
# probe today because every rune is one column wide; see the script.
ucs-detect:
	./scripts/ucs-detect.sh

# The only coverage of the renderer. Needs a display, so it runs under Xvfb
# rather than as a go test.
render:
	./scripts/render-golden.sh

render-update:
	./scripts/render-golden.sh --update

# The parser is an untrusted-input boundary; this looks for panics and hangs.
fuzz:
	go test ./internal/app/darktile/termutil -run '^$$' -fuzz FuzzParser -fuzztime 300s

# Write results to a file and compare two runs with benchstat before and after
# a change.
bench:
	go test ./internal/app/darktile/termutil -run '^$$' -bench Parser -benchmem -count=8

# Where the parse path actually spends its time. WORKLOAD names one of the
# benchmarks; the profiles stay in .cache so they can be opened again with
# `go tool pprof -http=: .cache/termutil.test .cache/cpu.prof`.
WORKLOAD ?= sgr-churn
profile:
	@mkdir -p .cache
	go test ./internal/app/darktile/termutil -run '^$$' \
		-bench 'Parser/$(WORKLOAD)$$$$' -benchtime 5s \
		-cpuprofile .cache/cpu.prof -memprofile .cache/mem.prof \
		-o .cache/termutil.test
	@echo
	go tool pprof -top -nodecount=20 .cache/termutil.test .cache/cpu.prof

.PHONY: default build test tag vuln esctest esctest-update wraptest wraptest-update \
	ucs-detect render render-update fuzz bench profile
