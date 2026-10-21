# Security Policy

## Supported versions

Security fixes are provided for the latest release in the current pre-1.0 line. Once
`v1.0.0` is published, the latest stable minor release will receive security fixes.

## Reporting a vulnerability

Do not open a public issue for a suspected vulnerability. Report it privately through
GitHub Security Advisories:

<https://github.com/IamDevTrivedi/useragent/security/advisories/new>

Include:

- a description of the issue and its impact
- the affected version or commit
- a minimal reproduction or User-Agent string
- whether the issue requires application-specific assumptions
- any suggested fix

The maintainer will acknowledge the report as soon as practical, investigate, and
coordinate disclosure. Please do not disclose the issue publicly until a fix or
mitigation is available.

## Scope

`useragent` parses untrusted text. Relevant security concerns include:

- panics or excessive resource use on malformed or very large input
- incorrect parsing that could cause a security decision to be made on spoofed data
- parser output being used unsafely in shell commands, SQL, HTML, file paths, or logs

The parser performs no network access, DNS lookup, file access, or command execution.

## Trust model

A User-Agent and Client Hints are **untrusted, client-controlled input**. They can be
spoofed, omitted, frozen, or shared by many devices. Do not use this package to make
authentication, authorization, fraud, or unique-device-identification decisions. Use
server-side session identifiers and application authentication for those purposes.

Client hints are also client-controlled. Treat them as higher-quality parsing hints,
not verified facts.

## Dependency policy

The core module has no third-party runtime dependencies. New dependencies require a
strong justification and review. CI runs CodeQL, `govulncheck`, fuzz smoke tests, and
the standard Go test suite.
