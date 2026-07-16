# NixOS VM test for the module in module.nix.
#
#   nix build .#checks.x86_64-linux.nixos      # flake
#   nix-build deploy/nix -A tests.nixos        # without flakes
#
# `machine` boots with nothing but `services.hyperuplink.enable = true`.
# PostgreSQL and Redis get provisioned, the migrations run against a real
# database over the peer-authenticated socket, and the board answers on its
# port. `secrets` covers the environmentFile path, where the config is rendered
# again at start-up.
#
# The test framework passes its own `pkgs` down to the nodes and makes
# nixpkgs.config and nixpkgs.overlays read-only there, so the package is handed
# to the module directly rather than through overlay.nix.
{
  pkgs ? import <nixpkgs> { },
  hyperuplink ? pkgs.callPackage ./package.nix { },
  ...
}:

pkgs.testers.runNixOSTest {
  name = "hyperuplink";

  nodes = {
    machine = {
      imports = [ ./module.nix ];

      services.hyperuplink = {
        enable = true;
        package = hyperuplink;
        settings.Users.PromoteAdmin = [ "admin@example.com" ];
      };

      # The VM default is stingy for a Go server next to PostgreSQL.
      virtualisation.memorySize = 2048;
    };

    secrets =
      { pkgs, ... }:
      {
        imports = [ ./module.nix ];

        services.hyperuplink = {
          enable = true;
          package = hyperuplink;
          # A real deployment keeps this out of the store, but the point here is
          # only that the placeholder below gets filled in from it.
          environmentFile = pkgs.writeText "hyperuplink.env" ''
            ADMIN_EMAIL=admin@example.com
          '';
          settings.Users.PromoteAdmin = [ "$ADMIN_EMAIL" ];
        };

        virtualisation.memorySize = 2048;
      };
  };

  testScript = ''
    start_all()

    with subtest("the defaults are enough to serve the board"):
        machine.wait_for_unit("hyperuplink.service")
        machine.wait_for_open_port(3000)

        # A rendered front page means the embedded views, the locales and the
        # migrated schema all came together.
        machine.succeed("curl -sSf http://localhost:3000/ | grep -i hyperuplink")

        # Static assets are served out of the binary, not off the disk.
        machine.succeed("curl -sSf -o /dev/null http://localhost:3000/static/css/hyperuplink.css")

        # The schema really landed in the provisioned database.
        machine.succeed("sudo -u postgres psql -d hyperuplink -c 'SELECT * FROM users LIMIT 1'")

        # Nothing secret to hide, so the config is read straight from the store.
        machine.succeed("systemctl show hyperuplink.service -p ExecStart | grep -F 'file:///nix/store'")

    with subtest("environmentFile placeholders are substituted at start-up"):
        secrets.wait_for_unit("hyperuplink.service")
        secrets.wait_for_open_port(3000)

        config = "/run/hyperuplink/hyperuplink.toml"
        secrets.succeed(f"grep -F 'admin@example.com' {config}")
        secrets.fail(f"grep -F 'ADMIN_EMAIL' {config}")

        # The rendered config is private to the service, unlike the store copy.
        secrets.succeed(f"test $(stat -c %a {config}) = 600")
        secrets.succeed(f"test $(stat -c %U {config}) = hyperuplink")
  '';
}
