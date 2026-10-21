#!/usr/bin/env bash
# Build the local release assets without creating tags, commits, or releases.
set -euo pipefail

version="${1:-${VERSION:-v0.1.0}}"
case "$version" in
  v*) ;;
  *) version="v${version}" ;;
esac

root="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
cd "$root"

export CGO_ENABLED=0
export GOTOOLCHAIN="${GOTOOLCHAIN:-go1.27.2}"

targets=(
  linux/amd64
  linux/arm64
  darwin/amd64
  darwin/arm64
  windows/amd64
  windows/arm64
)

rm -rf dist
mkdir -p dist

for target in "${targets[@]}"; do
  goos="${target%/*}"
  goarch="${target#*/}"
  ext=""
  if [[ "$goos" == "windows" ]]; then
    ext=".exe"
  fi
  output="dist/ua_${version}_${goos}_${goarch}${ext}"
  GOOS="$goos" GOARCH="$goarch" go build -trimpath -buildvcs=false -o "$output" ./cmd/ua
  [[ -s "$output" ]] || { printf 'missing build output: %s\n' "$output" >&2; exit 1; }
done

if [[ -r /proc/sys/kernel/random/uuid ]]; then
  serial="urn:uuid:$(cat /proc/sys/kernel/random/uuid)"
elif command -v uuidgen >/dev/null 2>&1; then
  serial="urn:uuid:$(uuidgen | tr '[:upper:]' '[:lower:]')"
else
  printf 'cannot generate a UUID for the SBOM\n' >&2
  exit 1
fi

created="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
purl="pkg:golang/github.com/IamDevTrivedi/useragent@${version}"
go_version="$(go env GOVERSION)"

cat > dist/sbom.cdx.json <<EOF
{
  "bomFormat": "CycloneDX",
  "specVersion": "1.5",
  "serialNumber": "${serial}",
  "version": 1,
  "metadata": {
    "timestamp": "${created}",
    "tools": {
      "components": [
        {
          "type": "application",
          "name": "go",
          "version": "${go_version}"
        }
      ]
    },
    "component": {
      "type": "application",
      "bom-ref": "${purl}",
      "name": "ua",
      "version": "${version}",
      "purl": "${purl}"
    }
  },
  "components": []
}
EOF

cd dist
if command -v sha256sum >/dev/null 2>&1; then
  sha256sum ua_* sbom.cdx.json > checksums.txt
else
  shasum -a 256 ua_* sbom.cdx.json > checksums.txt
fi

printf 'built %s release assets with %s\n' "$version" "$go_version"
cat checksums.txt
