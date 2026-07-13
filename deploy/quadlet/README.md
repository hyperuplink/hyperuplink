# Hyperuplink + Podman Quadlets (systemd)

These
[Quadlet](https://docs.podman.io/en/latest/markdown/podman-systemd.unit.5.html)
units are the systemd-native equivalent of the `podman-compose.yml` setup.
Podman generates a `systemd` service from each file, so the whole stack starts
on boot, restarts on failure, and is managed with `systemctl`.

## Stack

| File                             | Container / resource                 |
| -------------------------------- | ------------------------------------ |
| `hyperuplink.container`          | the app (`ghcr.io/hyperuplink/hyperuplink`) |
| `hyperuplink-postgres.container` | PostgreSQL 17                        |
| `hyperuplink-valkey.container`   | Valkey 8 (Redis-compatible)          |
| `hyperuplink-minio.container`    | MinIO / S3 — **optional**            |
| `hyperuplink.network`            | shared bridge network                |
| `hyperuplink-*.volume`           | persistent data volumes              |

Containers find each other through the `NetworkAlias` in each `.container`
(`postgres`, `valkey`, `minio`), which is exactly what `deploy/hyperuplink.toml`
points at, so that the config file is used unchanged.

Requires Podman 4.4+ (validated against 5.x).

---

## Rootless install (recommended)

Runs under your own user account without the need for root. The app listens on
`:3000` and everything runs as a non-root user.

1. Copy the sample config and edit the credentials (Postgres password,
   `PromoteAdmin` e-mail, etc.):

   ```sh
   mkdir -p ~/.config/hyperuplink
   cp deploy/hyperuplink.toml ~/.config/hyperuplink/hyperuplink.toml
   chmod 644 ~/.config/hyperuplink/hyperuplink.toml
   $EDITOR ~/.config/hyperuplink/hyperuplink.toml
   ```

   The `hyperuplink.container` unit bind-mounts this path
   (`%h/.config/hyperuplink/hyperuplink.toml`). `UserNS=keep-id` in that unit
   maps your user to the image's UID so the file stays readable inside the
   container.

2. Copy them into the rootless quadlet directory, skip `hyperuplink-minio.*`
   unless you want S3 storage:

   ```sh
   mkdir -p ~/.config/containers/systemd
   cp deploy/quadlet/hyperuplink.network \
      deploy/quadlet/hyperuplink-pgdata.volume \
      deploy/quadlet/hyperuplink-valkeydata.volume \
      deploy/quadlet/hyperuplink-media.volume \
      deploy/quadlet/hyperuplink-postgres.container \
      deploy/quadlet/hyperuplink-valkey.container \
      deploy/quadlet/hyperuplink.container \
      ~/.config/containers/systemd/
   ```

3. A daemon reload runs the quadlet generator, then start the app (its
   dependencies start automatically):

   ```sh
   systemctl --user daemon-reload
   systemctl --user start hyperuplink.service
   ```

4. Rootless user services stop when you log out unless lingering is enabled:

   ```sh
   loginctl enable-linger "$USER"
   ```

Hyperuplink is now on <http://localhost:3000>.

---

## Rootful install (system-wide)

Managed by the system `systemd`, starts at boot without lingering. Two edits to
`hyperuplink.container` first:

- Delete the `UserNS=keep-id:uid=10001,gid=10001` line (rootless only).
- Change the config `Volume=` host path from `%h/.config/hyperuplink/...` to an
  absolute path such as `/etc/hyperuplink/hyperuplink.toml`.

Then:

```sh
sudo mkdir -p /etc/hyperuplink
sudo cp deploy/hyperuplink.toml /etc/hyperuplink/hyperuplink.toml
sudo $EDITOR /etc/hyperuplink/hyperuplink.toml

sudo cp deploy/quadlet/hyperuplink.network \
        deploy/quadlet/hyperuplink-*.volume \
        deploy/quadlet/hyperuplink-postgres.container \
        deploy/quadlet/hyperuplink-valkey.container \
        deploy/quadlet/hyperuplink.container \
        /etc/containers/systemd/

sudo systemctl daemon-reload
sudo systemctl start hyperuplink.service
```

(drop `hyperuplink-miniodata.volume` from the `*.volume` glob if you're not
using MinIO/S3 storage)

---

## Optional: MinIO (S3 storage)

Only needed if you select the `remote-storage` (S3) provider in the admin UI.
Add its two files alongside the rest and reload:

```sh
cp deploy/quadlet/hyperuplink-miniodata.volume \
   deploy/quadlet/hyperuplink-minio.container \
   ~/.config/containers/systemd/          # or /etc/containers/systemd/ (rootful)
systemctl --user daemon-reload            # sudo systemctl ... for rootful
systemctl --user start hyperuplink-minio.service
```

Console: <http://localhost:9001>.

---

## Management

```sh
systemctl --user status hyperuplink.service        # is it up?
journalctl --user -u hyperuplink.service -f        # follow app logs
systemctl --user restart hyperuplink.service       # restart just the app
podman auto-update                                 # pull a newer app image + restart
```

Stop everything by stopping the datastores (the app depends on them):

```sh
systemctl --user stop \
  hyperuplink.service \
  hyperuplink-postgres.service \
  hyperuplink-valkey.service
```

(use `sudo systemctl ...` without `--user` for a rootful install)

---

## Notes

- The app publishes `3000`, MinIO `9000`/`9001`, which are all >1024, so
  rootless Podman binds them without extra privileges.
- Behind a reverse proxy, set `ProxyHeader`/`TrustProxy` in `hyperuplink.toml`.
- Change the Postgres password (and MinIO keys, if used) in both the
  `.container` file **and** the matching values in `hyperuplink.toml` before
  exposing this to the internet.
- The config bind mount carries `:z`. On hosts without SELinux (non-Fedora/RHEL)
  you can drop it, but it's harmless if left in.
- The app only _orders_ after Postgres/Valkey and it does not wait for them to
  become healthy. `Restart=always` retries the app until the datastores accept
  connections, which takes a few seconds on first boot.
