# Version and linker flags
VERSION=$(shell echo $$(ver=$$(git tag -l --points-at HEAD) && [ -z $$ver ] && ver=$$(git describe --always --dirty); printf $$ver))
LDFLAGS=-s -w -X github.com/mattolenik/hclq/cmd.version=${VERSION}
GOOS=darwin linux windows
GOARCH=amd64

# Dependency vars
UPX_URL=$(shell curl -sL https://api.github.com/repos/upx/upx/releases/latest | grep -e "browser_download_url.*amd64_linux" | awk -F'"' '{print $$4}')

default: test build

build:
	go build -i -ldflags="${LDFLAGS}" -gcflags='-N -l' -o dist/hclq

cideps:
	# Download and extract UPX
	[ -z "$$CI" ] || curl -sSL ${UPX_URL} | tar xJ --wildcards --strip-components=1 "*/upx"

clean:
	rm -rf dist/ vendor/

dist: get
	set -v; for goos in ${GOOS}; do GOOS=$$goos GOARCH=${GOARCH} go build -i -ldflags="${LDFLAGS}" -o dist/hclq-$$goos-${GOARCH}; done
	# Remove binary used for testing
	rm dist/hclq
	[ -n "$$CI" ] && ./upx dist/* || upx dist/*

get:
	go get -u github.com/golang/dep/cmd/dep
	$$GOPATH/bin/dep ensure
	# GitHub release tool
	go get -u github.com/tcnksm/ghr

install: get
	go install -ldflags="${LDFLAGS}"

brew:
	./mo homebrew/hclq.rb.mo > homebrew/hclq.rb

release: cideps test dist
	( \
		VERSION=${VERSION}; \
		LINUX_FILENAME="hclq-linux-amd64"; \
		DARWIN_FILENAME="hclq-darwin-amd64"; \
		LINUX_HASH=$$(shasum -a 256 dist/$$LINUX_FILENAME | awk '{print $$1}'); \
		DARWIN_HASH=$$(shasum -a 256 dist/$$DARWIN_FILENAME | awk '{print $$1}'); \
		shasum -a 256 dist/* > dist/hclq-${VERSION}-shasums; \
		[ -z "$$CI" ] && printf "CI var not set, skipping publish\n" && exit 0; \
		ghr -u "$$GITHUB_USER" -r hclq ${VERSION} dist/; \
	)

test: get build
	HCLQ_BIN=$$(pwd)/dist/hclq go test -v "./..."


.PHONY: get dist cideps release build install test clean
