BINARY := aminal
VERSION ?= vlatest

.PHONY: build
build: test 
	mkdir -p release
	go build -o release/$(BINARY)-$(VERSION)

.PHONY: test
test:
	go test -v ./...

.PHONY: install
install: build
	install -m 0755 release/$(BINARY)-$(VERSION) /usr/local/bin/aminal
