# Overlay that adds `hyperuplink` to a nixpkgs instance.
#
# On NixOS:
#
#   nixpkgs.overlays = [ (import /path/to/hyperuplink/deploy/nix/overlay.nix) ];
#
# The NixOS module in module.nix defaults to `pkgs.hyperuplink`, so this is what
# wires the two together when the flake is not used.
final: prev: {
  hyperuplink = final.callPackage ./package.nix { };
}
