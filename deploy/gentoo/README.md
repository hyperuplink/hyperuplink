# Gentoo ebuild for Hyperuplink

This directory is a small, drop-in [Gentoo][gentoo] ebuild repository (overlay)
that builds Hyperuplink from source and installs it as a single binary with
`emerge`. Every ebuild here is the same file, which branches on `${PV}`:

- `hyperuplink-9999.ebuild` is the canonical one, and it is also the live ebuild
  that builds the tip of the GitHub default branch. It needs no deps tarball,
  because it vendors the modules during `src_unpack`.
- `hyperuplink-0.1.0.ebuild`, and any other released version, is a byte-for-byte
  copy of it that `make release` writes. Portage takes `${PV}` from the
  filename, so the copy takes the `else` branch and builds from the GitHub tag
  tarball, with the Go module cache supplied as a separate _deps tarball_ that
  the Release workflow attaches (see below).

Edit `hyperuplink-9999.ebuild` and only that one, since anything else is
overwritten on the next release.

```
deploy/gentoo/
├── acct-group/hyperuplink/     # `hyperuplink` group (GLEP 81)
├── acct-user/hyperuplink/      # `hyperuplink` service user
├── www-apps/hyperuplink/       # package
│   ├── hyperuplink-0.1.0.ebuild
│   ├── hyperuplink-9999.ebuild
│   └── metadata.xml
├── licenses/SEGV-1.1           # SEGV license, v1.1
├── metadata/layout.conf
└── profiles/repo_name
```

## Requirements

- `>=dev-lang/go-1.26.4` (`GOTOOLCHAIN=local` forbids fetching a newer one).
- ImageMagick (`media-gfx/imagemagick`, for `convert`) is pulled in as a runtime
  dependency.
- A reachable PostgreSQL server, a Redis/Valkey server, and optionally an
  S3-compatible storage like MinIO, which are not hard dependencies and can be
  configured in `/etc/hyperuplink.toml`.

## Steps

### 1. Register the overlay

```sh
# As root, adjust the path to wherever this directory is located.
cat > /etc/portage/repos.conf/hyperuplink.conf <<'EOF'
[hyperuplink]
location = /path/to/hyperuplink/deploy/gentoo
masters = gentoo
auto-sync = no
EOF
```

### 2. License and version

```sh
# The SEGV license is non-OSI and must be accepted explicitly for this package.
echo 'www-apps/hyperuplink SEGV-1.1' >> /etc/portage/package.license

# Either the released version:
echo 'www-apps/hyperuplink ~amd64' >> /etc/portage/package.accept_keywords

# ... or the live ebuild instead:
echo 'www-apps/hyperuplink **' >> /etc/portage/package.accept_keywords
```

### 3a. Released version

The Go module cache (`hyperuplink-0.1.0-deps.tar.xz`) is built automatically by
`.github/workflows/release.yml` on every `v*` tag and attached to the GitHub
release, so the ebuild's `SRC_URI` fetches it with no manual step.

To build it by hand (e.g. testing an unreleased tag), reproduce what the CI
does:

```sh
git clone https://github.com/hyperuplink/hyperuplink.git
cd hyperuplink
git checkout v0.1.0

# Populate a module cache and pack it (-modcacherw keeps it writable so Portage
# can clean it up afterwards).
GOMODCACHE="${PWD}/go-mod" go mod download -modcacherw
tar --create --auto-compress --file hyperuplink-0.1.0-deps.tar.xz go-mod

# Make it available locally without a release upload:
cp hyperuplink-0.1.0-deps.tar.xz "$(portageq distdir)"/
```

Then generate the Manifest and install:

```sh
cd /path/to/hyperuplink/deploy/gentoo/www-apps/hyperuplink
ebuild hyperuplink-0.1.0.ebuild manifest
emerge -av www-apps/hyperuplink
```

### 3b. Live version

No deps tarball and no Manifest are needed, because the ebuild vendors modules
itself.

```sh
emerge -av "=www-apps/hyperuplink-9999"
```

## 4. Configure and start

Either way, `emerge` creates the `hyperuplink` user/group, installs the binary
to `/usr/bin/hyperuplink`, a starter config to `/etc/hyperuplink.toml`, an
OpenRC service, a systemd unit, and the data directory `/var/lib/hyperuplink`.

Edit `/etc/hyperuplink.toml` (database DSN, Redis address, `PromoteAdmin`, …)
as root, since it holds the database password and is installed `0640
root:hyperuplink`, then:

```sh
# OpenRC
rc-update add hyperuplink default
rc-service hyperuplink start

# systemd
systemctl enable --now hyperuplink.service
```

The service listens on port `3000`. Run it behind a TLS-terminating reverse
proxy, in `Mode = "production"` so that the session cookie is HTTPS-only.

## Notes

- Hyperuplink is under the custom, non-OSI **SEGV** license, whose text
  (version 1.1) is in `licenses/SEGV-1.1`. `::gentoo` does not carry it, which
  is why it is bundled here and must be accepted via `package.license`.
- `make release VERSION=<new>` writes `hyperuplink-<new>.ebuild` for you, as a
  copy of the live one. The Manifest is not part of that, because it hashes the
  deps tarball that only exists once the Release workflow has attached it to the
  new tag, so regenerate it here afterwards. The `-ldflags -X main.Version`
  value is wired to `${PV}`, so `hyperuplink -v` reports the ebuild version (the
  live ebuild reports `9999` plus the checked-out commit).
- GitHub's tag archive extracts to `hyperuplink-<version>/`, which is
  go-module.eclass's default `S`, so the released ebuild sets no explicit `S`.

[gentoo]: https://www.gentoo.org/
