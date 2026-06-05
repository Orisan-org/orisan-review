# Orisan Review

Orisan Review is a local-first CLI that analyzes pull request diffs for security-sensitive AI-assisted output risk.

It helps answer:

> What changed, why does it matter, and who needs to review this before it is merged or trusted?

Review v0.1 is diff-aware risk routing, not generic SAST.

## What Review Is

Review is a portfolio/community tool for validating AI-assisted output risk before trust. It is designed to run locally or in CI, inspect git diffs or patch files, classify security-sensitive changes, and route attention to the right reviewer.

## What Review Is Not

Review v0.1 is not SAST, SCA, DAST, or an AI code reviewer. It does not prove that a vulnerability exists, does not replace human review, does not upload source code, does not call an LLM, and does not post PR comments by default.

## Current Status

This repository is at the scaffold stage. The CLI command surface and package layout exist, but diff parsing and rule logic are intentionally not implemented yet.

## Build

```sh
go fmt ./...
go vet ./...
go test ./...
go build -o bin/orisan-review ./cmd/orisan-review
```

## Quickstart

```sh
orisan-review --help
orisan-review version
orisan-review list-rules
orisan-review list-categories
```

The following commands are present but return a not-implemented input error until the parser and engine are wired:

```sh
orisan-review diff --base main --head HEAD
orisan-review diff --staged
orisan-review diff --worktree
orisan-review scan-patch ./change.patch
git diff main...HEAD | orisan-review scan-patch -
```

## Privacy

Review is local-first by design. v0.1 must not upload source code, diffs, findings, or reports to Orisan. Findings must store short redacted evidence only and always set `payload_stored=false`.

## GitHub Integration

The recommended v0.1 integration is a GitHub Actions workflow that runs Review on pull requests and stores JSON or Markdown output. PR comments, GitHub Checks API integration, and inline annotations are deferred until the CLI is stable.
