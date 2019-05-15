# Version and linker flags
# This will return either the current tag, branch, or commit hash of this repo.
VERSION         = $(shell echo $$(ver=$$(git tag -l --points-at HEAD) && [ -z $$ver ] && ver=$$(git describe --always --dirty); printf $$ver))
LDFLAGS         = -s -w -X github.com/mattolenik/hclq/cmd.version=${VERSION}
PROJECT_ROOT    = $(shell cd -P -- '$(shell dirname -- "$0")' && pwd -P)
IS_PUBLISH      = $(APPVEYOR_REPO_TAG)
BUILD_CMD       = go build -mod=vendor -ldflags="${LDFLAGS}"
# Build tools
GHR             := github.com/tcnksm/ghr
GO_JUNIT_REPORT := github.com/jstemmer/go-junit-report

default: test build readme

build:
	$(BUILD_CMD) -i -gcflags='-N -l' -o dist/hclq

clean:
	rm -rf dist/

dist:
	# Delete files from testing
	rm -rf dist && mkdir -p dist
	# Make available for all the same platforms as Terraform
	export GOOS=darwin  GOARCH=amd64; $(BUILD_CMD) -o dist/hclq-$$GOOS-$$GOARCH
	export GOOS=freebsd GOARCH=amd64; $(BUILD_CMD) -o dist/hclq-$$GOOS-$$GOARCH
	export GOOS=freebsd GOARCH=386  ; $(BUILD_CMD) -o dist/hclq-$$GOOS-$$GOARCH
	export GOOS=freebsd GOARCH=arm  ; $(BUILD_CMD) -o dist/hclq-$$GOOS-$$GOARCH
	export GOOS=linux   GOARCH=amd64; $(BUILD_CMD) -o dist/hclq-$$GOOS-$$GOARCH
	export GOOS=linux   GOARCH=386  ; $(BUILD_CMD) -o dist/hclq-$$GOOS-$$GOARCH
	export GOOS=linux   GOARCH=arm  ; $(BUILD_CMD) -o dist/hclq-$$GOOS-$$GOARCH
	export GOOS=openbsd GOARCH=amd64; $(BUILD_CMD) -o dist/hclq-$$GOOS-$$GOARCH
	export GOOS=openbsd GOARCH=386  ; $(BUILD_CMD) -o dist/hclq-$$GOOS-$$GOARCH
	export GOOS=solaris GOARCH=amd64; $(BUILD_CMD) -o dist/hclq-$$GOOS-$$GOARCH
	export GOOS=windows GOARCH=amd64; $(BUILD_CMD) -o dist/hclq-$$GOOS-$$GOARCH
	export GOOS=windows GOARCH=386  ; $(BUILD_CMD) -o dist/hclq-$$GOOS-$$GOARCH
	cd dist && shasum -a 256 hclq-* > hclq-shasums

install: get
	go install -mod=vendor-ldflags="${LDFLAGS}"

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
	@[ ! -f .git/hooks/pre-commit ] && printf "Missing pre-commit hook for readme, be sure to copy it from hclq-pages repo" && exit 1
	erb README.md.rb > README.md

test: $(GO_JUNIT_REPOT) build
	#!/usr/bin/env bash
	set -euo pipefail
	mkdir -p test
	HCLQ_BIN=$(PROJECT_ROOT)/dist/hclq go test -mod=vendor -v "./..." | tee >(go-junit-report > test/TEST.xml)


.PHONY: build clean dist install publish test testci
