# Nix package for Hyperuplink

This directory packages Hyperuplink for [Nix][nix], so that it can be deployed,
without Docker or Podman, on NixOS as a service, and on any other distro with
Nix installed as a plain binary.

The flake that ties these together is `flake.nix` in the repository root,
because seemingly that is the only place Nix looks for one.

## Requirements

- Nix 2.4 or newer. Flakes are convenient but optional, see
  [Without flakes](#without-flakes).
- Go, ImageMagick and everything else the build needs come from nixpkgs, so
  nothing has to be installed by hand. The nixpkgs in use does have to carry Go
  1.26.4 or newer, because `go.mod` asks for it and the Go builder pins
  `GOTOOLCHAIN=local`, which forbids fetching another one. The flake pins a
  nixpkgs that does, so this only bites when the overlay is pointed at an older
  channel.
- A reachable PostgreSQL server and a Redis/Valkey server. On NixOS the module
  can provision both, and on other systems they are configured in the TOML
  config like everywhere else.

## License

Hyperuplink is under the custom, non-OSI SEGV license. Section 2.1 grants
copying and distribution, but section 3 conditions the grant on ethical
standards. That is not "free" in the sense nixpkgs means the word, so the
package is marked `free = false` and a stock nixpkgs refuses to build it until
it is told otherwise.

The flake and `default.nix` already exempt this one package, so `nix build` on
either of them just works. Only when the overlay is used against your own
nixpkgs do you have to say so yourself:

```nix
nixpkgs.config.allowUnfreePredicate = pkg: (pkg.pname or null) == "hyperuplink";
```

## Install the binary

On NixOS or on any other distro with Nix:

```sh
# Run it straight from the repository, without installing:
nix run github:hyperuplink/hyperuplink -- -v

# ... or build it and get a ./result/bin/hyperuplink:
nix build github:hyperuplink/hyperuplink

# ... or install it into a profile:
nix profile install github:hyperuplink/hyperuplink
```

The result is the single self-contained executable, wrapped so that
ImageMagick's `convert` is on its PATH. Point it at a config the usual way:

```sh
hyperuplink -c file:///etc/hyperuplink.toml
```

## Deploy on NixOS

Add the flake as an input and import the module:

```nix
{
  inputs.hyperuplink.url = "github:hyperuplink/hyperuplink";

  outputs = { nixpkgs, hyperuplink, ... }: {
    nixosConfigurations.board = nixpkgs.lib.nixosSystem {
      system = "x86_64-linux";
      modules = [
        hyperuplink.nixosModules.default
        {
          services.hyperuplink = {
            enable = true;
            settings.Users.PromoteAdmin = [ "me@example.com" ];
          };
        }
      ];
    };
  };
}
```

`enable = true` also provisions a PostgreSQL server and a Redis server, runs the
migrations on first start, and leaves the board listening on port `3000`.

| Path                            | Contents                                             |
| ------------------------------- | ---------------------------------------------------- |
| `/var/lib/hyperuplink`          | state directory, the service's home                  |
| `/var/lib/hyperuplink/media`    | local media storage (`0750 hyperuplink:hyperuplink`) |
| `/nix/store/…-hyperuplink.toml` | generated config, passed with `-c`                   |

There is no `/etc/hyperuplink.toml`. The config is generated from
`services.hyperuplink.settings` and passed to the binary out of the store, which
is what makes the deployment reproducible. Everything in
`deploy/hyperuplink.toml` is settable there, with the same key names:

```nix
services.hyperuplink.settings = {
  Web.ProxyHeader = "X-Forwarded-For";   # behind a reverse proxy
  Web.TrustProxy = true;
  Logging.Level = "debug";
};
```

The only mapping that is not literal is `[[Storage]]`, which is a TOML array of
tables and so a list of attribute sets here. The default is local storage under
the state directory, and replacing the list is how an S3 provider is added:

```nix
services.hyperuplink.settings.Storage = [
  {
    ID = "local-storage";
    Type = "Local";
    Local = { Path = "/var/lib/hyperuplink/media"; PublicURI = "/media"; };
  }
  {
    ID = "remote-storage";
    Type = "S3";
    S3 = {
      Endpoint = "https://s3.example.com";
      Region = "us-east-1";
      AccessKey = "$S3_ACCESS_KEY";        # see Secrets below
      SecretKey = "$S3_SECRET_KEY";
      PublicURL = "https://s3.example.com";
      PublicDownload = true;
    };
  }
];
```

Run it behind a TLS-terminating reverse proxy, and leave `General.Mode` at
`"production"` so the session cookie stays HTTPS-only:

```nix
services.nginx.virtualHosts."board.example.com" = {
  enableACME = true;
  forceSSL = true;
  locations."/".proxyPass = "http://127.0.0.1:3000";
};
```

### Database and Redis

By default the module provisions both, and connects to PostgreSQL over its Unix
socket as the `hyperuplink` role, which peer authentication accepts with no
password at all. That is why the default setup has no secrets to keep.

To use servers you run elsewhere, turn the provisioning off and say where they
are:

```nix
services.hyperuplink = {
  database.createLocally = false;
  redis.createLocally = false;
  settings = {
    Database.Connection = "postgres://hyperuplink@db.example.com:5432/hyperuplink?sslmode=require";
    Redis.Addresses = [ "cache.example.com:6379" ];
  };
};
```

The job queue speaks TCP only, so a Redis Unix socket is not an option, and
`redis.port` exists for when 6379 is already taken on the host.

### Secrets

`settings` is written to the Nix store, which is world-readable, so a password
put there is readable by every user on the machine. Keep it in an environment
file instead, and reference it as a placeholder:

```nix
services.hyperuplink = {
  environmentFile = "/var/lib/secrets/hyperuplink.env";  # DB_PASSWORD=hunter2
  settings.Database.Connection =
    "postgres://hyperuplink:$DB_PASSWORD@db.example.com:5432/hyperuplink";
};
```

With `environmentFile` set, the config is rendered again at start-up into
`/run/hyperuplink/hyperuplink.toml` (mode `0600`, private to the service) with
the placeholders filled in. It also makes every literal `$` in `settings`
significant, so write `$$` for a dollar sign that should survive.

## Without flakes

`default.nix` is the same packaging without the flake, for a system on channels:

```sh
nix-build deploy/nix -A hyperuplink
./result/bin/hyperuplink -v
```

For the NixOS module, import it by path and add the overlay so that
`pkgs.hyperuplink` exists:

```nix
{
  imports = [ /path/to/hyperuplink/deploy/nix/module.nix ];
  nixpkgs.overlays = [ (import /path/to/hyperuplink/deploy/nix/overlay.nix) ];
  nixpkgs.config.allowUnfreePredicate = pkg: (pkg.pname or null) == "hyperuplink";
  services.hyperuplink.enable = true;
}
```

## Testing

`test.nix` boots a NixOS VM with nothing but
`services.hyperuplink.enable = true`, waits for the service, and checks that it
answers, that the embedded assets are served and that the migrations reached the
provisioned database.

```sh
nix build .#checks.x86_64-linux.nixos      # needs KVM
nix-build deploy/nix -A tests.nixos        # without flakes

nix flake check                            # test and package together

# Evaluation only, which is how the arches without a runner are covered:
nix flake check --all-systems --no-build
```

`.github/workflows/nix.yml` runs all of that on changes to `deploy/nix/**`, the
flake, or `go.mod`/`go.sum`. The last two are in the path filter because
`vendorHash` is derived from `go.sum`, so a dependency bump invalidates it and
breaks every install from the flake. There is no release job, because unlike the
RPM or the container there is no artifact to attach, because a tag is installed
directly with `nix build github:hyperuplink/hyperuplink/vX.Y.Z`.

## Development shell

`nix develop` drops into a shell with the Go toolchain the build uses, plus
ImageMagick and the PostgreSQL and Valkey clients, which is enough to
`make build` and to talk to the local stack described in `DEVELOPMENT.md`. It
does not start any of the services, so the containers there are still the way to
run them.

## Notes

- The version, commit and build date are injected with `-ldflags -X` exactly as
  in the ebuild and the spec, so `hyperuplink -v` reports them. The flake fills
  in the revision it was built from.
- A flake cannot read the git tag. The version therefore has to be committed in
  the tree, because `nix build github:hyperuplink/hyperuplink/v0.1.3` has
  nothing else to read it from. So rather than keep the file and the tag in step
  by hand:

  ```sh
  make release VERSION=0.1.3
  ```

  That writes the version into `package.nix`, commits it, and tags `v0.1.3` from
  it, but it pushes nothing. `release.yml` still checks the tag against the
  file, which is what catches a tag made by hand.
- `package.nix` pins `vendorHash`, which is the hash of the Go module cache, and
  it changes whenever `go.sum` does. To refresh it, set it to `lib.fakeHash`,
  build, and copy the hash out of the mismatch the build reports.
- The Go test suite talks to a live PostgreSQL and Redis, which the build
  sandbox has not got, so `doCheck` is off and the VM test covers the package
  instead.
- `src` is a filtered file set listing only what the compiler and `go:embed`
  read, so editing this packaging, the docs or the other `deploy/` formats does
  not rebuild the binary.

[nix]: https://nixos.org/
