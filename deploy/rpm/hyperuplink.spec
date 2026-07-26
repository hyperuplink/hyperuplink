# RPM spec for Hyperuplink.
#
# Builds the single, statically-linked Go binary from source and installs it
# together with a production config, a systemd unit and a dedicated service
# user. It targets Fedora, but is written to also build and install cleanly on
# RHEL-family (Rocky/Alma/RHEL) and hopefully openSUSE.
# See deploy/rpm/README.md.
#
# The version, commit and build date are injected at build time. The CI
# workflow overrides them from the git tag:
#
#   rpmbuild -bb deploy/rpm/hyperuplink.spec \
#     --define "version 1.2.3" \
#     --define "commit  <sha>" \
#     --define "date    <iso-8601>"

%global debug_package %{nil}
%global __os_install_post %{nil}

# Overridable at build time, the defaults mirror the Gentoo ebuild so a plain
# `rpmbuild -bb` still produces a sane, self-identifying binary.
%{!?version: %global version 0.1.0}
%{!?commit:  %global commit  v%{version}}
%{!?date:    %global date    release}

Name:           hyperuplink
Version:        %{version}
Release:        1%{?dist}
Summary:        A super high speed internet bulletin board

# SEGV is a custom, non-OSI license. Its full text ships as the LICENSE file.
License:        SEGV
URL:            https://hyperup.link
Source0:        %{name}-%{version}.tar.gz
Source1:        %{name}.sysusers

BuildRequires:  golang
BuildRequires:  systemd-rpm-macros
%{?systemd_requires}

# useradd/groupadd at these paths on Fedora/RHEL (shadow-utils) and
# openSUSE (shadow) alike, so depend on the file rather than the package name.
Requires(pre):  /usr/sbin/useradd
Requires(pre):  /usr/sbin/groupadd

# convert(1) from ImageMagick is shelled out to for image processing.
Requires:       /usr/bin/convert

%description
Hyperuplink is a super high speed internet bulletin board. It renders modern,
JavaScript-free HTML5 directly from a PostgreSQL database and uses Redis/Valkey
for sessions.

All assets (static files, views, locales, templates, docs and migrations) are
embedded into the executable, so the package installs just the binary, a
configuration file and a systemd unit. A running PostgreSQL server and a
Redis/Valkey server, reachable per /etc/hyperuplink.toml, are required at
runtime. They are intentionally not hard dependencies because they may live on
another host.

%prep
%autosetup -n %{name}-%{version}

# The binary installs to /usr/bin here, not /usr/local/bin.
sed -i 's|/usr/local/bin/hyperuplink|/usr/bin/hyperuplink|g' \
    deploy/init/hyperuplink.service

# The shipped production config targets container hostnames.
sed -i \
    -e 's|"valkey:6379"|"127.0.0.1:6379"|g' \
    -e 's|@postgres:5432|@localhost:5432|g' \
    -e 's|http://minio:9000|http://localhost:9000|g' \
    deploy/hyperuplink.toml

%build
# Map the RPM target CPU onto a Go GOARCH so the CI workflow can cross-build
# foreign-arch RPMs on an x86_64 host with `rpmbuild --target <arch>`.
case %{_target_cpu} in
    x86_64)  export GOARCH=amd64             ;;
    aarch64) export GOARCH=arm64             ;;
    ppc64le) export GOARCH=ppc64le           ;;
    s390x)   export GOARCH=s390x             ;;
    riscv64) export GOARCH=riscv64           ;;
    i?86)    export GOARCH=386               ;;
    armv7hl) export GOARCH=arm; export GOARM=7 ;;
    *) echo "unsupported target CPU: %{_target_cpu}" >&2; exit 1 ;;
esac

export GOOS=linux
export CGO_ENABLED=0
export GOFLAGS=-mod=mod
export GOTOOLCHAIN=auto
export GOCACHE="%{_builddir}/.gocache"
export GOMODCACHE="%{_builddir}/.gomodcache"

# Every asset is compiled in via go:embed, so this yields a fully self-contained
# executable. -buildvcs=false because the source tarball carries no .git.
go build \
    -trimpath \
    -buildvcs=false \
    -ldflags "-s -w \
        -X main.Version=%{version} \
        -X main.Commit=%{commit} \
        -X main.Date=%{date}" \
    -o %{name} .

%install
install -Dpm 0755 %{name} %{buildroot}%{_bindir}/%{name}

install -Dpm 0640 deploy/hyperuplink.toml \
    %{buildroot}%{_sysconfdir}/%{name}.toml

install -Dpm 0644 deploy/init/hyperuplink.service \
    %{buildroot}%{_unitdir}/%{name}.service

install -Dpm 0644 %{SOURCE1} %{buildroot}%{_sysusersdir}/%{name}.conf

# Shipped renamed because it shares the README.md basename with the top-level
# one that %doc installs into the same directory.
install -Dpm 0644 deploy/init/README.md \
    %{buildroot}%{_pkgdocdir}/init.README.md

# Runtime data / local media storage, owned by the service user. Matches
# [[Storage]] Local.Path in the installed config and StateDirectory=
# in the unit.
install -dm 0750 %{buildroot}%{_sharedstatedir}/%{name}/media

%pre
getent group %{name} >/dev/null || groupadd -r %{name}
getent passwd %{name} >/dev/null || \
    useradd -r -g %{name} -d %{_sharedstatedir}/%{name} -s /usr/sbin/nologin \
        -c "Hyperuplink service" %{name}
exit 0

%post
%systemd_post %{name}.service

%preun
%systemd_preun %{name}.service

%postun
%systemd_postun_with_restart %{name}.service

%files
%license LICENSE
%doc README.md
%{_pkgdocdir}/init.README.md
%{_bindir}/%{name}
%config(noreplace) %attr(0640, root, %{name}) %{_sysconfdir}/%{name}.toml
%{_unitdir}/%{name}.service
%{_sysusersdir}/%{name}.conf
%dir %attr(0750, %{name}, %{name}) %{_sharedstatedir}/%{name}
%dir %attr(0750, %{name}, %{name}) %{_sharedstatedir}/%{name}/media

%changelog
* Mon Jul 13 2026 Hyperuplink Authors - 0.1.0-1
- Initial RPM packaging.
