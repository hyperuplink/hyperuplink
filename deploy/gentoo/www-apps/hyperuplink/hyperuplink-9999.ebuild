# Copyright 2025 Hyperuplink Authors
# Distributed under the terms of the SEGV License, Version 1.0

EAPI=8

inherit go-module systemd

DESCRIPTION="A super high speed internet bulletin board"
HOMEPAGE="
	https://xn--gckvb8fzb.com
	https://codeberg.org/hyperuplink/hyperuplink
	https://github.com/hyperuplink/hyperuplink
"

# This one file serves the live version and every release alike.
# `make release` copies it to hyperuplink-<version>.ebuild unchanged, and the
# branches below are all that the copy needs, since ${PV} comes from the
# filename.
if [[ ${PV} == 9999 ]]; then
  # Build straight from the GitHub default branch. There is no source archive
  # and no deps tarball, and go-module_live_vendor populates a vendor/
  # directory during src_unpack.
  inherit git-r3

  EGIT_REPO_URI="https://github.com/hyperuplink/hyperuplink.git"
  KEYWORDS=""

  src_unpack() {
    git-r3_src_unpack
    # Vendors the modules (needs network, hence src_unpack) so the sandboxed
    # src_compile can build fully offline against vendor/.
    go-module_live_vendor
  }
else
  # The main source archive is GitHub's auto-generated tag tarball.
  # ${P}-deps.tar.xz is the Go module cache, and it is not produced by GitHub
  # automatically. The Release workflow (.github/workflows/release.yml) builds
  # it and attaches it to the GitHub release. See deploy/gentoo/README.md to
  # build it by hand.
  #
  # There is no src_unpack here on purpose, because go-module.eclass unpacks
  # the deps tarball into the module cache itself.
  SRC_URI="
		https://github.com/hyperuplink/hyperuplink/archive/v${PV}.tar.gz -> ${P}.tar.gz
		https://github.com/hyperuplink/hyperuplink/releases/download/v${PV}/${P}-deps.tar.xz
	"

  # GitHub's archive extracts to ${PN}-${PV}, which is go-module.eclass's
  # default S (${WORKDIR}/${P}), so no explicit S is needed.

  KEYWORDS="~amd64 ~arm ~arm64 ~ppc64 ~riscv ~x86"
fi

# SEGV is a custom, non-OSI license. The full text ships as licenses/SEGV in
# this overlay, because ::gentoo does not carry it.
LICENSE="SEGV"
SLOT="0"

# go.mod requires this Go version, and GOTOOLCHAIN=local
# (set by go-module.eclass) forbids downloading a newer toolchain, so the
# system Go must be new enough.
BDEPEND=">=dev-lang/go-1.26.4"

# The service user/group are provided by the acct-* packages in this overlay.
DEPEND="
	acct-group/hyperuplink
	acct-user/hyperuplink
"
# convert(1) from ImageMagick is shelled out to for image processing.
RDEPEND="
	${DEPEND}
	media-gfx/imagemagick
"

src_prepare() {
  default

  # Gentoo installs the binary to /usr/bin
  sed -i \
    -e 's|/usr/local/bin/hyperuplink|/usr/bin/hyperuplink|g' \
    deploy/init/hyperuplink.openrc \
    deploy/init/hyperuplink.service || die

  # The shipped production config targets Docker service hostnames, so we
  # rewrite them to loopback, so that a single-host install works out of the
  # box. Edit /etc/hyperuplink.toml afterwards.
  sed -i \
    -e 's|"valkey:6379"|"127.0.0.1:6379"|g' \
    -e 's|@postgres:5432|@localhost:5432|g' \
    -e 's|http://minio:9000|http://localhost:9000|g' \
    deploy/hyperuplink.toml || die
}

src_compile() {
  # All embedded assets (static/, views/, locales/, templates/, docs/,
  # migrations/) are compiled into the binary via go:embed, so the single
  # executable is fully self-contained.
  local modpath="xn--gckvb8fzb.com/hyperuplink"
  local ldflags=(
    -s -w
    -X "${modpath}/runtime.Version=${PV}"
  )
  local goflags=()

  if [[ ${PV} == 9999 ]]; then
    ldflags+=(
      -X "${modpath}/runtime.Commit=${EGIT_VERSION}"
      -X "${modpath}/runtime.Date=live"
    )
    # -mod=vendor (explicit, overriding the eclass GOFLAGS) forces the offline
    # build against the vendor/ tree created in src_unpack.
    goflags+=(-mod=vendor)
  else
    ldflags+=(
      -X "${modpath}/runtime.Commit=v${PV}"
      -X "${modpath}/runtime.Date=release"
    )
  fi

  CGO_ENABLED=0 go build "${goflags[@]}" -trimpath \
    -ldflags "${ldflags[*]}" -o "${PN}" . || die
}

src_install() {
  dobin "${PN}"

  # Configuration under /etc.
  insinto /etc
  newins deploy/hyperuplink.toml hyperuplink.toml

  # Init integration for both supported init systems.
  newinitd deploy/init/hyperuplink.openrc "${PN}"
  systemd_dounit deploy/init/hyperuplink.service

  # Runtime data / local media storage, owned by the service user.
  # Matches [[Storage]] Local.Path in the installed config.
  keepdir /var/lib/${PN}/media
  fowners -R ${PN}:${PN} /var/lib/${PN}
  fperms 0750 /var/lib/${PN}
  fperms 0750 /var/lib/${PN}/media

  einstalldocs
  dodoc deploy/init/README.md
}

pkg_postinst() {
  elog "Hyperuplink needs a running PostgreSQL server and a Redis/Valkey"
  elog "server reachable per /etc/hyperuplink.toml. They are intentionally"
  elog "not hard dependencies (they may live on another host)."
  elog
  elog "Before first start:"
  elog "  1. Edit /etc/hyperuplink.toml (DB DSN, Redis, PromoteAdmin, etc.)."
  elog "  2. Run the database and Redis/Valkey services."
  elog
  elog "OpenRC:"
  elog "  rc-update add ${PN} default"
  elog "  rc-service ${PN} start"
  elog
  elog "systemd:"
  elog "  systemctl enable --now ${PN}.service"
  elog
  elog "The service listens on port 3000 by default. Terminate TLS in front"
  elog "of it with a reverse proxy, especially in Mode = \"production\"."
}
