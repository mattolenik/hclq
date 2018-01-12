LDFLAGS=$(shell echo -X main.version=$$(ver=$$(git tag -l --points-at HEAD) && [ -z $$ver ] && ver=$$(git describe --always --dirty); printf $$ver))
default: build

get:
	go get -u github.com/golang/dep/cmd/dep
	$$GOPATH/bin/dep ensure

build: get
	go build -i -ldflags="${LDFLAGS}" -gcflags='-N -l' -o dist/hclq

dist: get
	GOOS=darwin  GOARCH=amd64 go build -i -ldflags="${LDFLAGS}" -o dist/hclq-darwin-amd64
	GOOS=linux   GOARCH=amd64 go build -i -ldflags="${LDFLAGS}" -o dist/hclq-linux-amd64
	GOOS=windows GOARCH=amd64 go build -i -ldflags="${LDFLAGS}" -o dist/hclq-windows-amd64

test: build
	DLV_BIN=$$GOPATH/bin/dlv HCLQ_BIN=$$(pwd)/dist/hclq go test -v "./..."

debug-test-cmd: build
	go test -c "github.com/mattolenik/hclq/cmd" -o test/cmd.test
	dlv --listen=:2345 --headless=true --api-version=2 exec test/cmd.test

install:
	go install -ldflags="${LDFLAGS}"

clean:
	rm -rf dist

.PHONY: build dist clean get install test
