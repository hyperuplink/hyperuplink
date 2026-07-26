.PHONY: all help test swagger build db\:drop db\:seed run release manual\:screenshots site\:screenshots
PWD := $(shell pwd)
GOPATH := $(shell go env GOPATH)

GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

NAME := hyperuplink
VERSION := $(shell git describe --tags 2>/dev/null || echo "dev")

# The ebuild takes its version from its filename and the RPM spec from
# --define, but a Nix flake cannot read the git tag (`self` exposes rev, never
# the tag), so deploy/nix/package.nix has to carry it in the tree. `make
# release` writes it there and derives the tag from it, which is why the two
# cannot drift apart.
NIXPKG := deploy/nix/package.nix
NIXVERSION = sed -n -E 's/^[[:space:]]*version \? "([^"]*)".*/\1/p' $(NIXPKG)

# The live ebuild is the canonical one and branches on ${PV}, which Portage
# takes from the filename, so a release ebuild is a plain copy of it.
EBUILDDIR := deploy/gentoo/www-apps/hyperuplink
LIVEEBUILD := $(EBUILDDIR)/hyperuplink-9999.ebuild
EBUILD = $(EBUILDDIR)/hyperuplink-$(VERSION).ebuild
COMMIT := $(shell git rev-parse --verify HEAD)
DATE := $(shell date)

# Where the hyperup.link Hugo project keeps the screenshots its landing page
# shows. It is a sibling checkout, so override this if yours sits elsewhere.
SITE ?= $(abspath $(PWD)/../pub/static/screenshots)

all: build

help: ## print this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' Makefile | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-40s\033[0m %s\n", $$1, $$2}'

test: ## test
	go test -v ./...

# The specification lands in docs/, which is already embedded, so the build
# picks it up without any further wiring. Only the JSON is generated: the
# handler reads the embedded file, so swag's docs.go would be dead weight.
swagger: ## generate the OpenAPI specification the API serves at /_internal/swagger
	@go tool swag init \
		--quiet \
		--generalInfo http/api/api._.go \
		--dir ./ \
		--exclude build,deploy,docs,locales,migrations,static,templates,testdata,tools,views \
		--output docs \
		--outputTypes json \
		--parseInternal \
		--parseDependency \
		--parseDepth 2

# Version, Commit and Date are declared in hyperuplink.go. The linker names a
# main package's symbols `main.X` whatever the module path is, so an -X that
# spells out the import path silently sets nothing and `-v` prints blanks.
build: swagger ## build
	@echo "Building with the following parameters:"
	@echo "VERSION = $(VERSION)"
	@echo "COMMIT  = $(COMMIT)"
	@echo "DATE    = $(DATE)"
	@CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags "-s -w -X \"main.Version=${VERSION}\" -X \"main.Commit=${COMMIT}\" -X \"main.Date=${DATE}\"" -o $(PWD)/build/$(NAME)

db\:drop: ## clear development database (drop all tables and content, keep the database)
	psql -h localhost -p 5432 -U postgres -d hyperuplink_dev -c "DROP SCHEMA public CASCADE;" -c "CREATE SCHEMA public;"

db\:seed: build ## seed the development database with users and a board
	@tools/seed/seed.sh "file://$(PWD)/hyperuplink.toml" hyperuplink_dev

run: build ## build and run
	./build/hyperuplink -c "file://$(PWD)/hyperuplink.toml"

manual\:screenshots: ## regenerate the screenshots embedded in the manual
	@tools/screenshots/run.sh -set manual

site\:screenshots: ## regenerate the screenshots on the hyperup.link website (SITE=...)
	@test -d "$(dir $(SITE))" \
		|| { echo "error: '$(dir $(SITE))' not found, set SITE=/path/to/pub/static/screenshots"; exit 1; }
	@tools/screenshots/run.sh -set site -out "$(SITE)"

release: ## bump the nix package, commit and tag a release (VERSION=x.y.z)
	@echo "$(VERSION)" | grep -Eq '^[0-9]+\.[0-9]+\.[0-9]+$$' \
		|| { echo "VERSION must be x.y.z, e.g. make release VERSION=0.1.3"; exit 1; }
	@test -z "$$(git status --porcelain --untracked-files=no)" \
		|| { echo "there are uncommitted changes, commit or stash them first:"; \
		     git status --short --untracked-files=no; exit 1; }
	@git rev-parse -q --verify "refs/tags/v$(VERSION)" >/dev/null \
		&& { echo "tag v$(VERSION) already exists"; exit 1; } || true
	@sed -i -E 's/^([[:space:]]*version \? )"[^"]*"/\1"$(VERSION)"/' $(NIXPKG)
	@test "$$($(NIXVERSION))" = "$(VERSION)" \
		|| { git checkout -- $(NIXPKG); \
		     echo "could not set the version in $(NIXPKG), has the argument been renamed?"; exit 1; }
	@cp $(LIVEEBUILD) $(EBUILD)
	@git add $(NIXPKG) $(EBUILD) && git commit -S -q -m "Release v$(VERSION)" \
		|| { git reset -q -- $(NIXPKG) $(EBUILD); git checkout -- $(NIXPKG); rm -f $(EBUILD); \
		     echo "the commit failed, $(NIXPKG) and $(EBUILD) have been rolled back"; exit 1; }
	@git tag -s -m "v$(VERSION)" "v$(VERSION)"
	@echo "Tagged v$(VERSION), with $(NIXPKG) bumped and $(EBUILD) added."
	@echo "Nothing is published until you run:"
	@echo "    git push --follow-tags"
