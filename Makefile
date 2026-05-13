.PHONY: build install test

# intervals.icu's spec is checked in at the repo root as openapi-spec.json (JSON
# already), so there's no yaml→json `spec` target as in sibling CLIs that
# embed yq-converted output. The cmd/openapi-spec.json mirror is what
# `//go:embed` consumes; keep both in sync if the upstream spec changes.

test:
	go test ./...

build:
	go build -o intervals .

# `go install` names the binary after the module directory ("intervals-icu-cli"),
# so we alias it to "intervals" to match the CLI's invocation name.
install: build
	go install .
	cp $$HOME/go/bin/intervals-icu-cli $$HOME/go/bin/intervals
