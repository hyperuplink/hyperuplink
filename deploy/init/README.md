# Running Hyperuplink as Service

Here are service/init files for supervising the compiled `hyperuplink` binary
with your system's native init. Each file has full install steps in its header
comment.

| Init system            | File                  | Installs as                       |
| ---------------------- | --------------------- | --------------------------------- |
| systemd (Linux)        | `hyperuplink.service` | `/etc/systemd/system/`            |
| OpenRC (Alpine/Gentoo) | `hyperuplink.openrc`  | `/etc/init.d/hyperuplink`         |
| FreeBSD rc.d           | `hyperuplink.freebsd` | `/usr/local/etc/rc.d/hyperuplink` |
| OpenBSD rc.d           | `hyperuplink.openbsd` | `/etc/rc.d/hyperuplink`           |

These manage **only** the Hyperuplink process. PostgreSQL and Valkey/Redis must
be running and reachable per your config. Install and supervise them however you
normally would.

## Prerequisites (all platforms)

1. Grab a release archive from the GoReleaser builds, or `make build`, and
   install it to `/usr/local/bin/hyperuplink`
2. The service shells out to `convert` for image processing, so it must be on
   the service's `PATH` (`imagemagick` / `ImageMagick` package).
3. Copy `deploy/hyperuplink.toml` to the path the script expects
   (`/etc/hyperuplink.toml`, or `/usr/local/etc/hyperuplink.toml` on FreeBSD)
   and edit it. Point `[[Storage]]` `Local.Path` at the data directory the
   script creates for the service user.
4. Create a dedicated unprivileged user (`hyperuplink`, or `_hyperuplink` on
   OpenBSD), the exact command is in each file's header.

## Config path

The binary defaults to `-c file:///etc/hyperuplink.toml`, but the scripts pass
`-c` explicitly so you can point it elsewhere. The value is a URL, so a local
file always needs the `file://` prefix.

## Notes

- The port defaults to `3000` (non-privileged). To bind `<1024` directly, see
  the capability note in `hyperuplink.service`. On the BSDs/OpenRC run behind a
  reverse proxy instead. Actually, better always run behind a proxy, regardless
  of the system.
- In `Mode = "production"` the session cookie is HTTPS-only, so terminate TLS in
  front of the app (reverse proxy) or cookies will crumble.
- All scripts supervise the process and restart it on failure, which also covers
  the first-boot window while Postgres/Valkey come up.
