# Contributing

Thanks for helping improve `useragent`. This project values small, well-tested changes
that keep the parser predictable and dependency-free.

## Before you start

- Search existing issues and pull requests.
- Open an issue for API changes, new detector families, or behavior changes that affect
  many User-Agent strings.
- Security issues must follow [SECURITY.md](SECURITY.md), not the public issue tracker.

## Development setup

Requirements:

- Go 1.24 or newer
- `golangci-lint` v2.x
- Optional: `staticcheck`, `govulncheck`

```sh
git clone https://github.com/IamDevTrivedi/useragent.git
cd useragent
go test ./...
make check
```

The module has no third-party runtime dependencies. Do not add one to the core package
without a clear design discussion.

## Making a change

1. Fork the repository and create a focused branch.
2. Add or update a test first whenever behavior changes.
3. Keep the public API small and documented.
4. Update the relevant documentation and `CHANGELOG.md`.
5. Run the full local checks:

```sh
make check
```

At minimum, run:

```sh
gofmt -w .
go test ./...
go test -race ./...
go vet ./...
golangci-lint run ./...
```

## Parser changes

- Add a regression fixture for every reported User-Agent.
- Do not weaken an existing assertion to make a new fixture pass. Decide which behavior
  is correct and document the choice.
- Prefer adding a data-driven rule to adding a special case in the pipeline.
- Keep detector ordering deliberate: specific clients must win over generic Chromium,
  WebKit, or Gecko tokens.
- Unknown input should remain unknown. Do not return an arbitrary product or comment
  token as a browser name.
- Never let a User-Agent control network access, filesystem paths, shell commands, or
  logging formats without application-level sanitization.
- Add a conformance case under `testdata/conformance.json` for broad coverage changes.
- Add a benchmark when changing the hot path.

Custom rule examples and built-in detector guidance are in [docs/rules.md](docs/rules.md).
Architecture details are in [docs/architecture.md](docs/architecture.md).

## Tests

Use typed tests and named subtests. Name behavioral regressions after the client or
behavior under test and preserve the exact fixture from the report when possible.

```sh
go test ./...
go test -run TestDetectionRegressions -v
go test -run '^$' -bench . -benchmem
make fuzz-all FUZZTIME=30s
```

Coverage must remain at or above 90%:

```sh
make cover
```

## Documentation

- Use clear, direct language.
- Keep code examples compilable.
- Update the API guide or README for behavior or API changes.
- Add a changelog entry under `Unreleased`.
- Do not claim universal User-Agent coverage; describe limitations honestly.

## Commit and pull request style

Use focused commits with imperative subjects, for example:

```text
parser: detect Acme Browser
tests: add QQ Browser device model fixture
docs: explain client hint precedence
```

Pull request descriptions should include:

- the problem
- the chosen behavior
- tests added
- benchmark impact, if the hot path changed
- documentation and changelog updates

Use the repository pull request template.

## API stability

The project is pre-1.0. Minor releases may include breaking changes, but they must be
documented in `CHANGELOG.md`. After `v1.0.0`, the project follows semantic versioning.

## Code of conduct

Participation is covered by [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md). Be respectful and
assume good intent.
