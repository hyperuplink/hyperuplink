# Flake entry point. The packaging itself lives in deploy/nix.
#
#   nix run   github:hyperuplink/hyperuplink -- -v
#   nix build github:hyperuplink/hyperuplink
#
# See deploy/nix/README.md for deploying it on NixOS.
{
  description = "A super high speed internet bulletin board";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs =
    { self, nixpkgs }:
    let
      # Linux is what the service targets, and the arches match the ebuild's
      # KEYWORDS. Darwin is here so that `nix run` and the dev shell work on a
      # Mac.
      systems = [
        "x86_64-linux"
        "aarch64-linux"
        "riscv64-linux"
        "aarch64-darwin"
      ];

      # The dev shell and the formatter are conveniences, and nixfmt is a
      # Haskell program that GHC cannot bootstrap on riscv64.
      # The package itself builds fine on it.
      devSystems = builtins.filter (system: system != "riscv64-linux") systems;

      forAllSystems = nixpkgs.lib.genAttrs systems;
      forDevSystems = nixpkgs.lib.genAttrs devSystems;

      # SEGV is not a free license in the sense nixpkgs means it (see the
      # license comment in deploy/nix/package.nix), so a stock nixpkgs refuses
      # to build the package. This exempts this one package rather than allowing
      # unfree.
      pkgsFor =
        system:
        import nixpkgs {
          inherit system;
          config.allowUnfreePredicate = pkg: (pkg.pname or null) == "hyperuplink";
        };

      # `nix build` from a checkout has no tag to read, so the version stays the
      # one in package.nix and only the revision is filled in.
      rev = self.shortRev or self.dirtyShortRev or "unknown";
      date = self.lastModifiedDate or "unknown";
    in
    {
      packages = forAllSystems (
        system:
        let
          pkgs = pkgsFor system;
        in
        rec {
          hyperuplink = pkgs.callPackage ./deploy/nix/package.nix { inherit rev date; };
          default = hyperuplink;
        }
      );

      overlays.default = import ./deploy/nix/overlay.nix;

      # Defaults services.hyperuplink.package to this flake's build.
      nixosModules.hyperuplink =
        { pkgs, lib, ... }:
        {
          imports = [ ./deploy/nix/module.nix ];
          services.hyperuplink.package =
            lib.mkDefault
              self.packages.${pkgs.stdenv.hostPlatform.system}.hyperuplink;
        };
      nixosModules.default = self.nixosModules.hyperuplink;

      checks = forAllSystems (
        system:
        let
          pkgs = pkgsFor system;
        in
        {
          package = self.packages.${system}.hyperuplink;
        }
        // nixpkgs.lib.optionalAttrs (pkgs.stdenv.hostPlatform.isLinux) {
          nixos = pkgs.callPackage ./deploy/nix/test.nix {
            inherit pkgs;
            hyperuplink = self.packages.${system}.hyperuplink;
          };
        }
      );

      devShells = forDevSystems (
        system:
        let
          pkgs = pkgsFor system;
        in
        {
          default = pkgs.mkShell {
            packages = [
              pkgs.go
              pkgs.imagemagick
              pkgs.postgresql
              pkgs.valkey
            ];
          };
        }
      );

      formatter = forDevSystems (system: (pkgsFor system).nixfmt-tree);
    };
}
