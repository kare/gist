
NAME := kkn.fi/cmd/gist
.PHONY: build install clean test vet lint errcheck check

build:
	go build $(NAME)

install:
	go install $(NAME)

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

