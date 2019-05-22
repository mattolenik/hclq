# Version and linker flags
# This will return either the current tag, branch, or commit hash of this repo.
VERSION         = $(shell echo $$(ver=$$(git tag -l --points-at HEAD) && [ -z $$ver ] && ver=$$(git describe --always --dirty); printf $$ver))
LDFLAGS         = -s -w -X github.com/mattolenik/hclq/cmd.version=${VERSION}
PROJECT_ROOT    = $(shell cd -P -- '$(shell dirname -- "$0")' && pwd -P)
DIST            = dist
IS_PUBLISH      = $(APPVEYOR_REPO_TAG)
BUILD_CMD       = go build -mod=vendor -ldflags="${LDFLAGS}"
# Build tools
GHR             := github.com/tcnksm/ghr
GO_JUNIT_REPORT := github.com/jstemmer/go-junit-report

SOURCE := $(shell find $(PROJECT_ROOT) -name '*.go')
BINS   := $(shell find $(DIST) -name 'hclq-*')

default: build test README.md

build: dist/hclq

dist/hclq: $(SOURCE)
	$(BUILD_CMD) -i -gcflags='-N -l' -o dist/hclq

clean:
	rm -rf dist && mkdir dist

dist: dist/hclq-%
	cd dist && shasum -a 256 hclq-* > hclq-shasums

dist/hclq-darwin-amd64:
	export GOOS=darwin  GOARCH=amd64; $(BUILD_CMD) -o "$@"
	cd "$(@D)" && shasum -a 256 "$@" >> hclq-shasums

dist/hclq-freebsd-386:
	export GOOS=freebsd GOARCH=386  ; $(BUILD_CMD) -o "$@"
	cd "$(@D)" && shasum -a 256 "$@" >> hclq-shasums

dist/hclq-freebsd-amd64:
	export GOOS=freebsd GOARCH=amd64; $(BUILD_CMD) -o "$@"
	cd "$(@D)" && shasum -a 256 "$@" >> hclq-shasums

dist/hclq-freebsd-arm:
	export GOOS=freebsd GOARCH=arm  ; $(BUILD_CMD) -o "$@"
	cd "$(@D)" && shasum -a 256 "$@" >> hclq-shasums


dist/hclq-linux-386:
	export GOOS=linux   GOARCH=386  ; $(BUILD_CMD) -o "$@"
	cd "$(@D)" && shasum -a 256 "$@" >> hclq-shasums

dist/hclq-linux-amd64:
	export GOOS=linux   GOARCH=amd64; $(BUILD_CMD) -o "$@"
	cd "$(@D)" && shasum -a 256 "$@" >> hclq-shasums

dist/hclq-linux-arm:
	export GOOS=linux   GOARCH=arm  ; $(BUILD_CMD) -o "$@"
	cd "$(@D)" && shasum -a 256 "$@" >> hclq-shasums

dist/hclq-openbsd-386:
	export GOOS=openbsd GOARCH=amd64; $(BUILD_CMD) -o "$@"
	cd "$(@D)" && shasum -a 256 "$@" >> hclq-shasums

dist/hclq-openbsd-amd64:
	export GOOS=openbsd GOARCH=386  ; $(BUILD_CMD) -o "$@"
	cd "$(@D)" && shasum -a 256 "$@" >> hclq-shasums

dist/hclq-solaris-amd64:
	export GOOS=solaris GOARCH=amd64; $(BUILD_CMD) -o "$@"
	cd "$(@D)" && shasum -a 256 "$@" >> hclq-shasums

dist/hclq-windows-386:
	export GOOS=windows GOARCH=386  ; $(BUILD_CMD) -o "$@"
	cd "$(@D)" && shasum -a 256 "$@" >> hclq-shasums

dist/hclq-windows-amd64:
	export GOOS=windows GOARCH=amd64; $(BUILD_CMD) -o "$@"
	cd "$(@D)" && shasum -a 256 "$@" >> hclq-shasums

install:
	go install -mod=vendor ldflags="${LDFLAGS}"

# GitHub Release Tool
$(GHR):
	go install -mod=vendor $(GHR)

# Translates Go test results to JUnit XML
$(GO_JUNIT_REPORT):
	go install -mod=vendor $(GO_JUNIT_REPORT)

publish: $(GHR) test dist
	[ -n "$(IS_PUBLISH)" ] && ghr -replace -delete -u "$$GITHUB_USER" ${VERSION} dist/

readme: README.md
README.md: README.md.rb
	@if [ ! -f .git/hooks/pre-commit ]; then printf "Missing pre-commit hook for readme, be sure to copy it from hclq-pages repo"; exit 1; fi
	erb README.md.rb > README.md

test: $(GO_JUNIT_REPORT) build
	@mkdir -p test
	HCLQ_BIN=$(PROJECT_ROOT)/dist/hclq go test -mod=vendor -v "./..." | tee /dev/tty | go-junit-report > test/TEST.xml


.PHONY: clean install publish test testci
