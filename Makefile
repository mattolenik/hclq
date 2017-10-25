LDFLAGS=$(shell echo -X main.version=$$(git describe --always --dirty))
default: build

get:
	go get -u github.com/golang/dep/cmd/dep
	dep ensure

build: get
	GOOS=darwin  GOARCH=amd64 go build -i -ldflags="${LDFLAGS}" -o dist/hclq-macos-amd64
	GOOS=linux   GOARCH=amd64 go build -i -ldflags="${LDFLAGS}" -o dist/hclq-linux-amd64
	GOOS=windows GOARCH=amd64 go build -i -ldflags="${LDFLAGS}" -o dist/hclq-windows-amd64

build-brew: get
	GOOS=darwin  GOARCH=amd64 go build -i -ldflags="${LDFLAGS}" -o dist/hclq

install:
	go install -ldflags="${LDFLAGS}"

clean:
	rm -rf dist

.PHONY: build clean get install
