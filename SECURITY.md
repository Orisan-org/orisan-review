# Security Policy

Orisan Review is local-first security tooling. Please report suspected vulnerabilities privately to the maintainers before public disclosure.

## Scope

In scope:

- Orisan Review CLI behavior
- Incorrect handling of diff input
- Evidence redaction failures
- Raw secret leakage in reports
- Unexpected network or cloud calls

Out of scope:

- Findings produced by test fixtures
- Vulnerabilities in third-party repositories scanned by Review
- Requests to add SAST, SCA, DAST, PR comments, or cloud upload behavior to v0.1

## Privacy Expectations

Review must not upload source code, diffs, findings, or reports by default. Findings must not store full patch content or raw secrets.
