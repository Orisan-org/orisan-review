PARKED v0.1

# Orisan Review

Orisan Review is a local-first CLI that analyzes pull request diffs for security-sensitive AI-assisted output risk.

> **Status: alpha, parked for maintenance.** Review is a secondary community
> artifact. Maintenance scope is portability, CI health, security fixes,
> documentation honesty, and keeping deterministic diff-routing workflows
> working.

It helps answer:

> What changed, why does it matter, and who needs to review this before it is merged or trusted?

Review v0.1 is diff-aware risk routing, not generic SAST.

## What Review Is

Review is a portfolio/community tool for validating AI-assisted output risk before trust. It is designed to run locally or in CI, inspect git diffs or patch files, classify security-sensitive changes, and route attention to the right reviewer.

## What Review Is Not

Review v0.1 is not SAST, SCA, DAST, or an AI code reviewer. It does not prove that a vulnerability exists, does not replace human review, does not upload source code, does not call an LLM, and does not post PR comments by default.

## Current Status

This repository is an early v0.1 implementation. The `analyze` workflow can parse PR-like diffs, run deterministic routing rules, emit text/JSON/Markdown reports, and validate behavior against a fixture corpus. Rule coverage is intentionally narrow and focused on review routing, not vulnerability proof.

## Prerequisites

- Go 1.20 or newer to build and test from source.
- Git for `git diff` examples and direct `--repo --base --head` analysis.

## Build

```sh
go fmt ./...
go vet ./...
go test ./...
mkdir -p bin
go build -o bin/orisan-review ./cmd/orisan-review
```

## Quickstart

```sh
./bin/orisan-review --help
./bin/orisan-review version
./bin/orisan-review list-rules
./bin/orisan-review list-categories
./bin/orisan-review analyze --patch testdata/diffs/tls_verification_disabled.patch --format text
./bin/orisan-review analyze --patch testdata/diffs/tls_verification_disabled.patch --format json
```

Analyze a PR-like diff from a patch file:

```sh
git diff main...HEAD > /tmp/review.patch
./bin/orisan-review analyze --patch /tmp/review.patch --format text
./bin/orisan-review analyze --patch /tmp/review.patch --format json --out review-report.json
```

Generate a local HTML report:

```sh
./bin/orisan-review analyze --patch /tmp/review.patch --format html --out review-report.html
open review-report.html
```

Analyze a diff from stdin:

```sh
git diff main...HEAD | ./bin/orisan-review analyze --stdin --format text
```

Analyze a git ref range directly:

```sh
./bin/orisan-review analyze --repo . --base main --head HEAD --format text
```

## Privacy

Review is local-first by design. v0.1 must not upload source code, diffs, findings, or reports to Orisan. Findings must store short redacted evidence only and always set `payload_stored=false`.

## Validation

v0.1 success is measured by fixture routing, not vulnerability count. The release gate is:

```sh
go test ./...
test -z "$(gofmt -l .)"
go vet ./...
go build ./cmd/orisan-review
```

The validation corpus lives under `testdata/diffs` with semantic expectations under `testdata/expected`. Tests check reviewer routing, finding categories, `payload_stored=false`, redacted evidence, and raw secret leakage.

## GitHub Integration

The recommended v0.1 integration is a GitHub Actions workflow that runs Review on pull requests and stores JSON or Markdown output. PR comments, GitHub Checks API integration, and inline annotations are deferred until the CLI is stable.
