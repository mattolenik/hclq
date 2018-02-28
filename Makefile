# Version and linker flags
VERFLAG=$(shell echo -X main.version=$$(ver=$$(git tag -l --points-at HEAD) && [ -z $$ver ] && ver=$$(git describe --always --dirty); printf $$ver))
LDFLAGS=${VERFLAG} -s -w

# Dependency vars
UPX_URL=$(shell curl -sL https://api.github.com/repos/upx/upx/releases/latest | grep -e "browser_download_url.*amd64_linux" | awk -F'"' '{print $$4}')

default: test dist

get:
	go get -u github.com/golang/dep/cmd/dep
	$$GOPATH/bin/dep ensure

build: get
	GOOS=darwin  GOARCH=amd64 go build -i -ldflags="${LDFLAGS}" -o dist/hclq-darwin-amd64
	GOOS=linux   GOARCH=amd64 go build -i -ldflags="${LDFLAGS}" -o dist/hclq-linux-amd64
	GOOS=windows GOARCH=amd64 go build -i -ldflags="${LDFLAGS}" -o dist/hclq-windows-amd64

cideps:
	tar --version
	[ -z "$$CI" ] || curl -sSL ${UPX_URL} | tar xJ --wildcards --strip-components=1 "*/upx"

dist: build cideps
	[ -z "$$CI" ] || ./upx dist/*

debug:
	go build -i -gcflags='-N -l' -o dist/hclq

install: get
	go install -ldflags="${LDFLAGS}"

test: get debug
	HCLQ_BIN=$$(pwd)/dist/hclq go test -v "./..."

clean:
	rm -rf dist

.PHONY: clean debug debug-test-cmd dist get install test
