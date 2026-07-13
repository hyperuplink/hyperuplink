#!/usr/bin/env bash
# Build Hyperuplink RPMs from the current git checkout.
#
# Usage:
#   deploy/rpm/build.sh [ARCH ...]      # default arches: x86_64 aarch64
#
# The build cross-compiles the Go binary for each ARCH on whatever host runs
# this, so a single x86_64 machine can produce every RPM (see the spec's
# GOARCH mapping). Run it on Fedora, or in a fedora:* container, with rpm-build,
# golang and systemd-rpm-macros installed.
#
# Overridable environment:
#   VERSION  package version   (default: git describe, leading "v" stripped)
#   COMMIT   commit stamped in  (default: git HEAD sha)
#   DATE     build date stamped (default: HEAD commit date, ISO-8601)
#   OUTDIR   where .rpm files land (default: <repo>/dist)
set -euo pipefail

repo_root=$(git rev-parse --show-toplevel)
cd "$repo_root"

version=${VERSION:-$(git describe --tags --always 2>/dev/null || echo 0.0.0)}
version=${version#v}
# RPM forbids '-' in Version; map pre-release separators to '~' so that, e.g.,
# 1.2.0-rc1 sorts *before* 1.2.0 as intended.
version=${version//-/\~}

commit=${COMMIT:-$(git rev-parse HEAD)}
date=${DATE:-$(git show -s --format=%cI HEAD 2>/dev/null || date -u +%Y-%m-%dT%H:%M:%SZ)}
outdir=${OUTDIR:-$repo_root/dist}

if [ "$#" -gt 0 ]; then
    arches=("$@")
else
    arches=(x86_64 aarch64)
fi

topdir=$(mktemp -d)
trap 'rm -rf "$topdir"' EXIT
mkdir -p "$topdir"/{BUILD,RPMS,SOURCES,SPECS,SRPMS}

# Source0: a clean tarball of the tracked tree, prefixed to match the spec's
# %autosetup (hyperuplink-<version>/...).
git archive --format=tar.gz --prefix="hyperuplink-${version}/" \
    -o "$topdir/SOURCES/hyperuplink-${version}.tar.gz" HEAD

# Source1: the sysusers.d snippet (lives outside the source tarball).
cp deploy/rpm/hyperuplink.sysusers "$topdir/SOURCES/"

echo ">> building hyperuplink ${version} (commit ${commit}) for: ${arches[*]}"
for arch in "${arches[@]}"; do
    echo ">> == ${arch} =="
    rpmbuild -bb deploy/rpm/hyperuplink.spec \
        --define "_topdir $topdir" \
        --define "version $version" \
        --define "commit $commit" \
        --define "date $date" \
        --target "$arch"
done

mkdir -p "$outdir"
find "$topdir/RPMS" -name '*.rpm' -exec cp -v {} "$outdir/" \;
echo ">> done, RPMs in ${outdir}:"
ls -1 "$outdir"/*.rpm
