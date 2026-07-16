#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$REPO_ROOT"

UPSTREAM_REPO="https://github.com/golang/tools"
UPSTREAM_PKG="internal/diff"
MODULE_PATH="github.com/aymanbagabas/go-udiff"

# Clone upstream to a temporary directory.
tools=$(mktemp -d)/tools
trap 'rm -rf "$(dirname "$tools")"' EXIT
git clone --depth 1 "$UPSTREAM_REPO" "$tools"

# Copy the diff package to the current directory.
cp -r "$tools/$UPSTREAM_PKG/"* .

# Replace import paths. Use perl for portable in-place editing (GNU/BSD sed
# disagree on the -i syntax; perl -i behaves the same on Linux and macOS).
find . -type f -name '*.go' -exec perl -i -pe "s|golang.org/x/tools/${UPSTREAM_PKG}/|${MODULE_PATH}/|g" {} +
find . -type f -name '*.go' -exec perl -i -pe "s|\"golang.org/x/tools/${UPSTREAM_PKG}|diff \"${MODULE_PATH}|g" {} +

# Change package name to udiff.
perl -i -pe 's|package diff|package udiff|g' *.go

# Apply patches.
for p in _patches/*; do
	git apply "$p"
done

# Resolve upstream commit.
commit=$(cd "$tools" && git rev-parse HEAD)

# Update the upstream commit marker if there are changes.
if ! git update-index --refresh >/dev/null 2>&1 || ! git diff-index --quiet HEAD --; then
	echo "$commit" >.github/UPSTREAM
	echo "Upstream updated to $commit"
else
	echo "No changes from upstream"
fi
