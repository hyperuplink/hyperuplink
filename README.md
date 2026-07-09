# Hyperuplink

[![SEGV
LICENSE](https://img.shields.io/static/v1?label=SEGV%20LICENSE&message=1.0&labelColor=0060A8&color=ffffff)](https://xn--gckvb8fzb.com/segv/)

[<img src="https://xn--gckvb8fzb.com/images/chatroom.png" width="275">](https://xn--gckvb8fzb.com/contact/)

```
===============================================================================
                       NOW DIALING ...  ::  CARRIER DETECTED
===============================================================================

       █   █ █   █ ████  █████ ████  █   █ ████  █     █████ █   █ █   █
       █   █ █   █ █   █ █     █   █ █   █ █   █ █       █   ██  █ █  █
       █████  ███  ████  ███   ████  █   █ ████  █       █   █ █ █ ███
       █   █   █   █     █     █  █  █   █ █     █       █   █  ██ █  █
       █   █   █   █     █████ █   █  ███  █     █████ █████ █   █ █   █

      ::: A SUPER HIGH SPEED INTERNET BULLETIN BOARD AS SINGLE BINARY :::


      No PHP. No /var/www. No CGI. No FPM. No node_modules. No JavaScript.
     Just one self-contained, statically-linked executable that beams data
            right from your PostgreSQL cluster into beautiful, modern, 
                      server-side rendered HTML5 and CSS.


                        100% JavaScript-Free Browsing
                            Modern HTML5 and CSS
                                Runs Anywhere
                       Zero Loose Files Spilled On Disk
                           Themeable Like It's 1998
                                  And More!


                     Hyperuplink, it's the forum software 
                    your wife told you not to worry about.
                  And it, too, really whips the Llama's ass.

===============================================================================
```

Hyperuplink is a modern internet bulletin board reimagined as a single,
dependency-free binary, supporting PostgreSQL clusters and without any runtime
dependencies, 100% JavaScript free, and based on modern HTML5/CSS.

**More info here:** [hyperup.link](https://hyperup.link)

## Features

- A single static binary, compiled with `make build`. Without runtime
  dependencies, without interpreters, and without bloated frameworks. And
  because it's written purely in Go, it cross-compiles to Linux, macOS, FreeBSD,
  NetBSD, OpenBSD, and probably to whatever runs on your toaster these days.
- PostgreSQL-native, and cluster-friendly. Talks to Postgres over a plain
  connection string via [pgx](https://github.com/jackc/pgx), so it happily
  points at a single server _or_ a full-blown HA cluster. Hot read paths (forum
  and topic listings) are backed by _materialized views_ so the front page stays
  snappy under load. Schema migrations run automatically on startup and they're
  fully embedded, so there are no external migration files to take care of.
- 100% JavaScript-free. Simple and beautiful HTML5 and CSS, without the browser
  having to run a single line of code. And of course, nothing that tracks how
  your cursor gravitates for the third time towards that post on _"why pineapple
  actually belongs on pizza"_.
- Ships with a wardrobe of retro themes, and a few more _modern_ ones for
  everyone who's trying to host a ~~boring~~ serious discussion forum. Oh and
  the colorschemes, those are interchangable. Ever seen a _Gruvbox_ colored
  macOS 9 interface?
- Relatively fine-grained, additive access control with user groups, so you can
  make _the good stuff_ only available to _the good people_.
- Password sign-in with optional TOTP two-factor authentication (compatible with
  any standard authenticator app), plus OAuth sign-in, so you can easily
  convince your friends on other platforms to join your _Taylor Swift fans
  forum_.
- Feeling like email is too _boomer_ for you? No biggie, use XMPP for sign-ups
  (and notifications!) instead.
- Write posts in _Markdown_. Yeah, that's it, what more do you need?
- Profile pictures and post attachments, stored either on the local disk or in
  any S3-compatible object store.
- Moderation features included, so users can report posts and admins can take
  action if necessary.
- Ships with a UI in English and a handful of other languages that I wasn't
  afraid to butcher along the way. And for everything I couldn't translate it
  gracefully falls back to English.
- And there's a lot more where those parlor tricks came from!

## Building

Hyperuplink is actively developed on
[Codeberg](https://codeberg.org/mrus/hyperuplink).
[GitHub](https://github.com/mrusme/hyperuplink) is a mirror that provides
pre-built binaries.

### Requirements

- [Go](https://go.dev) 1.26 or newer
- A [PostgreSQL](https://www.postgresql.org) server
- A [Redis](https://redis.io)-compatible server (Redis, Valkey, KeyDB, ...)
- _Optional:_ an S3-compatible object store (e.g. [MinIO](https://min.io)) if
  you'd rather keep uploads off the local disk

### From Source

Clone this repository

- from [Codeberg](https://codeberg.org/mrus/hyperuplink) (primary):
  ```sh
  $ git clone https://codeberg.org/mrus/hyperuplink.git
  ```
- from [GitHub](https://github.com/mrusme/hyperuplink) (mirror):
  ```sh
  $ git clone https://github.com/mrusme/hyperuplink.git
  ```

Then `cd` into the cloned directory and build:

```sh
$ make build
```

The self-contained binary will be available at `./build/hyperuplink` and can be
moved wherever you please.

You can check the version of a build with:

```sh
$ ./build/hyperuplink -v
```

## Configuration

Hyperuplink is configured through a single TOML file, which you hand to the
binary as a `file://` URI (see [Running](#running)). A minimal configuration
looks roughly like this:

```toml
[General]
Mode = "production"

[Logging]
Level = "info"

[Redis]
Addresses = ["127.0.0.1:6379"]

[Database]
Connection = "postgres://user:pass@localhost:5432/hyperuplink?sslmode=disable"

[Server]
BindIP = "0.0.0.0"
Port   = 3000

[Users]
# The e-mail address(es) here are automatically promoted to admin on sign-up
PromoteAdmin = ["you@example.com"]

# --- Optional: sign in with GitHub ------------------------------------------
[[AuthProvider]]
Type   = "github"
Key    = "your-oauth-app-key"
Secret = "your-oauth-app-secret"
Scopes = ["read:user", "user:email"]

# --- Optional: where uploads live -------------------------------------------
[[Storage]]
ID   = "local-storage"
Type = "Local"
Local.Path      = "/var/lib/hyperuplink/media"
Local.PublicURI = "/media"

# --- Optional: an S3-compatible store instead of / in addition to local -----
[[Storage]]
ID   = "remote-storage"
Type = "S3"
S3.Endpoint  = "http://localhost:9000"
S3.Region    = "us-east-1"
S3.AccessKey = "minioadmin"
S3.SecretKey = "minioadmin"

# --- Optional: email notifications ------------------------------------------
[[Target]]
ID   = "notifications"
Type = "email"
Email.SMTPServer   = "smtp.example.org"
Email.SMTPUsername = "user"
Email.SMTPPassword = "pass"
Email.From.Email   = "reply@example.org"
Email.From.Name    = "Example.org"
```

A complete, commented reference lives in [`hyperuplink.toml`](hyperuplink.toml).
Almost everything else, like the board name, theme, colorscheme, permissions,
groups, is configured live from the `/admin` panel once the board is running and
you've created your admin account (the one you specified in `PromoteAdmin`).

> **Note:** For S3 attachment storage, use a **private** bucket. The gated
> attachment route serves objects with server-side credentials, so a private
> bucket keeps uploads from being guessable at the storage layer.

## Running

Point the binary at your configuration file and launch it:

```sh
$ ./build/hyperuplink -c "file:///etc/hyperuplink.toml"
```

On startup Hyperuplink connects to PostgreSQL, runs any pending migrations,
starts the web server and the background worker, and starts serving. Open your
browser at the address you configured (e.g. `http://localhost:3000`), sign up
with the e-mail you listed under `[Users] PromoteAdmin`, and you'll be set up as
the board's first administrator.

During development you can build and run in a single step:

```sh
$ make run
```

To stop the server, send it a `SIGINT` (Ctrl-C) or `SIGTERM`. This shuts the web
server and worker down cleanly.

## Thanks

A huge, heartfelt thank-you to
[nielssp/classic-stylesheets](https://github.com/nielssp/classic-stylesheets),
whose gorgeous retro CSS made my life _enormously_ easier when starting out with
Hyperuplink and gave the whole thing its unmistakable old-school look. Go star
it.

## License

Hyperuplink is released under the license specified in the [LICENSE](LICENSE)
file. Go read it, there will be a test on it on Monday.
