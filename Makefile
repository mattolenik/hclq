default: build

get:
	go get -u github.com/golang/dep/cmd/dep
	dep ensure

build: get
	GOOS=darwin GOARCH=amd64 go build -o dist/hclq-macos-amd64
	GOOS=linux GOARCH=amd64 go build -o dist/hclq-linux-amd64
	GOOS=windows GOARCH=amd64 go build -o dist/hclq-windows-amd64

install:
	go install

clean:
	rm -rf dist

.PHONY: build clean get install
