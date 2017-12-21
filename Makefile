
IMPORT_PATH := kkn.fi/cmd/gist

GOMETALINTER := $(GOPATH)/bin/gometalinter
GORELEASER := $(GOPATH)/bin/goreleaser

VERSION=$(shell git describe --tags --always --dirty="-dev")
DATE=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
VERSION_FLAGS := -ldflags='-X "main.version=$(VERSION)" -X "main.date=$(DATE)"'

.PHONY: build
build:
	go build $(VERSION_FLAGS) $(IMPORT_PATH)

.PHONY: install
install:
	go install $(VERSION_FLAGS) $(IMPORT_PATH)

.PHONY: clean
clean:
	@rm -rf gist dist

.PHONY: test
test:
	go test $(IMPORT_PATH)/...

.PHONY: lint
lint: $(GOMETALINTER)
	gometalinter --vendor ./...

$(GOMETALINTER):
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install

.PHONY: release
release: $(GORELEASER) clean
	goreleaser

$(GORELEASER):
	go get -u github.com/goreleaser/goreleaser

.PHONY: setup
setup:
	go get -u github.com/golang/dep/cmd/dep
