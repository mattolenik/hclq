# Version and linker flags
# This will return either the current tag, branch, or commit hash of this repo.
VERSION=$(shell echo $$(ver=$$(git tag -l --points-at HEAD) && [ -z $$ver ] && ver=$$(git describe --always --dirty); printf $$ver))
LDFLAGS=-s -w -X github.com/mattolenik/hclq/cmd.version=${VERSION}
GOPATH=$(HOME)/go
IS_PUBLISH=$(APPVEYOR_REPO_TAG)
BUILD_CMD=go build -ldflags="${LDFLAGS}"


default: test build readme

build:
	go build -i -ldflags="${LDFLAGS}" -gcflags='-N -l' -o dist/hclq

clean:
	rm -rf dist/ vendor/

dist: get
	# Delete files from testing
	rm -rf dist && mkdir -p dist
	# Make available for all the same platforms as Terraform.
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

get:
	go get -u github.com/golang/dep/cmd/dep
	$(GOPATH)/bin/dep ensure
	# GitHub release tool
	go get -u github.com/tcnksm/ghr
	go get -u github.com/jstemmer/go-junit-report

install: get
	go install -ldflags="${LDFLAGS}"

publish: readme test dist
	( \
		if [ -n "$(IS_PUBLISH)" ]; then \
			ghr -replace -delete -u "$$GITHUB_USER" ${VERSION} dist/; \
		fi; \
	)

readme:
	erb README.md.rb > README.md

test: get build
	HCLQ_BIN=$$(pwd)/dist/hclq go test -v "./..."

testci: get build
	mkdir -p test
	HCLQ_BIN=$$(pwd)/dist/hclq go test -v "./..." | go-junit-report | tee test/TEST.xml


.PHONY: get dist publish build install test clean
