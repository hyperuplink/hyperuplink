# Development

This document describes the process for running this application on your local
computer.

## Getting started

Everything except the following is embedded into the binary, which means that
templates, views, static assets, locales, migrations and docs all travel inside
the executable and there are no loose files to chase:

- [Go](https://go.dev), at the version named in the `go` directive in `go.mod`
- [podman](https://podman.io), or _docker_, in which case substitute the command
  everywhere below
- `psql` and `pg_dump`, which the `db:*` _Makefile_ targets shell out to
- [ImageMagick](https://imagemagick.org), specifically `convert` on your
  `$PATH`, which is used to convert profile pictures

```sh
git clone https://codeberg.org/hyperuplink/hyperuplink.git
cd hyperuplink
make build
```

## Building

Build with `make build` rather than a bare `go build`, because the _Makefile_
includes the version, commit and build date in the binary through `-ldflags`.

```sh
make build         # -> ./build/hyperuplink
make help          # every target
```

## The local stack

_PostgreSQL_ and a _Redis_-compatible server are the only two pieces that are
needed in order to boot, whereas _MinIO_ is only ever touched once a file is
actually uploaded, and _Prosody_ only once a real _XMPP_ target is configured,
which means that the latter two are worth starting only when you need them.

> **Note:** Do not use the `podman-compose.yml` for local development. It
> deploys the released image from `ghcr.io` against `deploy/hyperuplink.toml`
> with its own database name and its own credentials, and it does not publish
> _PostgreSQL_ to the host, so a binary that you have just built cannot reach
> it.

### PostgreSQL

```sh
podman run -d --name postgres \
  -p 5432:5432 \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=hyperuplink_dev \
  docker.io/library/postgres:18-alpine
```

Migrations run on startup, just make sure the database exists.

### Valkey

The job queue (_asynq_) and the session store (_fiber_) both live here, and any
_Redis_-compatible server will do.

```sh
podman run -d --name redis \
  -p 6379:6379 \
  docker.io/valkey/valkey:8-alpine
```

### MinIO

_Optional:_ only needed for profile pictures and attachments served from an `S3`
storage provider, since `hyperuplink.toml` also ships a `local-storage` provider
that can write to e.g. `/tmp/www` and requires no container at all.

```sh
podman run -d --name minio \
  -p 9000:9000 -p 9001:9001 \
  -e MINIO_ROOT_USER=minioadmin \
  -e MINIO_ROOT_PASSWORD=minioadmin \
  quay.io/minio/minio:latest server /data --console-address ":9001"
```

The console is on <http://localhost:9001> behind `minioadmin`/`minioadmin`, and
the bucket is the first segment of the storage path.

### Prosody

_Optional:_ only needed in order to verify the _XMPP_ path against a real
server, because day to day the `debug-xmpp` target renders _XMPP_ messages to
disk and no server is involved.

_Hyperuplink_ insists on _StartTLS_ and will not authenticate over an
unencrypted connection, which means that _Prosody_ needs a certificate even on
loopback, though a self-signed one is sufficient as long as the target sets
`XMPP.InsecureSkipVerify = true`.

> **Warning:** `XMPP.InsecureSkipVerify = true` turns off certificate
> verification, and it has no business anywhere other than a local development
> server. With verification disabled, anybody who is able to sit between the two
> ends of the connection may present whatever certificate they like and then
> read or rewrite the traffic, credentials included.

```sh
mkdir -p /tmp/prosody/certs

openssl req -x509 -newkey rsa:2048 -nodes -days 2 \
  -keyout /tmp/prosody/certs/localhost.key \
  -out    /tmp/prosody/certs/localhost.crt \
  -subj "/CN=localhost" -addext "subjectAltName=DNS:localhost"
chmod 644 /tmp/prosody/certs/*

cat > /tmp/prosody/prosody.cfg.lua <<'EOF'
admins = { }
plugin_paths = { }

modules_enabled = {
  "roster"; "saslauth"; "tls"; "disco"; "private";
  "vcard4"; "vcard_legacy"; "version"; "uptime";
  "time"; "ping"; "offline"; "posix";
}

allow_registration = false
c2s_require_encryption = true
s2s_require_encryption = false
authentication = "internal_plain"

log = { info = "*console"; }
pidfile = "/var/run/prosody/prosody.pid"

ssl = {
  key = "/etc/prosody/certs/localhost.key";
  certificate = "/etc/prosody/certs/localhost.crt";
}

VirtualHost "localhost"
EOF

podman run -d --name prosody \
  -p 5222:5222 \
  -v /tmp/prosody/prosody.cfg.lua:/etc/prosody/prosody.cfg.lua:ro,Z \
  -v /tmp/prosody/certs:/etc/prosody/certs:ro,Z \
  docker.io/prosody/prosody:latest

# one account for the board to send as, and one to receive with
podman exec prosody prosodyctl register hyperuplink localhost boardpass
podman exec prosody prosodyctl register postbob     localhost bobpass
```

The target itself is then declared in the config and selected at
`/admin/comms/xmpp`:

```toml
[[Target]]
ID = "jabber"
Type = "xmpp"
XMPP.Server = "localhost:5222"
XMPP.InsecureSkipVerify = true   # only ever when testing locally
XMPP.Username = "hyperuplink@localhost"
XMPP.Password = "boardpass"
```

### Starting and stopping

```sh
podman start postgres redis        # and minio, and prosody, if you use them
```

## Running

```sh
make run 
# ... which build and runs:
./build/hyperuplink -c "file://$PWD/hyperuplink.toml"
```

The board then answers on <http://localhost:3000>, and since both `SIGINT` and
`SIGTERM` are trapped, either `Ctrl-C` or a plain
`kill <pid>`/`killall hyperuplink` will shut it down cleanly.

> **Note:** `-c` takes a URL rather than a path, hence the `file://`, and it
> defaults to `/etc/hyperuplink.toml`.

> **Note:** if something is already bound to `:3000`, the instance will still
> start, but its web server will simply fail to bind without saying much about
> it, at which point your requests continue to land on whatever is listening on
> the port. Check with `ss -ltn | grep :3000` before concluding that anything
> else is wrong.

TODO: development mode watches `views/` and reloads the templates as you save
them, however `getWatcher` hands the directory straight to
`fsnotify.Watcher.Add`, which watches that one directory rather than the tree
beneath it, and since `views/root.html` is the only template that is not in a
subdirectory, it is also the only one the watcher ever fires for. Every other
template is parsed once at startup and then never again, which means that a
change to one of them is not visible until you restart, and that until you do,
the previous markup keeps being served with nothing in the log to suggest why.
Editing a template and finding the page unchanged is this, rather than a caching
problem or a mistake in the template. This needs to be fixed.

### Command line flags

| Flag                  | Purpose                                            |
| --------------------- | -------------------------------------------------- |
| `-c <url>`            | Config URL, default `file:///etc/hyperuplink.toml` |
| `-v`                  | Print version, commit and build date, then exit    |
| `-reset <HH:MM>`      | **Wipe the entire database** and exit              |
| `-create-user <json>` | Create an activated user and exit                  |

#### `-reset`

This one drops everything, and because muscle memory is what it is, it asks for
the current system time as confirmation, in 24 hour `HH:MM` format, and refuses
to do anything at all if the minute does not match:

```sh
./build/hyperuplink -c "file://$PWD/hyperuplink.toml" --reset 13:45
```

Of course, if you simply decide to use `"$(date +%H:%M)"`, and then one day
`ctrl+r set` in your shell, accidentally hit `enter` on the `-reset` command
that uses the `date` call, and through that shoot yourself in the foot, please
don't come complaining.

For local development you can also just use `make db:drop`, at least as long as
your database is called `hyperuplink_dev` and you have the `psql` login stored
in `~/.pgpass`.

#### `-create-user`

This takes the signup form as a JSON object and creates the user **already
activated**, which skips the e-mail confirmation step and, on a fresh database,
is the fastest way to an account you can log in with:

```sh
./build/hyperuplink -c "file://$PWD/hyperuplink.toml" \
  --create-user '{"username":"dummy1","email":"dummy1@example.com","password":"mypassword"}'
```

The fields are those of `logic/root/session.SignUpInput`, which is to say
`username`, `email` and `password`, plus the optional `password_repeat`,
`email_is_jid` and `language`.

For an admin, use an address that is listed under `[Users] PromoteAdmin` in
`hyperuplink.toml`, where `dummy1@example.com` already sits, since signup
promotes those addresses on its own.

## Configuration

`hyperuplink.toml` in the repository root is the development config.

- `[General] Mode = "development"` enables debug logging and drops the `Secure`
  flag from the session and CSRF cookies, which is what allows plain-HTTP `curl`
  to round-trip them, whereas production mode serves the embedded views and
  forces HTTPS-only cookies.
- `[Database] Connection` is the `postgres://` DSN.
- `[Redis] Addresses` takes one address for a single node or several for a
  cluster, and `MasterName` for sentinel.
- `[Users] PromoteAdmin` lists the addresses that become admin on signup.
- `[[Storage]]` covers `local-storage`, which writes to `/tmp/www` and needs
  nothing running, and `remote-storage`, which is an `S3`-compliant storage
  (e.g. _MinIO_).
- `[[Target]]` is where notifications go, and is covered below.

### Notification targets

A target is _where_ an outbound notification ends up, and there are three kinds:

| Type    | Delivers via                           | Templates                 |
| ------- | -------------------------------------- | ------------------------- |
| `email` | SMTP, through _go-mail_                | `.eml` (text) and `.html` |
| `xmpp`  | XMPP, through _go-xmpp_, over StartTLS | `.md` only                |
| `debug` | A file on disk, or the log             | whichever it emulates     |

They are declared in `hyperuplink.toml` and then assigned to a channel at
`/admin/comms/email` and `/admin/comms/xmpp`.

The development config ships two `debug` targets, which is why signup
confirmations and reply notifications work out of the box without an SMTP or
_XMPP_ server being involved anywhere:

```toml
[[Target]]
ID = "debug-email"
Type = "debug"
Debug.Emulates = "email"
Debug.Path = "/tmp/hyperuplink-messages"
```

A `debug` target renders through the _same_ template pipeline as the real thing,
so what lands on disk is what would have gone out, written one file per
recipient with mode `0600` inside a `0700` directory, since the messages carry
signup tokens.

```sh
ls -1 /tmp/hyperuplink-messages/
cat /tmp/hyperuplink-messages/*
```

```
Target: debug-email (emulates email)
Job: confirmation/signup id=019f64dc-7248-79f8-8ca7-b3cc23888087
To: alice <alice@example.org>
Lang: en
Subject: Signup Confirmation

--- text ---
Hello alice,
...
--- html ---
<p>Hello alice,</p>
...
```

Leaving `Debug.Path` empty sends the messages to the log instead.

When a comms setting has no target selected at all, the `debug` target that
emulates that channel is used automatically, which is why a freshly dropped
database delivers without anybody having to click through the admin UI first.
Production declares no `debug` target, so nothing falls back there.

## Seeding a board

`Startup()` seeds the _settings_ rows and nothing besides, meaning no users, no
categories, no forums and no topics, so a freshly dropped database is an empty
board that nobody is able to log into.

```sh
# an admin, since the address is in PromoteAdmin it lands as one
./build/hyperuplink -c "file://$PWD/hyperuplink.toml" \
  --create-user '{"username":"admin1","email":"dummy1@example.com","password":"mypassword"}'

# and a regular user
./build/hyperuplink -c "file://$PWD/hyperuplink.toml" \
  --create-user '{"username":"reguser","email":"reg@example.com","password":"mypassword"}'
```

Categories, forums and topics are then made through the UI (or with SQL if
you're really bored and don't know what else to do).

TODO: When inserting forums by hand, always set `description`. The column is
nullable, however `models/forum.Description` is a plain `string`, so a `NULL`
breaks the row scan and the forum page renders as an empty root page with HTTP
200, no error, and nothing in the log to suggest what happened. This needs to be
fixed.

The materialized views, which is to say `vforums`, `vtopics` and `vreplies`,
refresh through triggers on insert and update.

## Testing

```sh
make test              # go test -v ./...
go test ./...          # manually, without being verbose
go test -race ./...    # the one that matters for the worker and its caches
go vet ./...
```

### Golden template tests

The outbound message templates are rendered against fixture payloads and
compared against files under `testdata/golden/`. They live in `package main`
(`templates_test.go`).

```sh
go test .                     # verify
go test . -update-golden      # regenerate, after an intentional change
```

`TestTemplatesMatchGolden` renders every `(target, jobtype, subtype, lang)` case
and diffs it, while `TestEveryEmbeddedTemplateHasAGoldenCase` walks the embedded
templates and fails if one of them has no case at all, which is what stops a new
template from slipping in without the suite ever rendering it.

Templates reach into struct fields directly, and `text/template` treats a
missing struct field as a hard execution error, since `<no value>` only ever
happens for maps, which is the reason these tests exist at all.

## Verifying by hand

The unit tests obviously do not cover delivery, routing or permissions end to
end, so what follows is how the message pipeline is actually verified.

### Driving the app with curl

CSRF is enforced on every `POST`, as a double-submit of both a cookie and a form
field, which means you need a cookie jar and a scraped token:

```sh
JAR=/tmp/jar.txt

# scrape the token from a page that definitely has a form
curl -s -c $JAR http://localhost:3000/session/signin -o /tmp/p.html
CSRF=$(grep -o 'name="_csrf"[^>]*value="[^"]*"' /tmp/p.html \
       | grep -o 'value="[^"]*"' | cut -d'"' -f2 | head -1)

curl -s -b $JAR -c $JAR -o /dev/null -w "%{http_code}\n" \
  -X POST http://localhost:3000/session/signin \
  --data-urlencode "_csrf=$CSRF" \
  --data-urlencode "username=admin1" \
  --data-urlencode "password=mypassword"
```

Scrape the token from a page that has a form, because the home page has none and
a form that is hidden by permissions yields an empty token, so the `POST` sends
`_csrf=` and collects a `403` that looks like a permissions problem.
`/session/signin` (while logged out) and `/account/profile` (while logged in)
are both reliable.

Notes:

- Admin settings forms `POST` to `/admin/<section>/update` rather than to
  `/admin/<section>`, which 405s.
- Forms bind the whole struct, so a field you leave out is submitted as its zero
  value, which can trip `required` validations or can blank a setting.
- An unchecked checkbox is simply absent from the `POST`, and that absence is
  how it binds as `false`.
- File uploads want `-F "field=@file;type=mime"` rather than `--data`.

### Driving production-only behaviour

Any other value than `development` for the `Mode` configuration selects the
production behaviour. There, the session and CSRF cookies gain the `Secure`
flag, views and static assets are served out of the copies embedded in the
binary rather than from the working directory, the static cache-buster becomes
the build hash rather than a timestamp, and a global response cache is placed in
front of the routes.

Anything that happens only outside development, such as the `InsecureSkipVerify`
check on `/admin/health`, therefore has to be driven against a config whose
`Mode` says `production`. Luckily `curl` counts `localhost` as a trustworthy
origin and will send `Secure` cookies over plain HTTP to it regardless, so the
recipe above keeps working unchanged:

However, views and static assets are read from the binary in this mode, so a
change to a template needs a `make build` before it appears at all, whereas in
development they are read from the working directory. The response cache,
meanwhile, is keyed on the request path alone and expires on a timer, so when
you restart against an altered config in order to compare one behaviour against
another, an entry cached during the previous run is able to answer the request
and show you the state you were trying to change away from.

### Verifying notification delivery

With the `debug` targets in place, a signup or a reply logs out the
notifications:

```sh
rm -rf /tmp/hyperuplink-messages
# ... drive a signup, or post a reply ...
grep -i "^Target:\|^To:" /tmp/hyperuplink-messages/*
```

That covers the recipients, the channel routing, the rendered subject and body,
and the resulting URLs. _XMPP_ routing can be exercised without a server at all
by giving a user a JID:

```sql
UPDATE users SET email = 'bob@jabber.example.org', email_is_jid = true
 WHERE username = 'bob';
UPDATE settings SET json_value = '{"allowed_address_type": 2}' WHERE id = 'auth';
```

`allowed_address_type` is `0` for e-mail only, `1` for JID only and `2` for
both, the last of which is also what makes the signup form show the JID
checkbox.

### Verifying real XMPP delivery

In order to watch a notification actually arrive you want to log in as the
recipient with a throwaway client, for example in a scratch module that points
back at this repository:

```sh
mkdir -p /tmp/xr && cd /tmp/xr
cat > go.mod <<'EOF'
module xr

go 1.26.4
EOF
go mod edit -require=xn--gckvb8fzb.com/hyperuplink@v0.0.0 \
            -replace=xn--gckvb8fzb.com/hyperuplink=/path/to/hyperuplink
go mod tidy
```

```go
// main.go
package main

import (
	"crypto/tls"
	"fmt"

	goxmpp "github.com/xmppo/go-xmpp"
)

func main() {
	client, err := goxmpp.Options{
		Host: "localhost:5222", User: "postbob@localhost", Password: "bobpass",
		NoTLS: true, StartTLS: true,
		TLSConfig: &tls.Config{ServerName: "localhost", InsecureSkipVerify: true},
		Session: true, Status: "chat",
	}.NewClient()
	if err != nil {
		panic(err)
	}
	defer client.Close()

	fmt.Println("listening ...")
	for {
		stanza, err := client.Recv()
		if err != nil {
			panic(err)
		}
		if chat, ok := stanza.(goxmpp.Chat); ok && chat.Text != "" {
			fmt.Printf("from %s:\n%s\n", chat.Remote, chat.Text)
		}
	}
}
```

Run that, point `comms_xmpp` at the `jabber` target, set the recipient's `email`
to `postbob@localhost` with `email_is_jid = true`, and post a reply to a topic
they have taken part in, whereupon the message turns up in the receiver.

### Verifying common failure paths

#### _XMPP_ server down at boot

Point a target at an unreachable host, and the board should still come up, log
`failed to connect, will retry on first send`, and connect lazily on the first
send once the server is back.

#### Concurrency

`go test -race ./...`, and for the server itself `go build -race`, then fire a
pile of concurrent requests and grep the log for `DATA RACE`.

### Inspecting the queue

```sh
podman exec redis redis-cli --scan --pattern 'asynq:*' | head
podman exec redis redis-cli FLUSHDB          # clear queue state between runs
```

Jobs retry five times before they dead-letter, and a failing target only fails
its own job, because dispatch groups recipients per target and enqueues one job
for each, so an _XMPP_ outage cannot hold up e-mail.

### Inspecting the database

```sh
psql -h localhost -p 5432 -U postgres -d hyperuplink_dev
```

```sql
-- which target serves each channel?
SELECT id, json_value FROM settings WHERE id LIKE 'comms%';

-- per-user settings sit under a composite id: profile-<user uuid>
SELECT id, json_value FROM settings WHERE id LIKE 'profile-%';

-- category permissions, where bits::int is 0/4/6/7 = none/read/read-write/moderate
SELECT group_id, category_id, bits, bits::int AS level FROM permissions
 WHERE deleted_at IS NULL;
```

## Releasing

```sh
make release VERSION=0.1.3    # bumps, commits and tags, but pushes nothing
git push --follow-tags        # pushes
```

`make release` writes the version into `deploy/nix/package.nix`, copies
`hyperuplink-9999.ebuild` to `hyperuplink-0.1.3.ebuild`, commits the two as
`Release v0.1.3` and derives the signed tag from them rather than the other way
around, because a _Nix_ flake cannot read the tag it was built from.

The ebuild needs no editing at all, because _Portage_ takes `${PV}` from the
filename and the live ebuild branches on it, which is what makes the copy a
release ebuild and the original a live one, and is why there is only ever the
one file to keep up to date.

The `make release` command refuses to run when `VERSION` is not an `x.y.z`, when
tracked files have uncommitted changes, or when the tag exists already, and it
rolls both files back if the commit fails, so a half-finished release does not
survive.

Nothing leaves the machine until the push, which then starts the Release, which
cross-compiles the targets in `.goreleaser.yaml`, creates the GitHub release and
attaches the Go module cache that the ebuild needs as a distfile. It also
triggers the container build, which pushes the multi-arch image to `ghcr.io`.

The _RPM_ build waits for _Release_ to succeed first, since _GoReleaser_ creates
the release only once all targets are built and there is nothing to attach the
packages to before that.

`make release` also refuses to publish anything when `deploy/nix/package.nix`
and the tag disagree, which is the backstop for a tag that was made by hand
instead of by `make release`.

> **Note:** The Gentoo _Manifest_ hashes the deps tarball that only exists once
> _Release_ has built and attached it, so it cannot be generated at tag time and
> has to be done afterwards on a machine with _Portage_, as
> `deploy/gentoo/README.md` describes.

## Conventions

A handful of these are not apparent from reading the code:

### Layering DAG, which does not get inverted

`models` and `errs`, then `services`, then `runtime`, then `logic`, then `http`,
`cron`, `worker` and `main`. `logic` must not import `http`, so the _fiber_
`Ctx` and the session stay in the controller and plain values are handed down,
and anything that composes runtime and logic together, such as the CLI flags or
the cron registration, belongs in `main`.

### `logic/root/` mirrors `http/web/root/` one to one

And `logic/helpers/` holds the cross-controller pieces.

### A `v` prefix means materialized view

As in `vtopic` and `vreply`, and those models embed their base model.

### Startup-created singletons are never mutated from request code

Meaning the runtime, the services and the route controllers, since doing so is a
cross-request data leak.

### Migrations are edited in place

At least while the project is pre-release, and `make db:drop` re-runs them, so
there are no `ALTER` migrations yet.

### Settings are seeded in code

... rather than in migrations, in `services/repositories/setting/setting._.go`,
and per-user settings reuse the same table under a `profile-<uuid>` id, created
on demand.
