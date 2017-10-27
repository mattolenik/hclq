LDFLAGS=$(shell echo -X main.version=$$(ver=$$(git tag -l --points-at HEAD) && [ -z $$ver ] && ver=$$(git describe --always --dirty); printf $$ver))
default: build

get:
	go get -u github.com/golang/dep/cmd/dep
	$$GOPATH/bin/dep ensure

build: get
	GOOS=darwin  GOARCH=amd64 go build -i -ldflags="${LDFLAGS}" -o dist/hclq-darwin-amd64
	GOOS=linux   GOARCH=amd64 go build -i -ldflags="${LDFLAGS}" -o dist/hclq-linux-amd64
	GOOS=windows GOARCH=amd64 go build -i -ldflags="${LDFLAGS}" -o dist/hclq-windows-amd64

install:
	go install -ldflags="${LDFLAGS}"

clean:
	rm -rf dist

.PHONY: brew build clean get install
