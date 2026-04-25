# Security Policy

## Reporting a Vulnerability

**Please do not open public GitHub issues for security vulnerabilities.**

Email `security@costify.dev` with:

- A description of the vulnerability and its potential impact
- Reproduction steps (proof of concept where possible)
- Affected versions
- Your name / handle for credit (optional)

We acknowledge reports within **2 business days**. We aim to ship a fix within **30 days** for high-severity issues.

You may encrypt your report with our PGP key (fingerprint: TBD when key is published).

## Supported Versions

The latest minor release of `@costify/cost` receives security updates. Older versions are not supported.

## Scope

In scope:

- The `costify` binary (this repo)
- The `@costify/cost` npm wrapper

Out of scope:

- The Costify SaaS (report to `security@costify.dev` separately)
- Third-party dependencies (please report to upstream maintainers)

## Disclosure

We will publicly disclose accepted vulnerabilities after a fix is released, crediting the reporter unless they request otherwise.
