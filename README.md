# Hyperuplink

[![SEGV
LICENSE](https://img.shields.io/static/v1?label=SEGV%20LICENSE&message=1.1&labelColor=0060A8&color=ffffff)](https://xn--gckvb8fzb.com/segv/)

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

More info [here](https://xn--gckvb8fzb.com/hyperuplink-discuss-like-its-1998/)
and at [hyperup.link](https://hyperup.link).

## Building

Hyperuplink is actively developed on
[tty.fail](https://tty.fail/hyperuplink/hyperuplink).
[GitHub](https://github.com/hyperuplink/hyperuplink) is a mirror that provides
pre-built binaries.

### Requirements

- [Go](https://go.dev)
- [PostgreSQL](https://www.postgresql.org) server
- [Redis](https://redis.io)-compatible server (Redis, Valkey, KeyDB, ...)
- _Optional:_ S3-compatible object store (e.g. [MinIO](https://min.io))

### From Source

Clone this repository

- from [tty.fail](https://tty.fail/hyperuplink/hyperuplink) (primary):
  ```sh
  $ git clone https://tty.fail/hyperuplink/hyperuplink.git
  ```
- from [GitHub](https://github.com/hyperuplink/hyperuplink) (mirror):
  ```sh
  $ git clone https://github.com/hyperuplink/hyperuplink.git
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

Hyperuplink is configured through a single TOML file. A complete, commented
reference can be found in [`hyperuplink.toml`](hyperuplink.toml). Almost
everything else, like the board name, theme, colorscheme, permissions, groups,
is configured live from the `/admin` panel once the board is running and you've
created your admin account.

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

## Developing

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

Copyright © 2025-2026 [マリウス](https://xn--gckvb8fzb.com)

Hyperuplink is released under Version 1.1 of the
[SEGV License](https://xn--gckvb8fzb.com/segv/), whose full text is included in
the [LICENSE](LICENSE) file. Go read it, there will be a test on it on Monday.
