# Development

This document describes the process for running this application on your local
computer.

## Getting started

Everything except the following is embedded into the binary, which means that
templates, views, static assets, locales, migrations and docs are all contained
in the executable and there are no loose files to manage:

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
includes the version, commit and build date in the binary through `-ldflags`,
and because it regenerates the [OpenAPI specification](#openapi-specification)
the API embeds.

```sh
make build         # -> ./build/hyperuplink
make help          # every target
```

## Local stack

_PostgreSQL_ and a _Redis_-compatible server are the only two pieces that are
needed in order to boot, while _MinIO_ is only ever touched once a file is
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

The job queue (_asynq_) and the session store (_fiber_) both use this server,
and any _Redis_-compatible server will do.

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

_Hyperuplink_ requires _StartTLS_ and will not authenticate over an unencrypted
connection, which means that _Prosody_ needs a certificate even on loopback,
though a self-signed one is sufficient as long as the target sets
`XMPP.InsecureSkipVerify = true`.

> **Warning:** `XMPP.InsecureSkipVerify = true` turns off certificate
> verification, and it has no business anywhere other than a local development
> server. With verification disabled, anyone able to intercept traffic between
> the two ends of the connection may present whatever certificate they like and
> then read or rewrite the traffic, credentials included.

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

The board then answers on <http://localhost:3000>, the [JSON API](#json-api) on
<http://localhost:3001>, and since both `SIGINT` and `SIGTERM` are trapped,
either `Ctrl-C` or a plain `kill <pid>`/`killall hyperuplink` will shut it down
cleanly.

> **Note:** `-c` takes a URL rather than a path, hence the `file://`, and it
> defaults to `/etc/hyperuplink.toml`.

> **Note:** if something is already bound to `:3000`, the instance will still
> start, but its web server will simply fail to bind without saying much about
> it, at which point your requests continue to land on whatever is listening on
> the port. Check with `ss -ltn | grep :3000` before concluding that anything
> else is wrong.

TODO: development mode watches `views/` and reloads the templates as you save
them, however `getWatcher` passes the directory straight to
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
the current system time as confirmation, in 24 hour `HH:MM` format, and aborts
if the minute does not match:

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

The fields are those of `logic/root/session.SignUpInput`, namely `username`,
`email` and `password`, plus the optional `password_repeat`, `email_is_jid` and
`language`.

For an admin, use an address that is listed under `[Users] PromoteAdmin` in
`hyperuplink.toml`, where `dummy1@example.com` is already listed, since signup
promotes those addresses on its own.

## JSON API

Beside the web server _Hyperuplink_ starts a second _fiber_ app that speaks JSON
on <http://localhost:3001>, configured through an `[API]` section in
`hyperuplink.toml` whose keys mirror those of `[Web]`, though the API only comes
up when that section has `Enable = true`, and because `Enable` defaults to
`false` a missing `[API]` section leaves the JSON API switched off entirely,
while its port still falls back to `3001` whenever it is enabled without an
explicit `Port`.

### Authentication

Every request includes an API key, either as `Authorization: Bearer <key>` or as
`X-API-Key: <key>`, and there are no cookies, no server side sessions and no
CSRF anywhere on this port. A key is a `hup_`-prefixed random secret whose
SHA-256 is stored in the `apikeys` table, which keeps nothing besides that hash,
so a lost key is replaced rather than recovered. There is no route that signs
up, signs in or signs out, and a request without a valid key answers `401`
across the board, guest reads included.

Keys are issued at `/account/api` on the web server, which lists a user's keys,
issues new ones and soft-deletes old ones, and which shows the secret exactly
once, in the flash message right after issuance, since only the hash survives.
When you would rather not click through the UI, for example on a freshly seeded
board you are about to script against, mint one by hand:

```sh
SECRET="hup_$(head -c 32 /dev/urandom | basenc --base64url | tr -d '=')"
HASH=$(printf '%s' "$SECRET" | sha256sum | cut -d' ' -f1)
psql -h localhost -p 5432 -U postgres -d hyperuplink_dev \
  -c "INSERT INTO apikeys (user_id, name, secret_hash)
      SELECT id, 'dev', '$HASH' FROM users WHERE username = 'sysop';"
echo "$SECRET"   # the part the client keeps; the database never sees it
```

```sh
curl -s -H "Authorization: Bearer $SECRET" http://localhost:3001/session
```

The key resolves to its user on every request, so the caller's role and
permissions are the same ones the web session would have produced, and a banned
or deleted user takes their keys down with them.

### Routes

`http/api/root/` mirrors `http/web/root/`, and both are built on the same
`logic/root/` underneath, so an action behaves identically no matter which
server carried it in. Reads are `GET`, the forms that `POST` to `/update` on the
web side are `PUT` here, creation is `POST` and destruction is `DELETE`, with
the record's id moving into the path where the web form carried it as a hidden
field. The paths themselves come out of the very same `http/routes` table the
web server reads, which every controller consults through
`route.For("<id>").Pathname()` rather than spelling its segment out as a
literal, so the API inherits the sigils the web URLs use, a category reading
`/_general` and a profile `/~sysop`, and a segment that no web page has yet, of
which the flat topics feed is so far the only one, gets its own entry in
`http/routes/routes.go` rather than a string in the controller:

| Area         | Routes                                                                                      |
| ------------ | ------------------------------------------------------------------------------------------- |
| board        | `GET /`, `GET /_:cat`, `GET /_:cat/:forum[/:topic]` (`?page=`)                              |
| posting      | `GET`/`POST /new`, `POST /_:cat/:forum/:topic` (reply), `POST .../poll` (vote)              |
| topics feed  | `GET /topics` (`?forum_id=&page=`), `GET /topics/:id`, `POST /topics/:id/replies`           |
| session      | `GET /session`, which answers with the key's user                                           |
| account      | `GET /account`, `PUT /account/{password,profile,settings}`, `POST /account/profile/picture` |
| two-factor   | `GET /account/twofactor`, `POST /account/twofactor/{enable,disable}`                        |
| attachments  | `POST /attachments` (multipart), `GET /attachments/:id`                                     |
| search, docs | `GET /search?q=`, `GET /docs/{about,contact,privacy,terms}`                                 |
| report       | `GET /report?target=&id=`, `POST /report`                                                   |
| users        | `GET /~:username`, `POST /~:username/membership`                                            |
| admin        | `GET`/`PUT` per section under `/admin/...`, users and board records through path ids        |

Everything the web offers is here, except signup, signup confirmation, signin,
signout and the embedded manual. The request bodies take the same fields the web
forms take, however as JSON, and they run through the same struct validation,
which the shared input structs under `logic/root/` declare as `form`, `json` and
`validate` tags side by side.

Three flows differ from their web counterparts, in each case because there is no
session:

- Attachments are not part of the create request. Upload them first as multipart
  `attachments=@file` to `POST /attachments`, then pass the returned ids as
  `attachment_ids` when creating the topic or reply, whereupon the server
  verifies that each id belongs to the caller.
- Two-factor enrollment cannot store the pending secret in a session, so
  `GET /account/twofactor` returns a fresh enrollment whose `URL` the client
  holds on to and sends back as `otp_url`, together with the `otp_code`, to
  `POST /account/twofactor/enable`.
- The profile picture is sent as multipart to `/account/profile/picture` rather
  than inside the JSON `PUT`.

### Errors

Errors come back as `{"error": "<key>"}`, where the key is the same `errs`
string the web shows through its flash messages, and validation failures
additionally name their fields:

```json
{ "error": "validation", "fields": { "text": "validation_text_required" } }
```

The API returns `401` for a missing or unknown key, `403` for a role the route
does not accept, `404` for anything that does not exist or that the caller may
not see, `409` for unique violations such as voting twice, and `422` for
validation and other rejected input.

### OpenAPI specification

The API describes itself through a _Swagger_ UI at
<http://localhost:3001/_internal/swagger/>, which reads the specification from
`/_internal/swagger/doc.json` beside it, and since both live under the
`/_internal/` prefix that `srv.authenticate` waves through, they are reachable
without a key, though every operation the UI lists still needs one the moment
you press _Try it out_.

The specification is generated from annotations by
[swaggo/swag](https://github.com/swaggo/swag), which the _Makefile_ pins as a Go
tool so that nothing has to be installed by hand, and which writes
`docs/swagger.json` under the `swagger` target that `make`, `make build` and
`make run` all depend on:

```sh
make swagger       # -> docs/swagger.json
go tool swag fmt -d http/api -g api._.go   # realign the annotation comments
```

The general information, meaning the title, the description, the license and the
two API key schemes, is declared above `New` in `http/api/api._.go`, while each
routed action has its own `@Summary`, `@Param`, `@Success` and `@Router` block,
and the `@Router` paths are written exactly as the route table produces them,
sigils and all, so a category reads `/_{categories}` and a user reads
`/~{user}`. Types are named the way the file that references them imports them,
which is why the annotations say `logictopics.View` rather than `topics.View`,
however _swag_ also resolves a bare package name against everything it has
parsed, so `setting.Profiles` works from a file that never imports it.

`docs/` is embedded whole, which is how `loadSwagger` finds the specification
without a generated `docs.go` to import, and it means the file has to exist
before the binary is compiled. `make build` takes care of that here, however the
_Dockerfile_, the ebuild and the RPM spec all call `go build` straight, so
`docs/swagger.json` belongs in the repository rather than in `.gitignore`, and a
build that is missing it logs a warning and leaves the endpoint unregistered
instead of failing. The version is the one exception to the file being served
verbatim: it is written as the `{{.Version}}` placeholder that `swag.Spec` fills
in at startup from the same build version the `-ldflags` set, so the UI shows
`v0.1.2-15-g81df9f4` rather than whatever was true when the specification was
last regenerated.

Editing the annotations has two traps. Anything following a
`@securityDefinitions` block in the general information is swallowed by it, so
the `@tag` declarations have to come before the two schemes rather than after
them, and _swag_ hard-skips every directory named `docs` while it walks the
tree, which is why the four `/docs/*` static pages are the only routes the
specification does not describe.

## Configuration

`hyperuplink.toml` in the repository root is the development config.

- `[General] Mode = "development"` enables debug logging and drops the `Secure`
  flag from the session and CSRF cookies, which is what allows plain-HTTP `curl`
  to round-trip them, while production mode serves the embedded views and forces
  HTTPS-only cookies.
- `[Web] Enable` and `[API] Enable` each gate whether that server boots and both
  default to `false`, and since a missing section reads as `false` just the
  same, dropping either `[Web]` or `[API]` from the config leaves the
  corresponding server switched off rather than running on its defaults.
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
so what is written to disk is what would have gone out, one file per recipient
with mode `0600` inside a `0700` directory, since the messages contain signup
tokens.

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
database delivers without anyone having to click through the admin UI first.
Production declares no `debug` target, so nothing falls back there.

## Seeding a board

`Startup()` seeds the _settings_ rows and nothing besides, meaning no users, no
categories, no forums and no topics, so a freshly dropped database is an empty
board that nobody is able to log into. Since migrations are edited in place,
`db:drop` comes round often enough that putting the board back is one command:

```sh
make db:drop && make db:seed
```

That leaves four accounts, all with the password `hyperhyper!`:

| User     | Address                     | Notes                               |
| -------- | --------------------------- | ----------------------------------- |
| `sysop`  | `dummy1@example.com`        | admin, through `PromoteAdmin`       |
| `vera`   | `vera@example.org`          |                                     |
| `juno`   | `juno@example.org`          |                                     |
| `bitrot` | `bitrot@jabber.example.org` | a _JID_, for the XMPP delivery path |

It also leaves two categories, three forums, five topics of which one is a poll
that has been voted in and one has been reported, four replies, and two groups
with a category permission between them, which is enough for every page worth
looking at to have something on it.

### Seed

`tools/seed/seed.sh <config-url> <database>` is the whole of it, and
`make db:seed` is that script pointed at `hyperuplink.toml` and
`hyperuplink_dev`. The screenshot harness calls the same script and then lays
its own theme and port over the top, which is why the manual's screenshots and
the development board are the same board.

Each step is there because the one before it made it possible:

1. `sysop` is created first, since `-create-user` is what runs the migrations
   and seeds the settings rows, and there are no settings to change until
   something has done that.
2. `tools/seed/settings.sql` then runs, and it has to run _here_, because it
   sets the address type to _Email and JID_ and that is what decides whether
   `email_is_jid` means anything at all. Under the default "email only" a JID is
   validated as an email address, which it passes, and `bitrot` would be created
   as an email user without a single error to warn you.
3. The rest of the users are created.
4. `tools/seed/content.sql` runs last, because every topic and reply in it
   points at a user by username.

Neither SQL file sets the theme or the base URL: the first is taste and the
second depends on where the board is running, so both belong to whoever is
seeding. The default base URL is already `http://localhost:3000`, which is what
the development board needs anyway.

To seed something else by hand, `-create-user` still takes the signup form as
JSON, as described under [`-create-user`](#-create-user).

TODO: When inserting forums by hand, always set `description`. The column is
nullable, however `models/forum.Description` is a plain `string`, so a `NULL`
breaks the row scan and the forum page renders as an empty root page with HTTP
200, no error, and nothing in the log to suggest what happened. This needs to be
fixed.

The materialized views `vforums`, `vtopics` and `vreplies` refresh through
triggers on insert and update.

## Manual

The user and administrator manual lives in `docs/manual/`, it is plain Markdown,
one `index.md` per directory, and the directories loosely mirror
`http/web/root/`. It is served by `http/web/root/docs/manual/`, which reads the
page out of the `docs` embed, expands it as a template, converts it and renders
it into `views/docs/manual/index.html`.

**The pages are read from the embed and never from disk, so editing anything
under `docs/manual/` does nothing at all until you `make build`.** Development
mode serves the views from the filesystem, which makes it tempting to assume the
same is true here, and it is not.

### Linking from a page

Every page has to work whether the board sits at `mysite.com/` or at
`mysite.com/hyperuplink/`, so no link in the manual may be absolute. The
Markdown is run through `text/template` before it is converted, with two
functions in scope, both of which count the depth of the page being rendered and
prefix the right number of `../`:

| Function                       | Links to                   |
| ------------------------------ | -------------------------- |
| `{{ hrefTo "admin/general" }}` | a page of the board        |
| `{{ manual "admin/general" }}` | another page of the manual |

`hrefTo` is `Site.HrefTo`, so an absolute URL passed to it is returned untouched
and external links go through it unharmed. Images are plain relative Markdown
links, since a screenshot sits in the same directory as the page that shows it,
and the route serves any non-`.md` file straight out of the embed.

A malformed action, meaning a `{{` that does not parse, fails the page with a
500 rather than at build time, so it is worth clicking through a page you have
just written.

### Screenshots

```sh
make manual:screenshots
make site:screenshots
```

Both drop and recreate `hyperuplink_manual`, seed it through the shared
[`tools/seed/seed.sh`](#seeding-a-board) and then lay
`tools/screenshots/seed.sql` over the top, which sets the _macos9_ theme with
the _hyperuplink-light_ colorscheme that the manual shows throughout and the
port this board answers on. They then start the board against
`tools/screenshots/hyperuplink.toml` on port 3001 and Redis database 3, drive a
headless Chrome over the pages named in `tools/screenshots/shots.go` and write
the `.webp` files back out. Nothing either of them does touches
`hyperuplink_dev`.

Sharing the seed with the development board is deliberate: a screenshot of a
board nobody develops against goes stale without anyone noticing.

`tools/screenshots/` is a **separate Go module**, so `chromedp` stays out of the
board's own `go.mod` and off the dependency list of the binary you ship.

It needs a Chrome or Chromium on `PATH`, and it falls back to the one Playwright
keeps in `~/.cache/ms-playwright/` where there is none, or takes `-chrome` or
`$CHROME`. It also needs ImageMagick, which the board needs anyway for profile
pictures.

#### Two sets

`shots.go` holds a `set` per place the screenshots end up, and the set defines
the framing every shot in it is taken with, so the two share nothing beyond the
board they are pointed at:

| Set      | Written into                 | Framing                                                  |
| -------- | ---------------------------- | -------------------------------------------------------- |
| `manual` | `docs/manual/`               | clipped to `.container` at 1280x860, quality 55          |
| `site`   | `../pub/static/screenshots/` | 1280x800 below `.header`, resized to 960x600, quality 80 |

The `manual` set is clipped to `.container`, which drops the tiled desktop
background and cuts the files by more than half, while the `site` set takes a
fixed 1280x800 region and resizes it, since those four shots are shown as a grid
of tiles on the [hyperup.link](https://hyperup.link) landing page, where one
shared aspect ratio keeps the captions on a line rather than each tile ending
wherever its page happens to.

That region starts under whatever `Below` names, which for the `site` set is
`.header`, the div holding the banner, the menu bar and the breadcrumbs. Shot
from the top of the page instead, every tile would spend a third of its height
on the same banner and the reader would have four pictures of a banner and very
little board, while cutting under it leads each tile with its own title bar. It
is a crop rather than anything hidden in CSS, so the tile is a region of the
page as it really renders.

Two things affect where the crop ends up:

- The window is not the viewport. Headless Chrome counts its own tab strip
  against `--window-size`, so a window of 1280x800 gives you a viewport of
  1280x657, and an aspect-preserving resize then writes 960x493 tiles that no
  longer line up, without warning. Anything not clipped to an element goes
  through `EmulateViewport` for that reason.
- The page has to be tall enough to crop. `reset.css` gives `body` a
  `min-height` of `100vh` and the desktop grows into it, so the document is only
  as tall as the viewport once the content runs out, and a region that reaches
  past the bottom of it comes back with the tail painted black. `fitBelow`
  therefore measures the element and sizes the viewport to that plus the height
  of the region, which puts the whole crop inside a page that is actually
  painted.

The website is a Hugo project of its own beside this repository, so
`make site:screenshots` takes `SITE=`, defaulting to
`../pub/static/screenshots`, and exits with an error rather than creating a
directory somewhere surprising when that path is not there. The four shots it
writes are the ones `../pub/content/_index.md` names in its `[[screenshots]]`
blocks, and the file names are the contract between the two repositories: rename
one here and the landing page shows a broken image, so change both together.

#### Adding and iterating

Add a page by adding a `shot` to the set it belongs in. `-only <substring>`
shoots one while you are iterating on it, `-set <name>` shoots one set, and the
framing flags (`-out`, `-width`, `-height`, `-clip`, `-resize`, `-quality`)
override the set's own values only when you actually pass them, so
`tools/screenshots/run.sh -set site -only manual -quality 100` reshoots a single
tile without touching how the rest are taken.

A shot can override its set's `Clip` and `Below` where one page needs framing of
its own. The manual tile does: the manual is a page about the board with a
screenshot of the board in it, so cropping it under `.header` like the rest
would make the tile a picture of a picture, and it sits under `.markdown img`
instead, which starts it at the prose and the chapter list and makes it look
like a manual.

`As` picks which seeded user is signed in, which is what the poll needs, since a
poll shows its results only to someone who has already voted, and `Prep` runs
`chromedp` actions before the shutter, which is how the poll editor gets opened.

The users it signs in as are the seeded ones, which are listed under
[Seeding a board](#seeding-a-board).

## Testing

```sh
make test              # go test -v ./...
go test ./...          # manually, without being verbose
go test -race ./...    # the one that matters for the worker and its caches
go vet ./...
```

### Golden template tests

The outbound message templates are rendered against fixture payloads and
compared against files under `testdata/golden/`. They are in `package main`
(`templates_test.go`).

```sh
go test .                     # verify
go test . -update-golden      # regenerate, after an intentional change
```

`TestTemplatesMatchGolden` renders every `(target, jobtype, subtype, lang)` case
and diffs it, while `TestEveryEmbeddedTemplateHasAGoldenCase` iterates over the
embedded templates and fails if one of them has no case at all, which is what
stops a new template from being added without the suite ever rendering it.

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
a form that is hidden by permissions produces an empty token, so the `POST`
sends `_csrf=` and gets back a `403` that looks like a permissions problem.
`/session/signin` (while logged out) and `/account/profile` (while logged in)
are both reliable.

Notes:

- Admin settings forms `POST` to `/admin/<section>/update` rather than to
  `/admin/<section>`, which 405s.
- Forms bind the whole struct, so a field you leave out is submitted as its zero
  value, which can trip `required` validations or can blank a setting.
- An unchecked checkbox is simply absent from the `POST`, and that absence is
  how it binds as `false`.
- File uploads need `-F "field=@file;type=mime"` rather than `--data`.

### Driving production-only behaviour

Any other value than `development` for the `Mode` configuration selects the
production behaviour. There, the session and CSRF cookies gain the `Secure`
flag, views and static assets are served out of the copies embedded in the
binary rather than from the working directory, the static cache-buster becomes
the build hash rather than a timestamp, and a global response cache is placed in
front of the routes.

Anything that happens only outside development, such as the `InsecureSkipVerify`
check on `/admin/health`, therefore has to be driven against a config whose
`Mode` is `production`. Luckily `curl` counts `localhost` as a trustworthy
origin and will send `Secure` cookies over plain HTTP to it regardless, so the
recipe above keeps working unchanged:

However, views and static assets are read from the binary in this mode, so a
change to a template needs a `make build` before it appears at all, while in
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

In order to watch a notification actually arrive you need to log in as the
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

-- per-user settings are stored under a composite id: profile-<user uuid>
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

The `make release` command exits with an error when `VERSION` is not an `x.y.z`,
when tracked files have uncommitted changes, or when the tag exists already, and
it rolls both files back if the commit fails, so a half-finished release is not
left behind.

Nothing leaves the machine until the push, which then starts the Release, which
cross-compiles the targets in `.goreleaser.yaml`, creates the GitHub release and
attaches the Go module cache that the ebuild needs as a distfile. It also
triggers the container build, which pushes the multi-arch image to `ghcr.io`.

The _RPM_ build waits for _Release_ to succeed first, since _GoReleaser_ creates
the release only once all targets are built and there is nothing to attach the
packages to before that.

`make release` also does not publish anything when `deploy/nix/package.nix` and
the tag disagree, which is the backstop for a tag that was made by hand instead
of by `make release`.

> **Note:** The Gentoo _Manifest_ hashes the deps tarball that only exists once
> _Release_ has built and attached it, so it cannot be generated at tag time and
> has to be done afterwards on a machine with _Portage_, as
> `deploy/gentoo/README.md` describes.

## Conventions

A handful of these are not apparent from reading the code:

### Layering DAG

`models` and `errs`, then `services`, then `runtime`, then `logic`, then `http`,
`cron`, `worker` and `main`. `logic` must not import `http`, so the _fiber_
`Ctx` and the session stay in the controller and plain values are handed down,
and anything that composes runtime and logic together, such as the CLI flags or
the cron registration, belongs in `main`.

### Framework

The lower half of that stack consists of
[_glides_](https://tty.fail/mrus/glides), which is a framework that came out of
_Hyperuplink_ and serves as a foundation for it, as well as for other software.

`main` builds a `*glides/runtime.Runtime` with the framework services it asks
for in `runtime.Services` and then registers the board's own (e.g.
`repositories`, `activity`, `magick` and `dispatch`) through `rt.AddService`,
which the typed getters in `helpers` resolve back. `http/route` holds the
routing table and binds the generic controller types to the framework runtime.

Board-specific configuration has no getter of its own, because the framework
only declares the keys it reads itself. `Users.PromoteAdmin` is read with
`rt.Config().Strings()` and the `AuthProvider` blocks with
`rt.Config().Unmarshal()`, into the `AuthProviders` type that is in `helpers`.

### `logic/root/` mirrors `http/web/root/`

And `http/api/root/` mirrors both, so one action is three thin packages on the
same path: the use-case in `logic/root/`, which takes plain values and a
`*runtime.Runtime`, and one controller per server that binds, authenticates and
answers in its own format. The input structs are defined with the logic and
declare `form`, `json` and `validate` tags together, which is what keeps the two
servers validating identically. `logic/helpers/` holds the cross-controller
pieces, among them the permission resolver, the pagination arithmetic and the
attachment storage pipeline.

### A `v` prefix means materialized view

As in `vtopic` and `vreply`, and those models embed their base model.

The refresh triggers fire `AFTER INSERT OR UPDATE` on every table a view reads
from, and the `OR UPDATE` half is what keeps the soft deletes honest, since a
soft delete is an `UPDATE` and a view that only refreshed on `INSERT` would
continue serving a row that had already been marked as deleted.

### Soft deletes propagate downwards

Nothing is ever dropped: `deleted_at` is stamped, `common.QueryOptions` appends
`deleted_at IS NULL` to every read, and each read filters **its own** table.
That last part is why the delete has to walk the tree itself, which
`Category.Delete` and `Forum.Delete` do in one transaction, marking the forums,
the topics under them and the replies under those, all guarded by
`WHERE deleted_at IS NULL` so an earlier deletion keeps its own timestamp.

Leave a child unmarked and it stays visible, because its own `deleted_at` is
still `NULL` and no read path walks up the tree to check its parent. The one
place that does walk up is the reply branch of the search query, which joins up
to the topic, the forum and the category, and it is the exception rather than
the pattern.

A lookup that comes back `errs.ErrNoRows` is a 404 in the category, forum and
topic show controllers, so a deleted item is gone for anyone who kept the URL.

### Startup singletons are never mutated

Meaning the runtime, the services and the route controllers, since doing so is a
cross-request data leak.

### Migrations are edited in place

At least while the project is pre-release, and `make db:drop` re-runs them, so
there are no `ALTER` migrations yet.

### Settings are seeded in code

... rather than in migrations, in `services/repositories/setting/setting._.go`,
and per-user settings reuse the same table under a `profile-<uuid>` id, created
on demand.
