BINARY := googlemapscli
PKG    := ./cmd/googlemapscli
VERSION ?= dev

.PHONY: build test lint coverage clean

build:
	go build -trimpath -ldflags "-X github.com/krishnashahane/gmapcli/internal/terminal.BuildVersion=$(VERSION)" -o $(BINARY) $(PKG)

test:
	go test ./... -count=1

lint:
	golangci-lint run ./...

clean:
	rm -f $(BINARY)

coverage:
	go test ./... -coverprofile=coverage.txt -count=1
	go tool cover -func=coverage.txt
