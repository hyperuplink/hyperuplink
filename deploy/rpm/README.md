# RPM package for Hyperuplink

This directory builds Hyperuplink into an `.rpm` that installs the single,
statically-linked binary together with a production config, a systemd unit and a
dedicated service user.

```
deploy/rpm/
├── hyperuplink.spec       # RPM spec (builds from source)
├── hyperuplink.sysusers   # sysusers.d snippet (service account)
├── build.sh               # helper (build RPMs for one or more arches)
└── README.md
```

The packages target Fedora, but are written to install and run on the
RHEL-family (Rocky/Alma/RHEL) and hopefully openSUSE too.

## Installation

| Path                                          | Contents                                                       |
| --------------------------------------------- | -------------------------------------------------------------- |
| `/usr/bin/hyperuplink`                        | the binary                                                     |
| `/etc/hyperuplink.toml`                       | starter config (`0640 root:hyperuplink`, `%config(noreplace)`) |
| `/usr/lib/systemd/system/hyperuplink.service` | systemd unit                                                   |
| `/usr/lib/sysusers.d/hyperuplink.conf`        | `hyperuplink` user/group                                       |
| `/var/lib/hyperuplink/media`                  | local media storage (`0750 hyperuplink:hyperuplink`)           |

The `hyperuplink` system user/group is created on install. ImageMagick
(`convert`) is pulled in as a runtime dependency. PostgreSQL and Redis/Valkey
are not hard dependencies, configure them in `/etc/hyperuplink.toml`.

## Building

The build cross-compiles the Go binary for each target arch, so a single x86_64
host produces every RPM. It needs `rpm-build`, `golang` and
`systemd-rpm-macros`.

```sh
# Build default arches from a checked-out tag (or any commit):
deploy/rpm/build.sh

# ... or explicit:
deploy/rpm/build.sh x86_64 aarch64 ppc64le

# With version/commit/date override:
VERSION=1.2.0 COMMIT=$(git rev-parse HEAD) deploy/rpm/build.sh
```

RPMs land in `dist/` (override with `OUTDIR`). Under the hood `build.sh` just
drives `rpmbuild --target <arch>` against `hyperuplink.spec`, which maps the RPM
target CPU onto a Go `GOARCH`.

## Installing

```sh
# Fedora
sudo dnf install ./hyperuplink-*.x86_64.rpm

# openSUSE
sudo zypper install ./hyperuplink-*.x86_64.rpm
```

Note: The `.fcNN` in the filename only shows that the package was built on
Fedora, but the RPM is not tied to that release.

### Rocky / RHEL-family

ImageMagick lives in EPEL on RHEL-family distros, so enable it first:

```sh
sudo dnf install epel-release
sudo dnf install ./hyperuplink-*.x86_64.rpm
```

## Configure and start

Edit `/etc/hyperuplink.toml` (database DSN, Redis address, `PromoteAdmin`, …).
The shipped config points at loopback, so a single-host setup with local
PostgreSQL and Redis/Valkey works after filling in credentials. Then:

```sh
sudo systemctl enable --now hyperuplink.service
```

The service listens on port `3000`. Run it behind a TLS-terminating reverse
proxy, in `Mode = "production"` so the session cookie is HTTPS-only.
