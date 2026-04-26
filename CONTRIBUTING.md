# Contributing to Sevro CLI

Thanks for considering a contribution. The CLI is the open, auditable surface of Sevro — everything in here is Apache-2.0 and stays that way.

## Quickstart

```sh
git clone https://github.com/lowplane/sevro cli
cd cli
make build              # produces ./bin/sevro
./bin/sevro --version
make test
```

## How we work

- **Issues first.** Open an issue describing the problem or proposal before sending a large PR. Small fixes (typos, obvious bugs) can skip this.
- **Conventional Commits.** `feat(analyze): support kustomize overlays`, `fix(parser): handle templated names`, etc.
- **DCO sign-off.** Every commit must carry a `Signed-off-by:` line (use `git commit -s`). We do not require a CLA.
- **One change per PR.** Refactors, feature work, and dependency bumps go in separate PRs.
- **Tests required for behavior changes.** Detectors and parsers ship with golden tests in `testdata/fixtures/`.
- **Public surface (`pkg/`)** changes need a clear note in the PR description explaining the deprecation path if any.

## What we accept

- New deterministic cost or security detectors (Year 1: 15 + 15 in the SaaS — we're additive in the CLI)
- Parser improvements for real-world Helm chart patterns
- Renderer improvements (better tables, colored output, JSON shape stability)
- Bug fixes with golden test reproductions

## What we don't accept

- LLM calls or AI-generated suggestions in the CLI itself (those live in the SaaS)
- Telemetry / analytics by default
- Windows-specific code paths (we ship Linux + macOS via npm)
- Dependencies with non-permissive licenses (Apache 2.0 / MIT / BSD only)

## Running locally

```sh
make build           # build binary
make test            # run unit + golden tests
make lint            # golangci-lint
make release-dryrun  # GoReleaser snapshot build
```

## Reporting security issues

**Do not file public issues for security bugs.** See [SECURITY.md](SECURITY.md) for the disclosure process.

## Code of conduct

By participating you agree to abide by the [Code of Conduct](CODE_OF_CONDUCT.md).
