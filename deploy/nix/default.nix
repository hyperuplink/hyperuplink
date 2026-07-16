# Non-flake entry point for the Nix packaging.
#
# Works with a plain nixpkgs on NixOS and on any other distro with Nix
# installed, so flakes are not needed:
#
#   nix-build deploy/nix -A hyperuplink
#   nix-env  -f deploy/nix -iA hyperuplink
#
# SEGV is not a free license in the sense nixpkgs means it (see the license
# comment in package.nix), so a stock nixpkgs refuses to build the package. The
# default below exempts this one package rather than allowing unfree globally,
# and is also the snippet to copy into a system configuration.
{
  pkgs ? import <nixpkgs> {
    config.allowUnfreePredicate = pkg: (pkg.pname or null) == "hyperuplink";
  },
}:
rec {
  hyperuplink = pkgs.callPackage ./package.nix { };

  # NixOS VM test, `nix-build deploy/nix -A tests.nixos`.
  tests.nixos = pkgs.callPackage ./test.nix { inherit pkgs hyperuplink; };
}
