
NAME := kkn.fi/cmd/gist
.PHONY: build install clean test vet lint errcheck check

VERSION=$(shell cat version.txt)
DATE=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

build:
	go build -ldflags "-X main.version=${VERSION} -X main.date=${DATE}" $(NAME)

install:
	go install -ldflags "-X main.version=${VERSION} -X main.date=${DATE}" $(NAME)

clean:
	@rm -rf gist

test:
	go test $(NAME)/...

vet:
	go vet $(NAME)/...

lint:
	golint $(NAME)/...

errcheck:
	errcheck $(NAME)/...

check: vet lint errcheck test

