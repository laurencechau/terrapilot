BINARY   := terrapilot
VERSION  := $(shell cat version.go | grep 'version =' | sed 's/.*"\(.*\)".*/\1/')
LDFLAGS  := -ldflags "-X main.version=$(VERSION)"

.PHONY: build test lint clean install

build:
	go build $(LDFLAGS) -o bin/$(BINARY) .

test:
	go test ./...

lint:
	go vet ./...

clean:
	rm -rf bin/

install:
	go install $(LDFLAGS) .
