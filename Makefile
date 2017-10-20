default: build

build:
	GOOS=darwin GOARCH=amd64 go build -o dist/hclq-macos-amd64
	GOOS=linux GOARCH=amd64 go build -o dist/hclq-linux-amd64

build-win:
	GOOS=windows GOARCH=amd64 go build -o dist/hclq-windows-amd64

install:
	go install

clean:
	rm -rf dist

.PHONY: build clean install
