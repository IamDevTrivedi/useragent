# GitHub setup and release runbook

This document is for the repository owner. The local tree and repository automation are
prepared, but **no GitHub write operation has been performed**. No issue, pull request,
label, milestone, ruleset, release, tag, commit, or push was created.

Run the commands below from the repository root after reviewing the local tree.

## 0. Preconditions

- `gh` is installed and authenticated.
- The account has permission to create the repository and administer branch rules.
- The working tree is clean and the release checks pass.
- The final module path is confirmed.

```sh
gh auth status
git status
make check
```

If workflow files will be pushed, the GitHub token may need the `workflow` scope:

```sh
gh auth refresh -s workflow
```

## 1. Confirm the module path

The current module path is:

```text
github.com/IamDevTrivedi/useragent
```

If the final repository name or owner differs, update it everywhere before the first
commit:

```sh
grep -R "github.com/IamDevTrivedi/useragent" -n --exclude-dir=.git --exclude-dir=dist .
```

Update at least:

- `go.mod`
- Go import paths in `uahttp/` and `cmd/ua/`
- examples
- badges and documentation links
- the changelog release links
- this runbook

Then verify:

```sh
go list -m
go test ./...
```

## 2. Create the initial commit and repository

The tree is intentionally uncommitted locally. The owner performs the initial commit:

```sh
git add .
git commit -m "feat: initial useragent v0.1.0"
git branch -M main
```

Create the GitHub repository and push:

```sh
OWNER=IamDevTrivedi
REPO=useragent

gh repo create "$OWNER/$REPO" \
  --public \
  --source=. \
  --remote=origin \
  --description "Zero-dependency, typed User-Agent parser for Go"

git push -u origin main

gh repo edit "$OWNER/$REPO" \
  --add-topic go,golang,user-agent,ua-parser,browser-detection,bot-detection,client-hints,device-detection \
  --enable-discussions
```

If the repository already exists, add the remote manually:

```sh
git remote add origin "git@github.com:$OWNER/$REPO.git"
git push -u origin main
```

Confirm the workflows are visible under the Actions tab before continuing.

## 3. Labels

Create or update the repository labels. `--force` updates an existing label with the
same name.

```sh
gh label create "bug" --color "d73a4a" --description "Something is parsed or behaves incorrectly" --force
gh label create "enhancement" --color "a2eeef" --description "New capability or coverage" --force
gh label create "documentation" --color "0075ca" --description "Documentation and examples" --force
gh label create "parser" --color "5319e7" --description "Tokenizer or detector behavior" --force
gh label create "performance" --color "fbca04" --description "Speed, allocations, or benchmarks" --force
gh label create "bot-detection" --color "1d76db" --description "Bot, crawler, or HTTP client detection" --force
gh label create "client-hints" --color "0e8a16" --description "Sec-CH-UA and UA-CH handling" --force
gh label create "cli" --color "c2e0c6" --description "cmd/ua command-line tool" --force
gh label create "breaking-change" --color "b60205" --description "Requires an API break and a major release" --force
gh label create "dependencies" --color "0366d6" --description "Dependency or workflow updates" --force
gh label create "security" --color "ee0701" --description "Security-sensitive issue" --force
gh label create "needs-triage" --color "ededed" --description "Not yet reviewed" --force
gh label create "good first issue" --color "7057ff" --description "Good entry point for contributors" --force
gh label create "help wanted" --color "008672" --description "Community help requested" --force
gh label create "skip-changelog" --color "ffffff" --description "Exclude from generated release notes" --force
```

Verify:

```sh
gh label list --limit 100
```

The `.github/labeler.yml` configuration and `.github/workflows/labeler.yml` workflow
apply area labels to pull requests. The coverage issue form uses `parser` and
`needs-triage`.

## 4. Milestones

Create the release milestones:

```sh
gh api --method POST "repos/$OWNER/$REPO/milestones" \
  -f title="v0.1.0" \
  -f description="Initial parser release" \
  -f state="open"

gh api --method POST "repos/$OWNER/$REPO/milestones" \
  -f title="v0.2.0" \
  -f description="Coverage, performance, and feedback from the initial release" \
  -f state="open"

gh api --method POST "repos/$OWNER/$REPO/milestones" \
  -f title="v1.0.0" \
  -f description="Stable public API after pre-1.0 review" \
  -f state="open"
```

Verify:

```sh
gh api "repos/$OWNER/$REPO/milestones?state=all"
```

## 5. Create local tracking issues

Use GitHub Issues only for work in this repository. The project has no external issue
tracker or cross-repository mirroring process. When opening an issue, include:

- the exact User-Agent or header sample, with sensitive values removed;
- the expected and actual parsed fields;
- the library version and Go version;
- whether the sample is real-world or synthetic.

See [detection-matrix.md](detection-matrix.md) for the current supported surface. Local
issues and pull requests are the authoritative record for future work.

## 6. Branch protection

Create a ruleset for `main` after the first push. This ruleset requires pull requests,
blocks deletion and force pushes, and leaves the approval count at zero so a solo
maintainer is not locked out. If you add a second maintainer, raise
`required_approving_review_count` to `1`; GitHub does not allow an author to approve
their own pull request.

```sh
gh api --method POST "repos/$OWNER/$REPO/rulesets" --input - <<'JSON'
{
  "name": "main protection",
  "target": "branch",
  "enforcement": "active",
  "conditions": {
    "ref_name": {
      "include": ["refs/heads/main"],
      "exclude": []
    }
  },
  "rules": [
    { "type": "deletion" },
    { "type": "non_fast_forward" },
    {
      "type": "pull_request",
      "parameters": {
        "required_approving_review_count": 0,
        "dismiss_stale_reviews_on_push": true,
        "require_code_owner_reviews": false,
        "require_last_push_approval": false,
        "required_review_thread_resolution": true
      }
    }
  ]
}
JSON
```

After the first CI run, add required status checks using the exact check names shown in
the Actions UI. Do not guess names. The expected jobs are:

- `fmt-vet`
- `lint`
- `test (ubuntu-latest, 1.24.x)`
- `test (ubuntu-latest, 1.26.x)`
- `test (ubuntu-latest, 1.27.2)`
- `test (macos-latest, 1.27.2)`
- `test (windows-latest, 1.27.2)`
- `race`
- `coverage`
- `fuzz-smoke (FuzzParse)`
- `fuzz-smoke (FuzzTokenizer)`
- `fuzz-smoke (FuzzParseVersion)`
- `fuzz-smoke (FuzzParseBrandList)`
- `fuzz-smoke (FuzzParseRequest)`
- `bench`
- `build-cli (linux, amd64)`
- `build-cli (linux, arm64)`
- `build-cli (darwin, amd64)`
- `build-cli (darwin, arm64)`
- `build-cli (windows, amd64)`
- `build-cli (windows, arm64)`
- `govulncheck`
- `docs`
- `analyze (go)`

You can add checks in **Settings → Rules → Rulesets → main protection → Required status
checks**, or patch the ruleset through the API after inspecting the check names.

Optional release-tag protection:

```sh
gh api --method POST "repos/$OWNER/$REPO/rulesets" --input - <<'JSON'
{
  "name": "release tags",
  "target": "tag",
  "enforcement": "active",
  "conditions": {
    "ref_name": {
      "include": ["refs/tags/v*"],
      "exclude": []
    }
  },
  "rules": [
    { "type": "deletion" },
    { "type": "non_fast_forward" }
  ]
}
JSON
```

Verify:

```sh
gh api "repos/$OWNER/$REPO/rulesets"
```

## 7. Security settings

In **Settings → Code security and analysis**, enable:

- Dependabot alerts
- Dependabot security updates
- Secret scanning
- Push protection
- Private vulnerability reporting

The CodeQL workflow is in `.github/workflows/codeql.yml`. The nightly workflow runs
`govulncheck` and fuzz smoke tests. Dependabot configuration is in
`.github/dependabot.yml`.

After enabling private vulnerability reporting, confirm the link in `SECURITY.md` works.

## 8. First CI run

Push the initial commit or open a test pull request. Confirm:

- all hosted test jobs pass, including the Go 1.24 minimum-version job
- race and coverage jobs pass on Linux
- fuzz-smoke, benchmark, CLI-build, vulnerability, and docs jobs pass
- lint passes
- CodeQL completes
- the release workflow is syntactically valid
- the nightly workflow can be manually dispatched

Manual dispatch:

```sh
gh workflow run nightly.yml
gh run list --limit 10
```

If a workflow fails only because of a missing repository permission, check
**Settings → Actions → General → Workflow permissions**.

## 9. Release `v0.1.0`

Before tagging, re-run local checks and inspect the benchmark/changelog:

```sh
make check
git status
git diff --check
```

To build the six release binaries, checksums, and SBOM locally without publishing
anything, run:

```sh
make release-assets VERSION=v0.1.0
```

This writes to `dist/` only. It does not create a tag, commit, or GitHub release.

Create a `release` environment in **Settings → Environments** and add at least one
required reviewer. The release job will wait for approval before publishing. If the
repository has a second maintainer, require that person; otherwise approve the
deployment yourself.

Create an annotated tag:

```sh
git tag -a v0.1.0 -m "useragent v0.1.0"
git push origin v0.1.0
```

The release workflow builds Linux, macOS, and Windows binaries for `amd64` and `arm64`,
generates a checksum file, produces a CycloneDX SBOM, and publishes the matching
`CHANGELOG.md` section as the release notes. The `release` environment is configured
with required reviewers; approve the environment deployment after the build succeeds.

Verify:

```sh
gh release view v0.1.0
gh release download v0.1.0 --pattern 'checksums.txt' --dir /tmp/ua-release
cat /tmp/ua-release/checksums.txt
```

Expected assets:

- `ua_v0.1.0_linux_amd64`
- `ua_v0.1.0_linux_arm64`
- `ua_v0.1.0_darwin_amd64`
- `ua_v0.1.0_darwin_arm64`
- `ua_v0.1.0_windows_amd64.exe`
- `ua_v0.1.0_windows_arm64.exe`
- `checksums.txt`
- `sbom.cdx.json`

If the workflow fails after the tag is pushed, fix the workflow, delete the remote tag,
and re-tag only after the fix is merged. Tag deletion is a GitHub write and is the
owner's decision.

## 10. Post-release checklist

- [ ] Confirm pkg.go.dev indexes the module.
- [ ] Confirm the README badges resolve.
- [ ] Verify `go get github.com/IamDevTrivedi/useragent@v0.1.0`.
- [ ] Verify `go install github.com/IamDevTrivedi/useragent/cmd/ua@v0.1.0`.
- [ ] Verify the release checksum and SBOM assets.
- [ ] Open follow-up issues for any documented gaps.

## 11. Repository independence

This repository is self-contained. Do not add links to external trackers, reproduce
external issue numbers, or mirror external history in release notes, tests, or
documentation. Use local behavior names and local issue references only.

## 12. Rollback

- Workflow-only fix: merge a new commit to `main`.
- Failed release build: delete the remote tag, fix, and re-tag.
- Bad release artifact: mark the release as a pre-release or delete it, then publish a
  patch release.
- Ruleset mistake: disable the ruleset temporarily in repository settings rather than
  force-pushing over protected history.

All of these are owner-only GitHub writes. The local project automation does not perform
them.
