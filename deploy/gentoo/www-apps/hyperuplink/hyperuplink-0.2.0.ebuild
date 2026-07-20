# Copyright 2025 Hyperuplink Authors
# Distributed under the terms of the SEGV License, Version 1.1

EAPI=8

inherit go-module systemd

DESCRIPTION="A super high speed internet bulletin board"
HOMEPAGE="
	https://hyperup.link
	https://codeberg.org/hyperuplink/hyperuplink
	https://github.com/hyperuplink/hyperuplink
"

# This one file serves the live version and every release alike.
# `make release` copies it to hyperuplink-<version>.ebuild unchanged, and the
# branches below are all that the copy needs, since ${PV} comes from the
# filename.
if [[ ${PV} == *9999 ]]; then
  # Build straight from the GitHub default branch. There is no source archive
  # and no deps tarball, and go-module_live_vendor populates a vendor/
  # directory during src_unpack.
  inherit git-r3

  EGIT_REPO_URI="https://github.com/hyperuplink/hyperuplink.git"
  KEYWORDS=""

  # -mod=vendor (explicit, since the eclass GOFLAGS sets no module mode) forces
  # the offline build against the vendor/ tree created in src_unpack.
  HUP_GOFLAGS=(-mod=vendor)

  src_unpack() {
    git-r3_src_unpack
    go-env_set_compile_environment
    go-module_live_vendor
  }
else
  SRC_URI="
		https://github.com/hyperuplink/hyperuplink/archive/v${PV}.tar.gz -> ${P}.tar.gz
		https://github.com/hyperuplink/hyperuplink/releases/download/v${PV}/${P}-deps.tar.xz
	"

  HUP_GOFLAGS=()

  KEYWORDS="~amd64 ~arm ~arm64 ~ppc64 ~riscv ~x86"
fi

# SEGV is a custom, non-OSI license. The full text ships as licenses/SEGV-1.1
# in this overlay, because ::gentoo does not carry it.
LICENSE="SEGV-1.1"
SLOT="0"

# go.mod requires this Go version, and GOTOOLCHAIN=local
# (set by go-module.eclass) forbids downloading a newer toolchain, so the
# system Go must be new enough. Appended, because a plain assignment would drop
# the eclass's own BDEPEND.
BDEPEND+=" >=dev-lang/go-1.26.4:="

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

  if [[ ${PV} == *9999 ]]; then
    ldflags+=(
      -X "${modpath}/runtime.Commit=${EGIT_VERSION}"
      -X "${modpath}/runtime.Date=live"
    )
  else
    ldflags+=(
      -X "${modpath}/runtime.Commit=v${PV}"
      -X "${modpath}/runtime.Date=release"
    )
  fi

  local -x CGO_ENABLED=0
  ego build "${HUP_GOFLAGS[@]}" -trimpath -ldflags "${ldflags[*]}" -o "${PN}" .
}

src_test() {
  local -x CGO_ENABLED=0
  ego test "${HUP_GOFLAGS[@]}" ./...
}

src_install() {
  dobin "${PN}"

  insinto /etc
  insopts -m0640
  newins deploy/hyperuplink.toml hyperuplink.toml
  fowners root:${PN} /etc/hyperuplink.toml
  insopts -m0644

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
  # Renamed, because einstalldocs has already installed the project's own
  # README.md into the same directory.
  newdoc deploy/init/README.md README.init.md
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
