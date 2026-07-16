# Nix package for Hyperuplink.
#
# Builds the single, statically-linked Go binary from source and wraps it so
# that ImageMagick's `convert` is found on PATH. Every asset (static/, views/,
# locales/, templates/, docs/, migrations/) is compiled in via go:embed, so the
# result is one self-contained executable.
# See deploy/nix/README.md.
#
# The version, commit and build date are injected at build time, mirroring the
# Gentoo ebuild and the RPM spec. The flake passes the checked-out revision.
# By hand it is:
#
#   pkgs.hyperuplink.override { version = "1.2.3"; rev = "<sha>"; }
#
{
  lib,
  buildGoModule,
  makeBinaryWrapper,
  imagemagick,
  versionCheckHook,
  version ? "0.1.2",
  rev ? "v${version}",
  date ? "release",
}:

let
  goipath = "xn--gckvb8fzb.com/hyperuplink";

  # Only what the compiler and go:embed actually read, explicitly keeping out
  # deploy/, testdata/ and the working tree's build/
  # so editing this packaging does not rebuild the binary. Override it
  # with `.overrideAttrs (o: { src = ...; })`.
  source = lib.fileset.toSource {
    root = ../..;
    fileset = lib.fileset.unions [
      ../../go.mod
      ../../go.sum
      ../../hyperuplink.go
      ../../cron
      ../../docs
      ../../errs
      ../../http
      ../../locales
      ../../logic
      ../../migrations
      ../../models
      ../../runtime
      ../../services
      ../../static
      ../../templates
      ../../tools
      ../../views
      ../../worker
    ];
  };
in
buildGoModule (finalAttrs: {
  pname = "hyperuplink";
  inherit version;
  src = source;

  vendorHash = "sha256-cOggrFneuTep1+bdh+dcrXXN5AxAw2YSV7k5LxMuQ2M=";

  # go.mod requires Go 1.26.4 and the nixpkgs Go builder pins GOTOOLCHAIN=local,
  # so nixpkgs' Go has to be new enough.
  env.CGO_ENABLED = 0;

  subPackages = [ "." ];

  ldflags = [
    "-s"
    "-w"
    "-X ${goipath}/runtime.Version=${finalAttrs.version}"
    "-X ${goipath}/runtime.Commit=${rev}"
    "-X ${goipath}/runtime.Date=${date}"
  ];

  nativeBuildInputs = [ makeBinaryWrapper ];

  # convert(1) is resolved with exec.LookPath at startup (services/magick), so
  # it has to be on the service's PATH, not merely present at build time.
  postInstall = ''
    wrapProgram $out/bin/hyperuplink \
      --prefix PATH : ${lib.makeBinPath [ imagemagick ]}
  '';

  # The Go suite talks to a live PostgreSQL and Redis, which the build sandbox
  # has not got. The NixOS VM test in test.nix exercises the package instead.
  doCheck = false;

  doInstallCheck = true;
  nativeInstallCheckInputs = [ versionCheckHook ];
  versionCheckProgramArg = "-v";

  meta = {
    description = "A super high speed internet bulletin board";
    longDescription = ''
      Hyperuplink is a super high speed internet bulletin board. All assets
      (templates, static files, locales and SQL migrations) are embedded into
      the executable. It renders server-side HTML5/CSS with no client-side
      JavaScript and reads its data from a PostgreSQL database, using
      Redis/Valkey for background jobs.
    '';
    homepage = "https://xn--gckvb8fzb.com";
    downloadPage = "https://codeberg.org/hyperuplink/hyperuplink";
    # SEGV is a custom, non-OSI license. Section 2.1 grants copying and
    # distribution, but section 3 conditions the grant on ethical standards,
    # so is not free in the sense nixpkgs means it.
    # Redistributing the built binary is allowed.
    license = {
      shortName = "SEGV";
      fullName = "SEGV License, Version 1.0";
      url = "https://xn--gckvb8fzb.com/segv/";
      free = false;
      redistributable = true;
    };
    mainProgram = "hyperuplink";
    platforms = lib.platforms.unix;
    maintainers = [ ];
  };
})
