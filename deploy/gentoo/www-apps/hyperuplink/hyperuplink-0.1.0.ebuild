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

if [[ ${PV} == *9999 ]]; then
  inherit git-r3

  EGIT_REPO_URI="https://github.com/hyperuplink/hyperuplink.git"
  KEYWORDS=""

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

LICENSE="SEGV-1.1"
SLOT="0"

BDEPEND+=" >=dev-lang/go-1.26.4:="

DEPEND="
	acct-group/hyperuplink
	acct-user/hyperuplink
"
RDEPEND="
	${DEPEND}
	media-gfx/imagemagick
"

src_prepare() {
  default

  sed -i \
    -e 's|/usr/local/bin/hyperuplink|/usr/bin/hyperuplink|g' \
    deploy/init/hyperuplink.openrc \
    deploy/init/hyperuplink.service || die

  sed -i \
    -e 's|"valkey:6379"|"127.0.0.1:6379"|g' \
    -e 's|@postgres:5432|@localhost:5432|g' \
    -e 's|http://minio:9000|http://localhost:9000|g' \
    deploy/hyperuplink.toml || die
}

src_compile() {
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

  newinitd deploy/init/hyperuplink.openrc "${PN}"
  systemd_dounit deploy/init/hyperuplink.service

  keepdir /var/lib/${PN}/media
  fowners -R ${PN}:${PN} /var/lib/${PN}
  fperms 0750 /var/lib/${PN}
  fperms 0750 /var/lib/${PN}/media

  einstalldocs
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
